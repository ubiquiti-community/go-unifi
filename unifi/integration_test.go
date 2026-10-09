package unifi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

const integrationSitesPath = "/proxy/network/integration/v1/sites"

// integrationSitesTestServer serves the Integration API site listing. The
// handler receives the parsed offset/limit so tests can drive pagination, and
// the returned counter reports how many listing requests were made.
func integrationSitesTestServer(
	t *testing.T,
	handler func(w http.ResponseWriter, offset, limit int),
) (*httptest.Server, *int32) {
	t.Helper()
	var hits int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodGet || r.URL.Path != integrationSitesPath {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		atomic.AddInt32(&hits, 1)
		var offset, limit int
		_, _ = fmt.Sscanf(r.URL.Query().Get("offset"), "%d", &offset)
		_, _ = fmt.Sscanf(r.URL.Query().Get("limit"), "%d", &limit)
		w.Header().Set("Content-Type", "application/json")
		handler(w, offset, limit)
	}))
	t.Cleanup(srv.Close)
	return srv, &hits
}

const twoSitesBody = `{
	"offset": 0, "limit": 200, "count": 2, "totalCount": 2,
	"data": [
		{"id":"11111111-1111-1111-1111-111111111111","internalReference":"default","name":"Default"},
		{"id":"22222222-2222-2222-2222-222222222222","internalReference":"branch","name":"Branch Office"}
	]
}`

func TestResolveIntegrationSiteID_UUIDPassesThroughWithoutRequest(t *testing.T) {
	srv, hits := integrationSitesTestServer(t, func(w http.ResponseWriter, _, _ int) {
		_, _ = w.Write([]byte(twoSitesBody))
	})
	c := newTrafficMatchingListTestClient(t, srv)

	const uuid = "AABBCCDD-0000-1111-2222-333344445555"
	got, err := c.ResolveIntegrationSiteID(context.Background(), uuid)
	if err != nil {
		t.Fatalf("ResolveIntegrationSiteID errored: %v", err)
	}
	if got != IntegrationSiteID(uuid) {
		t.Errorf("got %q, want %q", got, uuid)
	}
	if n := atomic.LoadInt32(hits); n != 0 {
		t.Errorf("listing requests = %d, want 0 for a UUID input", n)
	}
}

func TestResolveIntegrationSiteID_ByLegacyName(t *testing.T) {
	srv, _ := integrationSitesTestServer(t, func(w http.ResponseWriter, _, _ int) {
		_, _ = w.Write([]byte(twoSitesBody))
	})
	c := newTrafficMatchingListTestClient(t, srv)

	got, err := c.ResolveIntegrationSiteID(context.Background(), "default")
	if err != nil {
		t.Fatalf("ResolveIntegrationSiteID errored: %v", err)
	}
	if want := IntegrationSiteID("11111111-1111-1111-1111-111111111111"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolveIntegrationSiteID_ByDisplayName(t *testing.T) {
	srv, _ := integrationSitesTestServer(t, func(w http.ResponseWriter, _, _ int) {
		_, _ = w.Write([]byte(twoSitesBody))
	})
	c := newTrafficMatchingListTestClient(t, srv)

	got, err := c.ResolveIntegrationSiteID(context.Background(), "Branch Office")
	if err != nil {
		t.Fatalf("ResolveIntegrationSiteID errored: %v", err)
	}
	if want := IntegrationSiteID("22222222-2222-2222-2222-222222222222"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolveIntegrationSiteID_LegacyNameWinsOverDisplayName(t *testing.T) {
	// One site's display name equals another's internalReference; the
	// internalReference match must win.
	srv, _ := integrationSitesTestServer(t, func(w http.ResponseWriter, _, _ int) {
		_, _ = w.Write([]byte(`{
			"offset": 0, "limit": 200, "count": 2, "totalCount": 2,
			"data": [
				{"id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","internalReference":"x","name":"default"},
				{"id":"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb","internalReference":"default","name":"Other"}
			]
		}`))
	})
	c := newTrafficMatchingListTestClient(t, srv)

	got, err := c.ResolveIntegrationSiteID(context.Background(), "default")
	if err != nil {
		t.Fatalf("ResolveIntegrationSiteID errored: %v", err)
	}
	if want := IntegrationSiteID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"); got != want {
		t.Errorf("got %q, want %q (internalReference must win)", got, want)
	}
}

func TestResolveIntegrationSiteID_NotFound(t *testing.T) {
	srv, _ := integrationSitesTestServer(t, func(w http.ResponseWriter, _, _ int) {
		_, _ = w.Write([]byte(twoSitesBody))
	})
	c := newTrafficMatchingListTestClient(t, srv)

	_, err := c.ResolveIntegrationSiteID(context.Background(), "nope")
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("err = %v, want *NotFoundError", err)
	}
	if nf.Type != "IntegrationSite" || nf.Value != "nope" {
		t.Errorf("NotFoundError = %+v, want IntegrationSite/nope", nf)
	}
}

func TestResolveIntegrationSiteID_AmbiguousDisplayName(t *testing.T) {
	srv, _ := integrationSitesTestServer(t, func(w http.ResponseWriter, _, _ int) {
		_, _ = w.Write([]byte(`{
			"offset": 0, "limit": 200, "count": 2, "totalCount": 2,
			"data": [
				{"id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","internalReference":"a","name":"Office"},
				{"id":"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb","internalReference":"b","name":"Office"}
			]
		}`))
	})
	c := newTrafficMatchingListTestClient(t, srv)

	_, err := c.ResolveIntegrationSiteID(context.Background(), "Office")
	if !errors.Is(err, ErrAmbiguousIntegrationSite) {
		t.Fatalf("err = %v, want ErrAmbiguousIntegrationSite", err)
	}
	if !strings.Contains(err.Error(), "ambiguous") {
		t.Errorf("err = %v, want the message to mention ambiguity", err)
	}
	var nf *NotFoundError
	if errors.As(err, &nf) {
		t.Errorf("ambiguity must not be reported as NotFoundError: %v", err)
	}
}

// TestResolveIntegrationSiteID_RequestFailureIsNeitherNotFoundNorAmbiguous guards
// against reporting a failed request (DNS, server error) as a site problem: the
// caller must be able to tell it apart from a site that simply does not exist.
func TestResolveIntegrationSiteID_RequestFailureIsNeitherNotFoundNorAmbiguous(t *testing.T) {
	srv, _ := integrationSitesTestServer(t, func(w http.ResponseWriter, _, _ int) {
		// A 4xx is not retried by the client, unlike a 5xx, which keeps this fast.
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"code":"api.forbidden","message":"nope"}`))
	})
	c := newTrafficMatchingListTestClient(t, srv)

	_, err := c.ResolveIntegrationSiteID(context.Background(), "default")
	if err == nil {
		t.Fatal("expected an error when the site listing request fails")
	}
	var nf *NotFoundError
	if errors.As(err, &nf) {
		t.Errorf("a failed request must not be reported as NotFoundError: %v", err)
	}
	if errors.Is(err, ErrAmbiguousIntegrationSite) {
		t.Errorf("a failed request must not be reported as ambiguous: %v", err)
	}
}

func TestListIntegrationSites_FollowsPagination(t *testing.T) {
	// Serve three sites two at a time; the client must request offset 0 then 2.
	all := []string{
		`{"id":"00000000-0000-0000-0000-000000000001","internalReference":"s1","name":"S1"}`,
		`{"id":"00000000-0000-0000-0000-000000000002","internalReference":"s2","name":"S2"}`,
		`{"id":"00000000-0000-0000-0000-000000000003","internalReference":"s3","name":"S3"}`,
	}
	srv, hits := integrationSitesTestServer(t, func(w http.ResponseWriter, offset, _ int) {
		end := offset + 2
		if end > len(all) {
			end = len(all)
		}
		page := all[offset:end]
		_, _ = fmt.Fprintf(
			w,
			`{"offset":%d,"limit":2,"count":%d,"totalCount":%d,"data":[%s]}`,
			offset, len(page), len(all), strings.Join(page, ","),
		)
	})
	c := newTrafficMatchingListTestClient(t, srv)

	got, err := c.ListIntegrationSites(context.Background())
	if err != nil {
		t.Fatalf("ListIntegrationSites errored: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3 sites across pages", len(got))
	}
	if got[2].InternalReference != "s3" {
		t.Errorf("last site = %+v, want s3", got[2])
	}
	if n := atomic.LoadInt32(hits); n != 2 {
		t.Errorf("listing requests = %d, want 2 (offset 0 then 2)", n)
	}
}

func TestListIntegrationSites_EmptyStopsWithoutLooping(t *testing.T) {
	srv, hits := integrationSitesTestServer(t, func(w http.ResponseWriter, _, _ int) {
		// totalCount lies (claims more); an empty page must still terminate.
		_, _ = w.Write([]byte(`{"offset":0,"limit":200,"count":0,"totalCount":5,"data":[]}`))
	})
	c := newTrafficMatchingListTestClient(t, srv)

	got, err := c.ListIntegrationSites(context.Background())
	if err != nil {
		t.Fatalf("ListIntegrationSites errored: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
	if n := atomic.LoadInt32(hits); n != 1 {
		t.Errorf("listing requests = %d, want 1 (no infinite loop)", n)
	}
}
