package unifi

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestFirewallPolicySourceMACBinding guards terraform-provider-unifi #575: the
// controller only keeps match_mac on a source with matching_target IP when the
// bound MAC is sent as `mac`; without it the API answers
// api.err.MissingFirewallPolicySourceIpMac.
func TestFirewallPolicySourceMACBinding(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		b, err := json.Marshal(FirewallPolicySource{
			ZoneID:         "zone1",
			MatchingTarget: "IP",
			IPs:            []string{"192.168.1.10"},
			MatchMAC:       true,
			MAC:            "aa:bb:cc:dd:ee:ff",
		})
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		s := string(b)
		for _, want := range []string{`"match_mac":true`, `"mac":"aa:bb:cc:dd:ee:ff"`} {
			if !strings.Contains(s, want) {
				t.Errorf("expected %s in %s", want, s)
			}
		}
	})

	t.Run("omitted when empty", func(t *testing.T) {
		b, err := json.Marshal(FirewallPolicySource{ZoneID: "zone1", MatchingTarget: "IP"})
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if strings.Contains(string(b), `"mac"`) {
			t.Errorf("mac must be omitted when empty: %s", b)
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		var src FirewallPolicySource
		in := `{"zone_id":"zone1","matching_target":"IP","ips":["192.168.1.10"],"match_mac":true,"mac":"aa:bb:cc:dd:ee:ff"}`
		if err := json.Unmarshal([]byte(in), &src); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if !src.MatchMAC || src.MAC != "aa:bb:cc:dd:ee:ff" {
			t.Errorf("got match_mac=%v mac=%q", src.MatchMAC, src.MAC)
		}
	})
}
