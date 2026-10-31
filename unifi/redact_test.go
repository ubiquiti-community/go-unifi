package unifi

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestRedactSensitivePayload guards that secrets embedded in a request body are
// not leaked when the body is included in an error message
// (ubiquiti-community/terraform-provider-unifi#256).
func TestRedactSensitivePayload(t *testing.T) {
	const (
		wgKey     = "fake-wireguard-key-value"
		pass      = "fake-passphrase-value"
		ipsecPSK  = "fake-ipsec-psk-value"
		radiusSec = "fake-radius-secret-value"
		pw        = "fake-password-value"
	)
	body := []byte(`{
		"name": "wgadmin",
		"purpose": "vpn-server",
		"x_wireguard_private_key": "` + wgKey + `",
		"x_passphrase": "` + pass + `",
		"x_ipsec_pre_shared_key": "` + ipsecPSK + `",
		"vlan": 50,
		"nested": {"radius_secret": "` + radiusSec + `", "ok": "keep"},
		"list": [{"password": "` + pw + `"}]
	}`)

	out := redactSensitivePayload(body)

	for _, leak := range []string{wgKey, pass, ipsecPSK, radiusSec, pw} {
		if strings.Contains(out, leak) {
			t.Errorf("redacted payload still leaks %q: %s", leak, out)
		}
	}

	// Non-sensitive fields must survive.
	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("redacted payload is not valid JSON: %v\n%s", err, out)
	}
	if m["name"] != "wgadmin" {
		t.Errorf("name was lost: %v", m["name"])
	}
	if m["x_wireguard_private_key"] != "REDACTED" {
		t.Errorf("private key not redacted: %v", m["x_wireguard_private_key"])
	}
	if nested, ok := m["nested"].(map[string]any); !ok || nested["ok"] != "keep" {
		t.Errorf("nested non-sensitive value lost: %v", m["nested"])
	}

	// Non-JSON body is omitted, not echoed.
	if got := redactSensitivePayload([]byte("not json")); strings.Contains(got, "not json") {
		t.Errorf("non-JSON body echoed: %q", got)
	}
}

// UniFi stores secrets in fields with the hidden-field prefix x_, and in token
// and preshared_key fields, under names the older substrings do not match.
func TestRedactSensitivePayloadKeys(t *testing.T) {
	tests := []struct {
		key    string
		value  string
		redact bool
	}{
		{key: "x_api_token", value: `"fake-api-token"`, redact: true},
		{key: "x_iapp_key", value: `"fake-iapp-key"`, redact: true},
		{key: "x_wep", value: `"fake-wep-key"`, redact: true},
		{key: "x_mgmt_key", value: `"fake-mgmt-key"`, redact: true},
		{key: "X_Shared_Client_Key", value: `"fake-client-key"`, redact: true},
		{key: "auth_token", value: `"fake-auth-token"`, redact: true},
		{key: "wireguard_client_preshared_key", value: `"fake-preshared-key"`, redact: true},
		{key: "x_ssh_enabled", value: `true`, redact: false},
		{key: "dot1x_ctrl", value: `"auto"`, redact: false},
		{key: "tx_power_mode", value: `"auto"`, redact: false},
		{key: "name", value: `"guest"`, redact: false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			out := redactSensitivePayload([]byte(`{"` + tt.key + `":` + tt.value + `}`))

			var m map[string]json.RawMessage
			if err := json.Unmarshal([]byte(out), &m); err != nil {
				t.Fatalf("redacted payload is not valid JSON: %v\n%s", err, out)
			}
			want := tt.value
			if tt.redact {
				want = `"REDACTED"`
			}
			if got := string(m[tt.key]); got != want {
				t.Errorf("%s = %s, want %s", tt.key, got, want)
			}
		})
	}
}
