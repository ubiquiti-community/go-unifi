package settings

import (
	"encoding/json"
	"testing"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
)

// TestMgmtSSHKeyDate checks an SSH key's date decodes whether the controller
// sends it as a number (UniFi Network 11.0.104) or as a string, and is sent
// back in the same form.
func TestMgmtSSHKeyDate(t *testing.T) {
	cases := []struct {
		name     string
		date     string
		wantDate types.NumberOrString
	}{
		{"number", `1732000000`, "1732000000"},
		{"string", `"2024-01-31T12:00:00Z"`, "2024-01-31T12:00:00Z"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := `{"key":"mgmt","auto_upgrade":true,"auto_upgrade_hour":3,` +
				`"x_ssh_enabled":true,"x_ssh_keys":[{"name":"admin","type":"ssh-ed25519",` +
				`"key":"AAAAC3NzaC1lZDI1NTE5AAAAIExample","comment":"admin@example",` +
				`"date":` + tc.date + `,"fingerprint":"SHA256:example"}]}`

			var s Mgmt
			if err := json.Unmarshal([]byte(raw), &s); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if s.GetKey() != "mgmt" {
				t.Errorf("GetKey() = %q, want mgmt", s.GetKey())
			}
			if !s.AutoUpgrade || s.AutoUpgradeHour == nil || *s.AutoUpgradeHour != 3 {
				t.Errorf("AutoUpgrade = %v, AutoUpgradeHour = %v, want true, 3", s.AutoUpgrade, s.AutoUpgradeHour)
			}
			if len(s.SSHKeys) != 1 {
				t.Fatalf("SSHKeys = %v, want one key", s.SSHKeys)
			}
			k := s.SSHKeys[0]
			if k.Date != tc.wantDate {
				t.Errorf("Date = %q, want %q", k.Date, tc.wantDate)
			}
			if k.Name != "admin" || k.KeyType != "ssh-ed25519" || k.Fingerprint != "SHA256:example" {
				t.Errorf("SSH key = %+v, want the other fields decoded", k)
			}

			b, err := json.Marshal(&s)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var m struct {
				SSHKeys []map[string]json.RawMessage `json:"x_ssh_keys"`
			}
			if err := json.Unmarshal(b, &m); err != nil {
				t.Fatalf("unmarshal marshalled setting: %v", err)
			}
			if len(m.SSHKeys) != 1 || string(m.SSHKeys[0]["date"]) != tc.date {
				t.Errorf("marshalled x_ssh_keys = %s, want date %s", b, tc.date)
			}
		})
	}
}
