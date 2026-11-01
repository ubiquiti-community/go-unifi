package unifi

import (
	"encoding/json"
	"testing"
)

// DevicePortTable decodes its numeric fields through types.Number and copies
// them only when the conversion succeeds. The check was inverted, so every
// numeric field of port_table stayed 0 and a port could not be looked up by
// port_idx.
func TestDevicePortTableUnmarshalNumbers(t *testing.T) {
	t.Run("numeric fields are decoded", func(t *testing.T) {
		body := `{
			"port_idx": 5, "poe_caps": 7, "speed_caps": 1048623, "anomalies": 1,
			"mac_table_count": 3,
			"rx_broadcast": 11, "rx_bytes": 123456, "rx_dropped": 12, "rx_errors": 13,
			"rx_multicast": 14, "rx_packets": 15,
			"satisfaction": 100, "satisfaction_reason": 2, "speed": 1000,
			"stp_pathcost": 20000,
			"tx_broadcast": 21, "tx_bytes": 654321, "tx_dropped": 22, "tx_errors": 23,
			"tx_multicast": 24, "tx_packets": 25,
			"stormctrl_bcast_rate": 31, "stormctrl_mcast_rate": 32, "stormctrl_ucast_rate": 33,
			"name": "Port 5"
		}`
		var pt DevicePortTable
		if err := json.Unmarshal([]byte(body), &pt); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		for _, tc := range []struct {
			field     string
			got, want int64
		}{
			{"PortIdx", pt.PortIdx, 5},
			{"PoeCaps", pt.PoeCaps, 7},
			{"SpeedCaps", pt.SpeedCaps, 1048623},
			{"Anomalies", pt.Anomalies, 1},
			{"MacTableCount", pt.MacTableCount, 3},
			{"RxBroadcast", pt.RxBroadcast, 11},
			{"RxBytes", pt.RxBytes, 123456},
			{"RxDropped", pt.RxDropped, 12},
			{"RxErrors", pt.RxErrors, 13},
			{"RxMulticast", pt.RxMulticast, 14},
			{"RxPackets", pt.RxPackets, 15},
			{"Satisfaction", pt.Satisfaction, 100},
			{"SatisfactionReason", pt.SatisfactionReason, 2},
			{"Speed", pt.Speed, 1000},
			{"StpPathcost", pt.StpPathcost, 20000},
			{"TxBroadcast", pt.TxBroadcast, 21},
			{"TxBytes", pt.TxBytes, 654321},
			{"TxDropped", pt.TxDropped, 22},
			{"TxErrors", pt.TxErrors, 23},
			{"TxMulticast", pt.TxMulticast, 24},
			{"TxPackets", pt.TxPackets, 25},
			{"StormctrlBcastRate", pt.StormctrlBcastRate, 31},
			{"StormctrlMcastRate", pt.StormctrlMcastRate, 32},
			{"StormctrlUcastRate", pt.StormctrlUcastRate, 33},
		} {
			if tc.got != tc.want {
				t.Errorf("%s = %d, want %d", tc.field, tc.got, tc.want)
			}
		}
		if pt.Name != "Port 5" {
			t.Errorf("Name = %q, want %q", pt.Name, "Port 5")
		}
	})

	t.Run("string-encoded and missing numbers", func(t *testing.T) {
		var pt DevicePortTable
		if err := json.Unmarshal([]byte(`{"port_idx":"8","speed":""}`), &pt); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if pt.PortIdx != 8 {
			t.Errorf("PortIdx = %d, want 8", pt.PortIdx)
		}
		if pt.Speed != 0 || pt.RxBytes != 0 {
			t.Errorf("Speed = %d, RxBytes = %d, want 0 for empty and missing values", pt.Speed, pt.RxBytes)
		}
	})
}
