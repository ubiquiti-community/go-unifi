package unifi_test

import (
	"encoding/json"
	"testing"

	"github.com/ubiquiti-community/go-unifi/unifi"
)

// TestDevicePortTable_AggregatedBy verifies that the aggregated_by field
// handles polymorphic JSON values from the UniFi API.
//
// The API returns aggregated_by as:
//   - false (boolean): port not in LAG, or is LAG master
//   - integer: LAG member port, value is the master port number
func TestDevicePortTable_AggregatedBy(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantVal interface{}
	}{
		{"bool false", `{"aggregated_by": false}`, false},
		{"bool true", `{"aggregated_by": true}`, true},
		{"integer", `{"aggregated_by": 15}`, float64(15)},
		{"zero", `{"aggregated_by": 0}`, float64(0)},
		{"omitted", `{}`, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dpt unifi.DevicePortTable
			if err := json.Unmarshal([]byte(tt.input), &dpt); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			assertAggregatedBy(t, dpt.AggregatedBy, tt.wantVal)
		})
	}
}

// TestDevicePortTable_AggregatedBy_LAGScenarios tests realistic LAG configurations.
func TestDevicePortTable_AggregatedBy_LAGScenarios(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantVal interface{}
	}{
		{
			name: "LAG member port",
			input: `{
				"media": "SFP+",
				"op_mode": "aggregate",
				"aggregated_by": 26
			}`,
			wantVal: float64(26),
		},
		{
			name: "LAG master port",
			input: `{
				"media": "SFP+",
				"op_mode": "aggregate",
				"aggregated_by": false
			}`,
			wantVal: false,
		},
		{
			name: "regular port",
			input: `{
				"media": "GE",
				"op_mode": "switch",
				"aggregated_by": false
			}`,
			wantVal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dpt unifi.DevicePortTable
			if err := json.Unmarshal([]byte(tt.input), &dpt); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			assertAggregatedBy(t, dpt.AggregatedBy, tt.wantVal)
		})
	}
}

func assertAggregatedBy(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("AggregatedBy = %v (%T), want %v (%T)", got, got, want, want)
	}
}
