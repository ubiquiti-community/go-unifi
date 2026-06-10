package unifi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFirewallPolicyPortMarshalsAsString(t *testing.T) {
	cases := []struct {
		name     string
		port     *int64
		wantPort string // exact JSON fragment expected, or "" when port must be absent
	}{
		{name: "specific port", port: ptrInt64(161), wantPort: `"port":"161"`},
		{name: "no port", port: nil, wantPort: ""},
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

// Round-trips a string-encoded port back through decode to confirm the tolerant
// UnmarshalJSON still reads what we now write.
func TestFirewallPolicyPortRoundTrip(t *testing.T) {
	src := FirewallPolicySource{Port: ptrInt64(443)}
	b, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got FirewallPolicySource
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Port == nil || *got.Port != 443 {
		t.Fatalf("round-trip port mismatch: got %v, want 443", got.Port)
	}
}
