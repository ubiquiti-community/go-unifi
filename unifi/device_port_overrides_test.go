package unifi

import (
	"encoding/json"
	"strings"
	"testing"
)

// terraform-provider-unifi#567: port_overrides is a full-replace array on the
// controller PUT, so a key the body leaves out resets to the controller
// default (for autoneg: enabled). With a plain bool and omitempty an explicit
// false was dropped; the booleans are *bool so false stays on the wire, and a
// decoded false survives a round-trip while an unset key stays omitted.
func TestDevicePortOverridesBoolPointers(t *testing.T) {
	t.Run("explicit false is serialized", func(t *testing.T) {
		body, err := json.Marshal(DevicePortOverrides{Autoneg: boolPtr(false)})
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if !strings.Contains(string(body), `"autoneg":false`) {
			t.Errorf("explicit false was dropped\nbody: %s", body)
		}
	})

	t.Run("decoded false survives a round-trip", func(t *testing.T) {
		var po DevicePortOverrides
		if err := json.Unmarshal([]byte(`{"port_idx":49,"autoneg":false}`), &po); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if po.Autoneg == nil || *po.Autoneg {
			t.Fatalf("Autoneg = %v, want pointer to false", po.Autoneg)
		}
		body, err := json.Marshal(po)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if !strings.Contains(string(body), `"autoneg":false`) {
			t.Errorf("decoded false was dropped on re-marshal\nbody: %s", body)
		}
	})

	t.Run("unset stays omitted", func(t *testing.T) {
		var po DevicePortOverrides
		if err := json.Unmarshal([]byte(`{"port_idx":1}`), &po); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if po.Autoneg != nil {
			t.Fatalf("Autoneg = %v, want nil", po.Autoneg)
		}
		body, err := json.Marshal(po)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if strings.Contains(string(body), `"autoneg"`) {
			t.Errorf("unset autoneg was serialized\nbody: %s", body)
		}
	})
}
