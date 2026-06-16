package unifi

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestFirewallPolicyPortMarshalsAsString guards terraform-provider-unifi #288
// and #286: the zone-based firewall API expects the source/destination `port`
// as a quoted string (e.g. "161" or a comma-separated list "80,443"), and a
// portless endpoint (port_matching_type ANY) must omit the field entirely — a
// `"0"` makes the gateway reject the whole firewall config.
func TestFirewallPolicyPortMarshalsAsString(t *testing.T) {
	cases := []struct {
		name     string
		port     string
		wantPort string // exact JSON fragment expected, or "" when port must be absent
	}{
		{name: "specific port", port: "161", wantPort: `"port":"161"`},
		{name: "comma separated", port: "80,443", wantPort: `"port":"80,443"`},
		{name: "no port", port: "", wantPort: ""},
	}

	for _, tc := range cases {
		t.Run("source/"+tc.name, func(t *testing.T) {
			b, err := json.Marshal(FirewallPolicySource{
				ZoneID:           "zone1",
				MatchingTarget:   "IP",
				PortMatchingType: "SPECIFIC",
				Port:             tc.port,
			})
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			assertPort(t, string(b), tc.wantPort)
		})

		t.Run("destination/"+tc.name, func(t *testing.T) {
			b, err := json.Marshal(FirewallPolicyDestination{
				ZoneID:           "zone2",
				MatchingTarget:   "IP",
				PortMatchingType: "SPECIFIC",
				Port:             tc.port,
			})
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			assertPort(t, string(b), tc.wantPort)
		})
	}
}

func assertPort(t *testing.T, got, wantPort string) {
	t.Helper()
	// The numeric form must never appear on the wire.
	if strings.Contains(got, `"port":161`) {
		t.Fatalf("port encoded as a number, firmware would reject it: %s", got)
	}
	if wantPort == "" {
		if strings.Contains(got, `"port"`) {
			t.Fatalf("expected port to be absent, got: %s", got)
		}
		return
	}
	if !strings.Contains(got, wantPort) {
		t.Fatalf("expected %s in payload, got: %s", wantPort, got)
	}
}

// TestFirewallPolicyPortUnmarshal confirms the tolerant decoder accepts a port
// sent as a JSON number, a quoted string, a comma-separated list, or omitted.
func TestFirewallPolicyPortUnmarshal(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
	}{
		{name: "numeric", body: `{"port":443}`, want: "443"},
		{name: "string", body: `{"port":"443"}`, want: "443"},
		{name: "comma separated", body: `{"port":"1812,1813"}`, want: "1812,1813"},
		{name: "empty string", body: `{"port":""}`, want: ""},
		{name: "absent", body: `{}`, want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got FirewallPolicySource
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.Port != tc.want {
				t.Fatalf("Port = %q, want %q", got.Port, tc.want)
			}
		})
	}
}
