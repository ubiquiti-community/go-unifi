package unifi

import (
	"encoding/json"
	"strings"
	"testing"
)

// Regression for terraform-provider-unifi zone CREATE 500: an empty NetworkIDs
// must serialize as "network_ids":[] (the controller 500s when the field is
// omitted). Before dropping omitempty this body had no network_ids key.
func TestFirewallZoneMarshalIncludesEmptyNetworkIDs(t *testing.T) {
	b, err := json.Marshal(&FirewallZone{Name: "tf-canary", NetworkIDs: []string{}})
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !strings.Contains(got, `"network_ids":[]`) {
		t.Fatalf("expected network_ids to be present as []; got %s", got)
	}
	if strings.Contains(got, "default_zone") {
		t.Fatalf("default_zone should be omitted on write; got %s", got)
	}
}
