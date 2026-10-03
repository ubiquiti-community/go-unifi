package unifi

import (
	"encoding/json"
	"strings"
	"testing"
)

// terraform-provider-unifi#542: config_network is a full replace on the
// controller - a key the PUT body leaves out is deleted from the stored object
// rather than left untouched. With omitempty on the fields whose ordinary
// value is the zero value (bonding_enabled=false, an empty dns2), the body
// omits them, the controller drops them, and the provider then reads them back
// as null and fails the apply with "inconsistent result after apply".
//
// Measured on a UCG-Fiber running Network 10.6.106: a config_network write
// without bonding_enabled and dns2 removed both from the stored object, and
// re-sending them explicitly restored it.
func TestDeviceConfigNetworkSerializesZeroValues(t *testing.T) {
	body, err := json.Marshal(&DeviceConfigNetwork{
		Type:    "static",
		IP:      "10.0.0.2",
		Netmask: "255.255.255.0",
		// gateway, dns1, dns2, dnssuffix and bonding_enabled left at their
		// zero values, as a configuration that does not set them produces.
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	for _, key := range []string{
		"bonding_enabled",
		"dns1",
		"dns2",
		"dnssuffix",
		"gateway",
	} {
		if !strings.Contains(string(body), `"`+key+`"`) {
			t.Errorf(
				"%s is missing from the body: omitempty on it makes the controller "+
					"delete the stored value (terraform-provider-unifi#542)\nbody: %s",
				key, body,
			)
		}
	}
}

// ip, netmask and type deliberately keep omitempty: the controller rejects a
// static config_network carrying an empty one, and the provider never builds
// the object with one missing. This pins that decision so a future sweep does
// not widen the exception without measuring it.
func TestDeviceConfigNetworkOmitsRequiredFieldsWhenEmpty(t *testing.T) {
	body, err := json.Marshal(&DeviceConfigNetwork{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	for _, key := range []string{"ip", "netmask", "type"} {
		if strings.Contains(string(body), `"`+key+`"`) {
			t.Errorf("%s should stay omitempty\nbody: %s", key, body)
		}
	}
}
