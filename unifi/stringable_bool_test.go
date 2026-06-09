package unifi

import (
	"encoding/json"
	"testing"
)

// TestNetworkDHCPDEnabledTolerantUnmarshal verifies that Network.DHCPDEnabled
// decodes whether the controller sends a native JSON boolean or a quoted
// string (terraform-provider-unifi #65).
func TestNetworkDHCPDEnabledTolerantUnmarshal(t *testing.T) {
	cases := map[string]bool{
		`{"dhcpd_enabled": true}`:    true,
		`{"dhcpd_enabled": false}`:   false,
		`{"dhcpd_enabled": "true"}`:  true,
		`{"dhcpd_enabled": "false"}`: false,
		`{"dhcpd_enabled": ""}`:      false,
		`{}`:                         false,
	}
	for body, want := range cases {
		var n Network
		if err := json.Unmarshal([]byte(body), &n); err != nil {
			t.Errorf("Unmarshal(%s) error: %v", body, err)
			continue
		}
		if n.DHCPDEnabled != want {
			t.Errorf("Unmarshal(%s): DHCPDEnabled = %v, want %v", body, n.DHCPDEnabled, want)
		}
	}
}

// TestClientBlockedTolerantUnmarshal verifies that Client.Blocked decodes both
// native and string-encoded booleans, and stays nil when absent
// (terraform-provider-unifi #132).
func TestClientBlockedTolerantUnmarshal(t *testing.T) {
	truthy := []string{`{"blocked": true}`, `{"blocked": "true"}`}
	for _, body := range truthy {
		var c Client
		if err := json.Unmarshal([]byte(body), &c); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", body, err)
		}
		if c.Blocked == nil || !*c.Blocked {
			t.Errorf("Unmarshal(%s): Blocked = %v, want true", body, c.Blocked)
		}
	}

	falsy := []string{`{"blocked": false}`, `{"blocked": "false"}`, `{"blocked": ""}`}
	for _, body := range falsy {
		var c Client
		if err := json.Unmarshal([]byte(body), &c); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", body, err)
		}
		if c.Blocked == nil || *c.Blocked {
			t.Errorf("Unmarshal(%s): Blocked = %v, want false", body, c.Blocked)
		}
	}

	var c Client
	if err := json.Unmarshal([]byte(`{}`), &c); err != nil {
		t.Fatalf("Unmarshal({}) error: %v", err)
	}
	if c.Blocked != nil {
		t.Errorf("Unmarshal({}): Blocked = %v, want nil", c.Blocked)
	}
}
