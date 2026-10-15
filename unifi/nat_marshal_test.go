package unifi

import (
	"encoding/json"
	"testing"
)

func TestNatFilter_MarshalOmitsEmptyAddressAndSendsGroupIDs(t *testing.T) {
	out, err := json.Marshal(&Nat{
		SourceFilter:      &NatSourceFilter{FilterType: "NONE"},
		DestinationFilter: &NatDestinationFilter{FilterType: "NONE"},
	})
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]json.RawMessage
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"source_filter", "destination_filter"} {
		var f map[string]json.RawMessage
		if err := json.Unmarshal(got[name], &f); err != nil {
			t.Fatal(err)
		}
		if string(f["firewall_group_ids"]) != "[]" {
			t.Errorf("%s firewall_group_ids = %s, want []", name, f["firewall_group_ids"])
		}
		if _, ok := f["address"]; ok {
			t.Errorf("%s address = %s, want it omitted when empty", name, f["address"])
		}
		if string(f["filter_type"]) != `"NONE"` {
			t.Errorf("%s filter_type = %s, want NONE", name, f["filter_type"])
		}
	}
}

func TestNatFilter_MarshalRoundTrips(t *testing.T) {
	port := int64(8443)
	in := &Nat{
		DestinationFilter: &NatDestinationFilter{
			FilterType:       "ADDRESS_AND_PORT",
			Address:          "203.0.113.0/24",
			Port:             &port,
			InvertAddress:    true,
			FirewallGroupIDs: []string{"g1"},
		},
	}
	out, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	var back Nat
	if err := json.Unmarshal(out, &back); err != nil {
		t.Fatal(err)
	}
	got := back.DestinationFilter
	if got == nil || got.FilterType != "ADDRESS_AND_PORT" || got.Address != "203.0.113.0/24" ||
		got.Port == nil || *got.Port != 8443 || !got.InvertAddress ||
		len(got.FirewallGroupIDs) != 1 || got.FirewallGroupIDs[0] != "g1" {
		t.Errorf("round trip = %+v", got)
	}
}
