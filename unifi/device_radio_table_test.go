package unifi

import (
	"encoding/json"
	"testing"
)

// Regression test for ubiquiti-community/terraform-provider-unifi#112:
// UniFi 10.x controllers return radio_table channel/tx_power as JSON numbers,
// while older ones (and the schema) use strings such as "auto". The
// DeviceRadioTable unmarshaler must accept either form for these string fields.
func TestDeviceRadioTableUnmarshalChannelTxPower(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		wantChannel string
		wantTxPower string
	}{
		{
			name:        "numeric channel and tx_power (UniFi 10.x)",
			raw:         `{"radio":"na","channel":36,"tx_power":23}`,
			wantChannel: "36",
			wantTxPower: "23",
		},
		{
			name:        "string auto (older controllers)",
			raw:         `{"radio":"ng","channel":"auto","tx_power":"auto"}`,
			wantChannel: "auto",
			wantTxPower: "auto",
		},
		{
			name:        "numeric string channel",
			raw:         `{"radio":"na","channel":"149","tx_power":"high"}`,
			wantChannel: "149",
			wantTxPower: "high",
		},
		{
			name:        "fractional channel (6GHz)",
			raw:         `{"radio":"6e","channel":1.5}`,
			wantChannel: "1.5",
			wantTxPower: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var rt DeviceRadioTable
			if err := json.Unmarshal([]byte(tc.raw), &rt); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if rt.Channel != tc.wantChannel {
				t.Errorf("Channel = %q, want %q", rt.Channel, tc.wantChannel)
			}
			if rt.TxPower != tc.wantTxPower {
				t.Errorf("TxPower = %q, want %q", rt.TxPower, tc.wantTxPower)
			}
		})
	}
}
