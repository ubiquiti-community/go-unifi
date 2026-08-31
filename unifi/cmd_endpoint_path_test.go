package unifi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// cmdPathTestServer stands up a new-style (UniFi OS) controller that records the
// path of the first cmd request it sees and answers it with an empty success
// body. Everything else answers 200 so client setup completes.
func cmdPathTestServer(t *testing.T, got *string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodPost && *got == "" {
			*got = r.URL.Path
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newCmdPathTestClient(t *testing.T, srv *httptest.Server) *ApiClient {
	t.Helper()
	c, err := New(context.Background(), &Config{BaseURL: srv.URL, Username: "admin", Password: "admin"})
	if err != nil {
		t.Fatalf("client init: %v", err)
	}
	return c
}

// The cmd endpoints live under the API root like every other call, so their
// relative paths must carry the `api/` segment themselves: apiPath supplies only
// the deployment prefix, which is "/proxy/network" on UniFi OS and "/" on a
// classic controller.
func TestCmdEndpointsCarryAPIPrefix(t *testing.T) {
	const site = "default"

	tests := []struct {
		name string
		want string
		call func(*ApiClient) error
	}{
		{
			name: "CreateSite",
			want: "/proxy/network/api/s/default/cmd/sitemgr",
			call: func(c *ApiClient) error {
				_, err := c.CreateSite(context.Background(), "a description")
				return err
			},
		},
		{
			name: "DeleteSite",
			want: "/proxy/network/api/s/default/cmd/sitemgr",
			call: func(c *ApiClient) error {
				_, err := c.DeleteSite(context.Background(), "abcdef0123456789abcdef01")
				return err
			},
		},
		{
			name: "UpdateSite",
			want: "/proxy/network/api/s/" + site + "/cmd/sitemgr",
			call: func(c *ApiClient) error {
				_, err := c.UpdateSite(context.Background(), site, "a new description")
				return err
			},
		},
		{
			name: "ReorderFirewallRules",
			want: "/proxy/network/api/s/" + site + "/cmd/firewall",
			call: func(c *ApiClient) error {
				return c.ReorderFirewallRules(context.Background(), site, "WAN_IN", nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got string
			srv := cmdPathTestServer(t, &got)
			c := newCmdPathTestClient(t, srv)

			if err := tc.call(c); err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			if got != tc.want {
				t.Errorf("%s requested %q, want %q", tc.name, got, tc.want)
			}
		})
	}
}
