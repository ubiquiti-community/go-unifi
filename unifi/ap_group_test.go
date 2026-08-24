package unifi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// apGroupTestServer stands up a new-style (UniFi OS) controller that serves the
// AP group collection at v2/.../apgroups and, crucially, answers a per-id GET
// the way a real controller does: HTTP 405 with an HTML body. A GetAPGroup that
// hits the per-id path therefore fails with `invalid character '<'`, so these
// tests lock GetAPGroup to the list-and-filter read.
func apGroupTestServer(t *testing.T, site, listJSON string) *httptest.Server {
	t.Helper()
	listPath := "/proxy/network/v2/api/site/" + site + "/apgroups"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == listPath:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(listJSON))
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, listPath+"/"):
			// The controller has no per-id GET for apgroups: 405 + HTML page.
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("<html><body>405 Not Allowed</body></html>"))
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newAPGroupTestClient(t *testing.T, srv *httptest.Server) *ApiClient {
	t.Helper()
	c, err := New(context.Background(), &Config{BaseURL: srv.URL, Username: "admin", Password: "admin"})
	if err != nil {
		t.Fatalf("client init: %v", err)
	}
	return c
}

// TestGetAPGroup_ReadsViaList proves GetAPGroup resolves a group by listing the
// collection and filtering by ID. If it ever regresses to a per-id GET it will
// hit the 405/HTML path and error with `invalid character '<'`, failing here.
func TestGetAPGroup_ReadsViaList(t *testing.T) {
	const site = "default"
	srv := apGroupTestServer(t, site, `[
		{"_id":"a1b2c3d4","name":"Building A","device_macs":["aa:bb:cc:dd:ee:ff"]},
		{"_id":"e5f6a7b8","name":"All APs","device_macs":[]}
	]`)
	c := newAPGroupTestClient(t, srv)

	got, err := c.GetAPGroup(context.Background(), site, "a1b2c3d4")
	if err != nil {
		t.Fatalf("GetAPGroup errored (regressed to per-id GET?): %v", err)
	}
	if got.Name != "Building A" {
		t.Errorf("name = %q, want Building A", got.Name)
	}
	if len(got.DeviceMacs) != 1 || got.DeviceMacs[0] != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("device_macs = %v, want [aa:bb:cc:dd:ee:ff]", got.DeviceMacs)
	}
}

// TestGetAPGroup_NotFound verifies a missing ID yields a typed NotFoundError so
// the Terraform resource's Read can drop it from state instead of erroring.
func TestGetAPGroup_NotFound(t *testing.T) {
	const site = "default"
	srv := apGroupTestServer(t, site, `[{"_id":"a1b2c3d4","name":"Building A","device_macs":[]}]`)
	c := newAPGroupTestClient(t, srv)

	_, err := c.GetAPGroup(context.Background(), site, "does-not-exist")
	if !errors.As(err, new(*NotFoundError)) {
		t.Errorf("expected *NotFoundError for missing id, got %v", err)
	}
}
