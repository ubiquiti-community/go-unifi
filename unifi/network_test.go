package unifi_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ubiquiti-community/go-unifi/unifi"
)

func TestNetworkUnmarshalJSON(t *testing.T) {
	for n, c := range map[string]struct {
		expected func(n *unifi.Network)
		json     string
	}{
		"int64 vlan": {
			expected: func(n *unifi.Network) { n.VLAN = unifi.Ptr[int64](1) },
			json:     `{ "vlan": 1 }`,
		},
		"string vlan": {
			expected: func(n *unifi.Network) { n.VLAN = unifi.Ptr[int64](1) },
			json:     `{ "vlan": "1" }`,
		},
		// vlan's validation accepts "" but not 0, so "" means unset.
		"empty string vlan": {
			expected: func(n *unifi.Network) {},
			json:     `{ "vlan": "" }`,
		},

		"int64 dhcpd_leasetime": {
			expected: func(n *unifi.Network) { n.DHCPDLeaseTime = unifi.Ptr[int64](1) },
			json:     `{ "dhcpd_leasetime": 1 }`,
		},
		"string dhcpd_leasetime": {
			expected: func(n *unifi.Network) { n.DHCPDLeaseTime = unifi.Ptr[int64](1) },
			json:     `{ "dhcpd_leasetime": "1" }`,
		},
		"empty string dhcpd_leasetime": {
			expected: func(n *unifi.Network) { n.DHCPDLeaseTime = unifi.Ptr[int64](0) },
			json:     `{ "dhcpd_leasetime": "" }`,
		},

		"int64 wan_egress_qos": {
			expected: func(n *unifi.Network) { n.WANEgressQOS = unifi.Ptr[int64](1) },
			json:     `{ "wan_egress_qos": 1 }`,
		},
		"string wan_egress_qos": {
			expected: func(n *unifi.Network) { n.WANEgressQOS = unifi.Ptr[int64](1) },
			json:     `{ "wan_egress_qos": "1" }`,
		},
		// wan_egress_qos is [1-7]|^$: "" means unset and must not become 0,
		// which the controller rejects (#105).
		"empty string wan_egress_qos": {
			expected: func(n *unifi.Network) {},
			json:     `{ "wan_egress_qos": "" }`,
		},

		"int64 wan_vlan": {
			expected: func(n *unifi.Network) { n.WANVLAN = unifi.Ptr[int64](1) },
			json:     `{ "wan_vlan": 1 }`,
		},
		"string wan_vlan": {
			expected: func(n *unifi.Network) { n.WANVLAN = unifi.Ptr[int64](1) },
			json:     `{ "wan_vlan": "1" }`,
		},
		// wan_vlan's validation accepts 0, so "" keeps decoding to 0.
		"empty wan_vlan vlan": {
			expected: func(n *unifi.Network) { n.WANVLAN = unifi.Ptr[int64](0) },
			json:     `{ "wan_vlan": "" }`,
		},
	} {
		t.Run(n, func(t *testing.T) {
			// set some non-zero value defaults
			expected := unifi.Network{
				InternetAccessEnabled: true,
			}
			c.expected(&expected)
			var actual unifi.Network
			err := json.Unmarshal(([]byte)(c.json), &actual)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("not equal:\nexpected: %#v\nactual: %#v", expected, actual)
			}
		})
	}
}

// TestNetworkEmptyEgressQOSRoundTrip covers #105: a WAN stored with
// wan_egress_qos "" must not be re-sent as the invalid value 0.
func TestNetworkEmptyEgressQOSRoundTrip(t *testing.T) {
	stored := `{
		"_id": "1",
		"purpose": "wan",
		"name": "WAN",
		"wan_type": "pppoe",
		"wan_egress_qos": "",
		"wan_egress_qos_enabled": false,
		"wan_prefixlen": "",
		"wan_load_balance_weight": "",
		"wan_dhcpv6_pd_size": "",
		"wan_vlan": ""
	}`
	var n unifi.Network
	if err := json.Unmarshal([]byte(stored), &n); err != nil {
		t.Fatal(err)
	}
	for name, v := range map[string]*int64{
		"wan_egress_qos":          n.WANEgressQOS,
		"wan_prefixlen":           n.WANPrefixlen,
		"wan_load_balance_weight": n.WANLoadBalanceWeight,
		"wan_dhcpv6_pd_size":      n.WANDHCPv6PDSize,
	} {
		if v != nil {
			t.Errorf("%s: want nil, got %d", name, *v)
		}
	}

	b, err := json.Marshal(&n)
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal(b, &body); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"wan_egress_qos", "wan_prefixlen", "wan_load_balance_weight", "wan_dhcpv6_pd_size"} {
		if v, ok := body[k]; ok {
			t.Errorf("%s: want omitted, got %v", k, v)
		}
	}

	// A configured value still round-trips.
	n.WANEgressQOS = unifi.Ptr[int64](5)
	b, err = json.Marshal(&n)
	if err != nil {
		t.Fatal(err)
	}
	body = nil
	if err := json.Unmarshal(b, &body); err != nil {
		t.Fatal(err)
	}
	if got := body["wan_egress_qos"]; got != float64(5) {
		t.Errorf("wan_egress_qos: want 5, got %v", got)
	}
}

// TestDiffEmptyEgressQOSIsNoChange checks that an unset wan_egress_qos
// decoded from "" does not show up as a change in a sparse PUT.
func TestDiffEmptyEgressQOSIsNoChange(t *testing.T) {
	var prior unifi.Network
	if err := json.Unmarshal([]byte(`{"_id":"1","purpose":"wan","name":"WAN","wan_egress_qos":""}`), &prior); err != nil {
		t.Fatal(err)
	}
	target := prior
	target.Name = unifi.Ptr("WAN renamed")
	patch, err := unifi.Diff(&prior, &target, unifi.DiffOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := patch["wan_egress_qos"]; ok {
		t.Errorf("wan_egress_qos must not be in the patch: %v", patch)
	}
	if patch["name"] != "WAN renamed" {
		t.Errorf("name: want change, got %v", patch)
	}
}
