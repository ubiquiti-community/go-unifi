package unifi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubiquiti-community/go-unifi/unifi"
)

func TestPowerSupervisorMarshalJSON(t *testing.T) {
	// Create body shape: power_sources may be sent empty.
	ps := unifi.PowerSupervisor{
		ClientMAC: "00:11:22:33:44:55",
		Enabled:   true,
		Settings: unifi.PowerSupervisorSettings{
			HeartbeatInterval: 60,
			SilenceThreshold:  900,
			PowerOffDuration:  120,
		},
		PowerSources: []unifi.PowerSupervisorSource{},
	}

	actual, err := json.Marshal(&ps)
	if err != nil {
		t.Fatal(err)
	}
	assert.JSONEq(t,
		`{"client_mac":"00:11:22:33:44:55","enabled":true,`+
			`"settings":{"heartbeat_interval":60,"silence_threshold":900,"power_off_duration":120},`+
			`"power_sources":[]}`,
		string(actual))
}

func TestPowerSupervisorUnmarshalJSON(t *testing.T) {
	// Resting shape returned by GET, with the controller-resolved power source.
	raw := `{
		"id": "000000000000000000000001",
		"site_id": "000000000000000000000002",
		"client_mac": "00:11:22:33:44:55",
		"enabled": true,
		"consecutive_failures": 0,
		"settings": {"heartbeat_interval": 60, "silence_threshold": 900, "power_off_duration": 120},
		"power_sources": [
			{"client_psu_index": 1, "power_source_index": 4, "power_source_mac": "66:77:88:99:aa:bb", "power_source_type": "poe_port"}
		]
	}`

	var ps unifi.PowerSupervisor
	if err := json.Unmarshal([]byte(raw), &ps); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "000000000000000000000001", ps.ID)
	assert.Equal(t, "00:11:22:33:44:55", ps.ClientMAC)
	assert.True(t, ps.Enabled)
	assert.Equal(t, 900, ps.Settings.SilenceThreshold)
	if assert.Len(t, ps.PowerSources, 1) {
		assert.Equal(t, "66:77:88:99:aa:bb", ps.PowerSources[0].PowerSourceMAC)
		assert.Equal(t, 4, ps.PowerSources[0].PowerSourceIndex)
		assert.Equal(t, "poe_port", ps.PowerSources[0].PowerSourceType)
	}
}
