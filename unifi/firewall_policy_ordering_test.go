package unifi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
)

const (
	orderingSiteID = "11111111-1111-1111-1111-111111111111"
	orderingSrc    = "aaaaaaaa-0000-0000-0000-000000000001"
	orderingDst    = "bbbbbbbb-0000-0000-0000-000000000002"
	orderingPolA   = "c0000000-0000-0000-0000-00000000000a"
	orderingPolB   = "c0000000-0000-0000-0000-00000000000b"
	orderingPolC   = "c0000000-0000-0000-0000-00000000000c"
)

// orderingRequest records what the fake controller received.
type orderingRequest struct {
	Method, Path, APIKey string
	Query                map[string]string
	Body                 map[string]json.RawMessage
}

// firewallPolicyOrderingTestServer fakes an API-key UniFi OS controller serving
// the Integration API site listing and the ordering endpoint for one zone pair.
// Like the real controller (Network 10.6.x), a PUT is rejected with 400
// api.firewall.policy.ordering.not-all-firewall-policies-specified unless it
// lists every user-defined policy of the pair; otherwise it is stored and
// echoed back.
func firewallPolicyOrderingTestServer(
	t *testing.T,
	initial FirewallPolicyOrdering,
) (*httptest.Server, func() []orderingRequest) {
	t.Helper()
	var (
		mu      sync.Mutex
		reqs    []orderingRequest
		current = initial
	)
	path := integrationSitesPath + "/" + orderingSiteID + "/firewall/policies/ordering"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet && r.URL.Path == integrationSitesPath {
			_, _ = w.Write([]byte(twoSitesBody))
			return
		}
		if r.URL.Path != path {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		rec := orderingRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			APIKey: r.Header.Get("X-Api-Key"),
			Query:  map[string]string{},
		}
		for k := range r.URL.Query() {
			rec.Query[k] = r.URL.Query().Get(k)
		}
		if r.Method == http.MethodPut {
			_ = json.NewDecoder(r.Body).Decode(&rec.Body)
		}
		mu.Lock()
		defer mu.Unlock()
		reqs = append(reqs, rec)

		if rec.Query["sourceFirewallZoneId"] != orderingSrc ||
			rec.Query["destinationFirewallZoneId"] != orderingDst {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"statusCode":400,"statusName":"BAD_REQUEST",` +
				`"code":"api.request.argument-type-mismatch","message":"bad zone pair"}`))
			return
		}

		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(current)
		case http.MethodPut:
			raw, _ := json.Marshal(rec.Body)
			var in FirewallPolicyOrdering
			_ = json.Unmarshal(raw, &in)
			got := in.PolicyIDs()
			want := current.PolicyIDs()
			sent := make(map[string]bool, len(got))
			for _, id := range got {
				sent[id] = true
			}
			complete := len(got) == len(want)
			for _, id := range want {
				complete = complete && sent[id]
			}
			if !complete {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"statusCode":400,"statusName":"BAD_REQUEST",` +
					`"code":"api.firewall.policy.ordering.not-all-firewall-policies-specified",` +
					`"message":"Not all firewall policies specified",` +
					`"timestamp":"2026-10-31T10:00:00Z",` +
					`"requestPath":"/integration/v1/sites/` + orderingSiteID + `/firewall/policies/ordering",` +
					`"requestId":"3fa85f64-5717-4562-b3fc-2c963f66afa6"}`))
				return
			}
			current = in
			_ = json.NewEncoder(w).Encode(current)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, func() []orderingRequest {
		mu.Lock()
		defer mu.Unlock()
		return append([]orderingRequest(nil), reqs...)
	}
}

func newOrderingTestClient(t *testing.T, srv *httptest.Server) *ApiClient {
	t.Helper()
	c, err := New(context.Background(), &Config{BaseURL: srv.URL, APIKey: "test-key"})
	if err != nil {
		t.Fatalf("client init: %v", err)
	}
	return c
}

func threePolicyOrdering() FirewallPolicyOrdering {
	return FirewallPolicyOrdering{OrderedFirewallPolicyIDs: FirewallPolicyOrderedIDs{
		BeforeSystemDefined: []string{orderingPolA, orderingPolB},
		AfterSystemDefined:  []string{orderingPolC},
	}}
}

func TestGetFirewallPolicyOrdering_RoutesAndParses(t *testing.T) {
	srv, reqs := firewallPolicyOrderingTestServer(t, threePolicyOrdering())
	c := newOrderingTestClient(t, srv)

	got, err := c.GetFirewallPolicyOrdering(context.Background(), orderingSiteID, orderingSrc, orderingDst)
	if err != nil {
		t.Fatalf("GetFirewallPolicyOrdering errored: %v", err)
	}
	if want := threePolicyOrdering(); !reflect.DeepEqual(*got, want) {
		t.Errorf("got %+v, want %+v", *got, want)
	}
	if ids := got.PolicyIDs(); !reflect.DeepEqual(ids, []string{orderingPolA, orderingPolB, orderingPolC}) {
		t.Errorf("PolicyIDs = %v, want before-system ids then after-system ids", ids)
	}

	r := reqs()
	if len(r) != 1 {
		t.Fatalf("requests = %d, want 1", len(r))
	}
	if r[0].Method != http.MethodGet || r[0].APIKey != "test-key" {
		t.Errorf("request = %s with key %q, want GET with X-API-Key", r[0].Method, r[0].APIKey)
	}
}

func TestGetFirewallPolicyOrdering_ResolvesLegacySiteName(t *testing.T) {
	srv, reqs := firewallPolicyOrderingTestServer(t, threePolicyOrdering())
	c := newOrderingTestClient(t, srv)
	ctx := context.Background()

	// "default" is the legacy name of the site whose UUID is orderingSiteID.
	siteID, err := c.ResolveIntegrationSiteID(ctx, "default")
	if err != nil {
		t.Fatalf("ResolveIntegrationSiteID errored: %v", err)
	}
	if _, err := c.GetFirewallPolicyOrdering(ctx, siteID, orderingSrc, orderingDst); err != nil {
		t.Fatalf("GetFirewallPolicyOrdering errored: %v", err)
	}
	if r := reqs(); len(r) != 1 || !strings.Contains(r[0].Path, "/sites/"+orderingSiteID+"/") {
		t.Errorf("requests = %+v, want one ordering call on the resolved site UUID", r)
	}
}

func TestUpdateFirewallPolicyOrdering_SendsFullDocument(t *testing.T) {
	srv, reqs := firewallPolicyOrderingTestServer(t, threePolicyOrdering())
	c := newOrderingTestClient(t, srv)

	// Move every policy before the system-defined ones, in a new order. The
	// AfterSystemDefined list is left nil and must still be sent as [].
	desired := &FirewallPolicyOrdering{OrderedFirewallPolicyIDs: FirewallPolicyOrderedIDs{
		BeforeSystemDefined: []string{orderingPolC, orderingPolA, orderingPolB},
	}}
	got, err := c.UpdateFirewallPolicyOrdering(
		context.Background(), orderingSiteID, orderingSrc, orderingDst, desired,
	)
	if err != nil {
		t.Fatalf("UpdateFirewallPolicyOrdering errored: %v", err)
	}
	if !reflect.DeepEqual(got.OrderedFirewallPolicyIDs.BeforeSystemDefined, desired.OrderedFirewallPolicyIDs.BeforeSystemDefined) ||
		len(got.OrderedFirewallPolicyIDs.AfterSystemDefined) != 0 {
		t.Errorf("stored ordering = %+v, want %+v", *got, *desired)
	}
	if desired.OrderedFirewallPolicyIDs.AfterSystemDefined != nil {
		t.Error("caller's ordering was mutated")
	}

	r := reqs()
	if len(r) != 1 || r[0].Method != http.MethodPut {
		t.Fatalf("requests = %+v, want a single PUT", r)
	}
	if r[0].Query["sourceFirewallZoneId"] != orderingSrc ||
		r[0].Query["destinationFirewallZoneId"] != orderingDst {
		t.Errorf("query = %v, want the zone pair", r[0].Query)
	}
	var ids map[string]json.RawMessage
	if err := json.Unmarshal(r[0].Body["orderedFirewallPolicyIds"], &ids); err != nil {
		t.Fatalf("body = %v, want an orderedFirewallPolicyIds object: %v", r[0].Body, err)
	}
	if string(ids["afterSystemDefined"]) != "[]" {
		t.Errorf("afterSystemDefined = %s, want []", ids["afterSystemDefined"])
	}
}

func TestUpdateFirewallPolicyOrdering_SubsetSurfacesControllerError(t *testing.T) {
	srv, reqs := firewallPolicyOrderingTestServer(t, threePolicyOrdering())
	c := newOrderingTestClient(t, srv)

	subset := &FirewallPolicyOrdering{OrderedFirewallPolicyIDs: FirewallPolicyOrderedIDs{
		BeforeSystemDefined: []string{orderingPolB, orderingPolA},
	}}
	_, err := c.UpdateFirewallPolicyOrdering(
		context.Background(), orderingSiteID, orderingSrc, orderingDst, subset,
	)
	if err == nil {
		t.Fatal("expected an error for a partial ordering")
	}
	if !errors.Is(err, ErrFirewallPolicyOrderingIncomplete) {
		t.Errorf("error %v does not wrap ErrFirewallPolicyOrderingIncomplete", err)
	}
	apiErr := &APIError{}
	if !errors.As(err, &apiErr) {
		t.Fatalf("error %v does not wrap *APIError", err)
	}
	if apiErr.RC != "api.firewall.policy.ordering.not-all-firewall-policies-specified" {
		t.Errorf("APIError.RC = %q", apiErr.RC)
	}
	for _, want := range []string{"Not all firewall policies specified", "400"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q is missing %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "giving up") {
		t.Errorf("error is opaque: %v", err)
	}
	if n := len(reqs()); n != 1 {
		t.Errorf("requests = %d, want 1 (a 400 must not be retried)", n)
	}
}

func TestUpdateFirewallPolicyOrdering_OtherErrorsAreNotIncomplete(t *testing.T) {
	srv, _ := firewallPolicyOrderingTestServer(t, threePolicyOrdering())
	c := newOrderingTestClient(t, srv)

	order := threePolicyOrdering()
	_, err := c.UpdateFirewallPolicyOrdering(
		context.Background(), orderingSiteID, orderingSrc, "not-the-pair", &order,
	)
	if err == nil {
		t.Fatal("expected an error")
	}
	if errors.Is(err, ErrFirewallPolicyOrderingIncomplete) {
		t.Errorf("unrelated error %v wraps ErrFirewallPolicyOrderingIncomplete", err)
	}
	if !strings.Contains(err.Error(), "api.request.argument-type-mismatch: bad zone pair") {
		t.Errorf("error %q does not carry the controller message", err)
	}
}
