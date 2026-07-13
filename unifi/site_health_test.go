package unifi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// siteHealthTestServer stands up a new-style (UniFi OS) controller that serves
// the site health metrics at s/{site}/stat/health. When rawBody is non-empty
// it is served verbatim with the given status; otherwise dataJSON is wrapped
// in an rc:ok envelope.
func siteHealthTestServer(t *testing.T, site, dataJSON, rawBody string, status int) *httptest.Server {
	t.Helper()
	healthPath := "/proxy/network/api/s/" + site + "/stat/health"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == healthPath {
			w.Header().Set("Content-Type", "application/json")
			if status != http.StatusOK {
				w.WriteHeader(status)
			}
			if rawBody != "" {
				_, _ = w.Write([]byte(rawBody))
				return
			}
			if status != http.StatusOK {
				_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"api.err.NoSiteContext"},"data":[]}`))
				return
			}
			_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":` + dataJSON + `}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newSiteHealthTestClient(t *testing.T, srv *httptest.Server) *ApiClient {
	t.Helper()
	c, err := New(context.Background(), &Config{BaseURL: srv.URL, Username: "admin", Password: "admin"})
	if err != nil {
		t.Fatalf("client init: %v", err)
	}
	return c
}

// TestGetHealth_ParsesSubsystems verifies GetHealth returns one entry per
// subsystem with the per-subsystem fields decoded (WAN gateway details, WLAN
// device/user counts).
func TestGetHealth_ParsesSubsystems(t *testing.T) {
	const site = "default"
	srv := siteHealthTestServer(t, site, `[
		{"subsystem":"wan","status":"ok","wan_ip":"203.0.113.10","gw_mac":"aa:bb:cc:dd:ee:ff","gw_version":"4.4.56","latency":12,"uptime":86400,"gateways":["203.0.113.1"]},
		{"subsystem":"wlan","status":"ok","num_ap":3,"num_adopted":3,"num_user":17,"num_guest":2},
		{"subsystem":"vpn","status":"unknown"}
	]`, "", http.StatusOK)
	c := newSiteHealthTestClient(t, srv)

	got, err := c.GetHealth(context.Background(), site)
	if err != nil {
		t.Fatalf("GetHealth errored: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3 subsystems", len(got))
	}
	wan := got[0]
	if wan.Subsystem != "wan" || wan.Status != "ok" {
		t.Errorf("wan subsystem/status = %q/%q, want wan/ok", wan.Subsystem, wan.Status)
	}
	if wan.WanIP != "203.0.113.10" || wan.GwMac != "aa:bb:cc:dd:ee:ff" || fmt.Sprint(wan.Latency) != "12" {
		t.Errorf("wan details = %q/%q/%v, want 203.0.113.10/aa:bb:cc:dd:ee:ff/12", wan.WanIP, wan.GwMac, wan.Latency)
	}
	if len(wan.Gateways) != 1 || wan.Gateways[0] != "203.0.113.1" {
		t.Errorf("wan gateways = %v, want [203.0.113.1]", wan.Gateways)
	}
	wlan := got[1]
	if fmt.Sprint(wlan.NumAp) != "3" || fmt.Sprint(wlan.NumUser) != "17" || fmt.Sprint(wlan.NumGuest) != "2" {
		t.Errorf("wlan counts = %v/%v/%v, want 3/17/2", wlan.NumAp, wlan.NumUser, wlan.NumGuest)
	}
}

// TestGetHealth_Error verifies a controller error status surfaces as an error
// rather than an empty result.
func TestGetHealth_Error(t *testing.T) {
	const site = "default"
	srv := siteHealthTestServer(t, site, "", "", http.StatusBadRequest)
	c := newSiteHealthTestClient(t, srv)

	_, err := c.GetHealth(context.Background(), site)
	if err == nil {
		t.Fatal("expected error from rc:error response, got nil")
	}
}

// TestGetHealth_MixedNumberRepresentations verifies that numeric health fields
// decode whether the controller emits them as JSON numbers or as strings —
// real controllers emit both representations across builds.
func TestGetHealth_MixedNumberRepresentations(t *testing.T) {
	const site = "default"
	srv := siteHealthTestServer(t, site, `[
		{"subsystem":"wan","status":"ok","latency":"12","xput_up":92.5,"tx_bytes-r":123456,
		 "gw_system-stats":{"cpu":4.2,"mem":"38.9","uptime":"86400"}},
		{"subsystem":"wlan","status":"ok","num_user":"17","num_ap":3}
	]`, "", http.StatusOK)
	c := newSiteHealthTestClient(t, srv)

	got, err := c.GetHealth(context.Background(), site)
	if err != nil {
		t.Fatalf("GetHealth errored on mixed number representations: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2 subsystems", len(got))
	}
	wan := got[0]
	if fmt.Sprint(wan.Latency) != "12" || fmt.Sprint(wan.XputUp) != "92.5" || fmt.Sprint(wan.TxBytesR) != "123456" {
		t.Errorf("wan numerics = %v/%v/%v, want 12/92.5/123456", wan.Latency, wan.XputUp, wan.TxBytesR)
	}
	stats := wan.GwSystemStats
	if fmt.Sprint(stats.CPU) != "4.2" || fmt.Sprint(stats.Mem) != "38.9" || fmt.Sprint(stats.Uptime) != "86400" {
		t.Errorf("gw system stats = %v/%v/%v, want 4.2/38.9/86400", stats.CPU, stats.Mem, stats.Uptime)
	}
	wlan := got[1]
	if fmt.Sprint(wlan.NumUser) != "17" || fmt.Sprint(wlan.NumAp) != "3" {
		t.Errorf("wlan counts = %v/%v, want 17/3", wlan.NumUser, wlan.NumAp)
	}
}

// TestGetHealth_MetaErrorWithHTTP200 verifies that an in-band controller error
// (HTTP 200 with meta.rc "error") is surfaced instead of an empty success.
func TestGetHealth_MetaErrorWithHTTP200(t *testing.T) {
	const site = "default"
	srv := siteHealthTestServer(t, site, "",
		`{"meta":{"rc":"error","msg":"api.err.NoSiteContext"},"data":[]}`, http.StatusOK)
	c := newSiteHealthTestClient(t, srv)

	_, err := c.GetHealth(context.Background(), site)
	if err == nil {
		t.Fatal("expected error from HTTP 200 rc:error response, got nil")
	}
	if !strings.Contains(err.Error(), "api.err.NoSiteContext") {
		t.Errorf("error = %q, want it to carry the controller message", err)
	}
}
