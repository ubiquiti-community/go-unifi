package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// trafficMatchingListTestServer stands up a new-style (UniFi OS) controller that
// serves the Integration API traffic-matching-lists endpoints under
// /proxy/network/integration/v1/sites/{siteID}/traffic-matching-lists. The
// handler echoes create/update bodies back (assigning an id on create) and
// tracks the last method/path seen so tests can assert routing.
func trafficMatchingListTestServer(
	t *testing.T,
	siteID string,
) (*httptest.Server, *struct{ Method, Path string }) {
	t.Helper()
	base := "/proxy/network/integration/v1/sites/" + siteID + "/traffic-matching-lists"
	last := &struct{ Method, Path string }{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		if !strings.HasPrefix(r.URL.Path, base) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		last.Method, last.Path = r.Method, r.URL.Path
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			if r.URL.Path == base {
				_, _ = w.Write([]byte(`{
					"offset": 0, "limit": 25, "count": 2, "totalCount": 2,
					"data": [
						{"type":"IPV4_ADDRESSES","id":"c487f45e","name":"DNS Resolvers",
						 "items":[{"type":"IP_ADDRESS","value":"192.0.2.4"}]},
						{"type":"PORTS","id":"177b5b16","name":"HTTP(S)",
						 "items":[{"type":"PORT_NUMBER","value":80},{"type":"PORT_NUMBER","value":443}]}
					]
				}`))
				return
			}
			_, _ = w.Write([]byte(`{"type":"PORTS","id":"177b5b16","name":"HTTP(S)",
				"items":[{"type":"PORT_NUMBER","value":80}]}`))
			return
		case http.MethodPost:
			var in TrafficMatchingList
			_ = json.NewDecoder(r.Body).Decode(&in)
			in.ID = "new-id"
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(in)
			return
		case http.MethodPut:
			// The real API rejects an "id" in the body (it is addressed by the URL
			// path): api.request.unknown-property: Unknown request body property '$.id'.
			var raw map[string]json.RawMessage
			_ = json.NewDecoder(r.Body).Decode(&raw)
			if _, ok := raw["id"]; ok {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"code":"api.request.unknown-property",` +
					`"message":"Unknown request body property '$.id'"}`))
				return
			}
			body, _ := json.Marshal(raw)
			var in TrafficMatchingList
			_ = json.Unmarshal(body, &in)
			in.ID = strings.TrimPrefix(r.URL.Path, base+"/")
			_ = json.NewEncoder(w).Encode(in)
			return
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	t.Cleanup(srv.Close)
	return srv, last
}

func newTrafficMatchingListTestClient(t *testing.T, srv *httptest.Server) *ApiClient {
	t.Helper()
	c, err := New(
		context.Background(),
		&Config{BaseURL: srv.URL, Username: "admin", Password: "admin"},
	)
	if err != nil {
		t.Fatalf("client init: %v", err)
	}
	return c
}

// TestListTrafficMatchingLists_ParsesPage verifies the offset/limit pagination
// envelope is unwrapped and the polymorphic item Value decodes as a string for
// addresses and a number for ports.
func TestListTrafficMatchingLists_ParsesPage(t *testing.T) {
	const siteID = "00000000-0000-0000-0000-000000000000"
	srv, _ := trafficMatchingListTestServer(t, siteID)
	c := newTrafficMatchingListTestClient(t, srv)

	got, err := c.ListTrafficMatchingLists(context.Background(), siteID)
	if err != nil {
		t.Fatalf("ListTrafficMatchingLists errored: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2 lists", len(got))
	}
	addr := got[0]
	if addr.Type != "IPV4_ADDRESSES" || addr.ID != "c487f45e" || addr.Name != "DNS Resolvers" {
		t.Errorf(
			"addr list = %q/%q/%q, want IPV4_ADDRESSES/c487f45e/DNS Resolvers",
			addr.Type,
			addr.ID,
			addr.Name,
		)
	}
	if len(addr.Items) != 1 || addr.Items[0].Type != "IP_ADDRESS" ||
		addr.Items[0].Value != "192.0.2.4" {
		t.Errorf("addr item = %+v, want {IP_ADDRESS 192.0.2.4}", addr.Items)
	}
	ports := got[1]
	if len(ports.Items) != 2 {
		t.Fatalf("ports items = %d, want 2", len(ports.Items))
	}
	// Numeric JSON values decode as float64 through any.
	if ports.Items[0].Type != "PORT_NUMBER" || fmt.Sprint(ports.Items[0].Value) != "80" {
		t.Errorf("port item = %+v, want {PORT_NUMBER 80}", ports.Items[0])
	}
}

// TestGetTrafficMatchingList_Routes verifies a single-resource GET hits the
// {id} path and returns the bare object (no pagination envelope).
func TestGetTrafficMatchingList_Routes(t *testing.T) {
	const siteID = "00000000-0000-0000-0000-000000000000"
	srv, last := trafficMatchingListTestServer(t, siteID)
	c := newTrafficMatchingListTestClient(t, srv)

	got, err := c.GetTrafficMatchingList(context.Background(), siteID, "177b5b16")
	if err != nil {
		t.Fatalf("GetTrafficMatchingList errored: %v", err)
	}
	if got.ID != "177b5b16" || got.Name != "HTTP(S)" {
		t.Errorf("got = %q/%q, want 177b5b16/HTTP(S)", got.ID, got.Name)
	}
	wantPath := "/proxy/network/integration/v1/sites/" + siteID + "/traffic-matching-lists/177b5b16"
	if last.Method != http.MethodGet || last.Path != wantPath {
		t.Errorf("request = %s %s, want GET %s", last.Method, last.Path, wantPath)
	}
}

// TestCreateTrafficMatchingList_ClearsIDAndPosts verifies Create clears any
// caller-supplied ID, POSTs to the collection path, and returns the server id.
func TestCreateTrafficMatchingList_ClearsIDAndPosts(t *testing.T) {
	const siteID = "00000000-0000-0000-0000-000000000000"
	srv, last := trafficMatchingListTestServer(t, siteID)
	c := newTrafficMatchingListTestClient(t, srv)

	got, err := c.CreateTrafficMatchingList(context.Background(), siteID, &TrafficMatchingList{
		ID:   "should-be-cleared",
		Type: "PORTS",
		Name: "DNS",
		Items: []TrafficMatchingListItem{
			{Type: "PORT_NUMBER", Value: 53},
		},
	})
	if err != nil {
		t.Fatalf("CreateTrafficMatchingList errored: %v", err)
	}
	if got.ID != "new-id" || got.Name != "DNS" {
		t.Errorf("got = %q/%q, want new-id/DNS", got.ID, got.Name)
	}
	wantPath := "/proxy/network/integration/v1/sites/" + siteID + "/traffic-matching-lists"
	if last.Method != http.MethodPost || last.Path != wantPath {
		t.Errorf("request = %s %s, want POST %s", last.Method, last.Path, wantPath)
	}
}

// TestUpdateTrafficMatchingList_Routes verifies Update PUTs to the {id} path.
// The test server rejects a body containing "id" (as the real API does), so a
// pass also proves the id is only sent in the URL.
func TestUpdateTrafficMatchingList_Routes(t *testing.T) {
	const siteID = "00000000-0000-0000-0000-000000000000"
	srv, last := trafficMatchingListTestServer(t, siteID)
	c := newTrafficMatchingListTestClient(t, srv)

	in := &TrafficMatchingList{
		ID:   "177b5b16",
		Type: "PORTS",
		Name: "HTTP(S) updated",
	}
	got, err := c.UpdateTrafficMatchingList(context.Background(), siteID, in)
	if err != nil {
		t.Fatalf("UpdateTrafficMatchingList errored: %v", err)
	}
	if got.Name != "HTTP(S) updated" {
		t.Errorf("got name = %q, want HTTP(S) updated", got.Name)
	}
	if got.ID != "177b5b16" {
		t.Errorf("got id = %q, want 177b5b16 from the response", got.ID)
	}
	if in.ID != "177b5b16" {
		t.Errorf(
			"caller's ID = %q, want it left intact (Update must not mutate its argument)",
			in.ID,
		)
	}
	wantPath := "/proxy/network/integration/v1/sites/" + siteID + "/traffic-matching-lists/177b5b16"
	if last.Method != http.MethodPut || last.Path != wantPath {
		t.Errorf("request = %s %s, want PUT %s", last.Method, last.Path, wantPath)
	}
}

// TestDeleteTrafficMatchingList_Routes verifies Delete issues a DELETE to the
// {id} path.
func TestDeleteTrafficMatchingList_Routes(t *testing.T) {
	const siteID = "00000000-0000-0000-0000-000000000000"
	srv, last := trafficMatchingListTestServer(t, siteID)
	c := newTrafficMatchingListTestClient(t, srv)

	if err := c.DeleteTrafficMatchingList(context.Background(), siteID, "177b5b16"); err != nil {
		t.Fatalf("DeleteTrafficMatchingList errored: %v", err)
	}
	wantPath := "/proxy/network/integration/v1/sites/" + siteID + "/traffic-matching-lists/177b5b16"
	if last.Method != http.MethodDelete || last.Path != wantPath {
		t.Errorf("request = %s %s, want DELETE %s", last.Method, last.Path, wantPath)
	}
}
