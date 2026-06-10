package settings

import (
	"encoding/json"
	"testing"
)

// TestIgmpSnoopingRoundTrip checks the site-level igmp_snooping setting
// (un)marshals correctly, using a payload shaped like a real UniFi 10.x
// controller response. Guards ubiquiti-community/terraform-provider-unifi#164.
func TestIgmpSnoopingRoundTrip(t *testing.T) {
	raw := `{
		"_id": "69d1908dd5c33da485ee2ea2",
		"site_id": "681268bd01e36a7836e2153f",
		"key": "igmp_snooping",
		"enabled": true,
		"flood_known_protocols": true,
		"forward_unknown_mcast_router_ports": false,
		"subscription_mode": "ALL",
		"querier_mode": "CUSTOM",
		"querier_subscription_mode": "ALL",
		"querier_switches": ["d8:b3:70:11:a9:5c"],
		"network_ids": ["681268c001e36a7836e21559", "6813e64a4ee8cb0f1f486ac8"]
	}`

	var s IgmpSnooping
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.GetKey() != "igmp_snooping" {
		t.Errorf("GetKey() = %q, want igmp_snooping", s.GetKey())
	}
	if !s.Enabled {
		t.Error("Enabled = false, want true")
	}
	if len(s.NetworkIDs) != 2 || s.NetworkIDs[0] != "681268c001e36a7836e21559" {
		t.Errorf("NetworkIDs = %v", s.NetworkIDs)
	}
	if s.SubscriptionMode != "ALL" || s.QuerierMode != "CUSTOM" {
		t.Errorf("subscription_mode=%q querier_mode=%q", s.SubscriptionMode, s.QuerierMode)
	}

	// GetSettingKey must resolve the type to the correct endpoint key.
	if k, err := GetSettingKey(&s); err != nil || k != "igmp_snooping" {
		t.Errorf("GetSettingKey = (%q, %v), want (igmp_snooping, nil)", k, err)
	}

	// Re-marshal and ensure key + enabled survive.
	b, err := json.Marshal(&s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back map[string]any
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
	if back["key"] != "igmp_snooping" || back["enabled"] != true {
		t.Errorf("round-trip lost fields: key=%v enabled=%v", back["key"], back["enabled"])
	}
}
