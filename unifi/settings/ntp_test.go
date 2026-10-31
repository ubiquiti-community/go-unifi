package settings

import (
	"encoding/json"
	"testing"
)

func strPtr(s string) *string { return &s }

// TestNtpServersNullable checks the NTP server fields distinguish "leave
// unchanged" (nil, omitted from the body) from "clear" (pointer to "", sent
// as an explicit empty string).
func TestNtpServersNullable(t *testing.T) {
	s := Ntp{
		NtpServer1:        strPtr("time.cloudflare.com"),
		NtpServer2:        strPtr(""),
		SettingPreference: "manual",
	}

	b, err := json.Marshal(&s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal map: %v", err)
	}
	if got := m["ntp_server_1"]; got != "time.cloudflare.com" {
		t.Errorf("ntp_server_1 = %v, want time.cloudflare.com", got)
	}
	if got, ok := m["ntp_server_2"]; !ok || got != "" {
		t.Errorf("ntp_server_2 = %v (present=%v), want explicit empty string", got, ok)
	}
	for _, k := range []string{"ntp_server_3", "ntp_server_4"} {
		if _, ok := m[k]; ok {
			t.Errorf("%s present in %s, want omitted when nil", k, b)
		}
	}
}

func TestNtpServersUnmarshal(t *testing.T) {
	raw := `{"key":"ntp","ntp_server_1":"pool.ntp.org","ntp_server_2":"","setting_preference":"manual"}`

	var s Ntp
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.NtpServer1 == nil || *s.NtpServer1 != "pool.ntp.org" {
		t.Errorf("NtpServer1 = %v, want pool.ntp.org", s.NtpServer1)
	}
	if s.NtpServer2 == nil || *s.NtpServer2 != "" {
		t.Errorf("NtpServer2 = %v, want pointer to empty string", s.NtpServer2)
	}
	if s.NtpServer3 != nil || s.NtpServer4 != nil {
		t.Errorf("NtpServer3/4 = %v/%v, want nil when absent", s.NtpServer3, s.NtpServer4)
	}
	if s.GetKey() != "ntp" {
		t.Errorf("GetKey() = %q, want ntp", s.GetKey())
	}
}
