package unifi_test

import (
	"encoding/json"
	"testing"

	"github.com/ubiquiti-community/go-unifi/unifi"
)

func TestDeviceRadioTableUnmarshalJSONChannel(t *testing.T) {
	for name, tc := range map[string]struct {
		json    string
		channel string
	}{
		"numeric channel":        {json: `{ "channel": 1 }`, channel: "1"},
		"string numeric channel": {json: `{ "channel": "36" }`, channel: "36"},
		"auto channel":           {json: `{ "channel": "auto" }`, channel: "auto"},
		"empty channel":          {json: `{ "channel": "" }`, channel: ""},
	} {
		t.Run(name, func(t *testing.T) {
			var actual unifi.DeviceRadioTable
			if err := json.Unmarshal([]byte(tc.json), &actual); err != nil {
				t.Fatal(err)
			}

			if actual.Channel != tc.channel {
				t.Fatalf("channel = %q, want %q", actual.Channel, tc.channel)
			}
		})
	}
}
