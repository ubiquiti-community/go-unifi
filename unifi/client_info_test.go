package unifi

import (
	"encoding/json"
	"testing"
)

func TestClientInfoDeserialization(t *testing.T) {
	tests := []struct {
		name          string
		raw           string
		wantPort      *int64
		wantUplinkMAC string
		wantDisplay   string
	}{
		{
			name:          "direct object",
			raw:           `{"display_name":"talos","mac":"2c:cf:67:0a:a3:33","last_uplink_mac":"f4:e2:c6:50:60:bb","last_uplink_remote_port":5,"status":"DISCONNECTED"}`,
			wantPort:      ptrInt64(5),
			wantUplinkMAC: "f4:e2:c6:50:60:bb",
			wantDisplay:   "talos",
		},
		{
			name:          "meta/data array wrapper silently fails",
			raw:           `{"meta":{"rc":"ok"},"data":[{"display_name":"talos","mac":"2c:cf:67:0a:a3:33","last_uplink_mac":"f4:e2:c6:50:60:bb","last_uplink_remote_port":5}]}`,
			wantPort:      nil,
			wantUplinkMAC: "",
			wantDisplay:   "",
		},
		{
			name:          "data object wrapper silently fails",
			raw:           `{"data":{"display_name":"talos","mac":"2c:cf:67:0a:a3:33","last_uplink_mac":"f4:e2:c6:50:60:bb","last_uplink_remote_port":5}}`,
			wantPort:      nil,
			wantUplinkMAC: "",
			wantDisplay:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ci ClientInfo
			err := json.Unmarshal([]byte(tt.raw), &ci)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			t.Logf("DisplayName=%q LastUplinkMac=%q LastUplinkRemotePort=%v SwPort=%v Status=%q",
				ci.DisplayName, ci.LastUplinkMac, ci.LastUplinkRemotePort, ci.SwPort, ci.Status)

			if tt.wantPort == nil {
				if ci.LastUplinkRemotePort != nil {
					t.Errorf("expected nil port, got %d", *ci.LastUplinkRemotePort)
				}
			} else {
				if ci.LastUplinkRemotePort == nil {
					t.Errorf("expected port %d, got nil", *tt.wantPort)
				} else if *ci.LastUplinkRemotePort != *tt.wantPort {
					t.Errorf("expected port %d, got %d", *tt.wantPort, *ci.LastUplinkRemotePort)
				}
			}

			if ci.LastUplinkMac != tt.wantUplinkMAC {
				t.Errorf("expected uplink MAC %q, got %q", tt.wantUplinkMAC, ci.LastUplinkMac)
			}
			if ci.DisplayName != tt.wantDisplay {
				t.Errorf("expected display %q, got %q", tt.wantDisplay, ci.DisplayName)
			}
		})
	}
}

// The controller reports channel_width as a bare number for some clients and a
// quoted string for others; a list decode fails wholesale on the mismatch.
func TestClientInfoChannelWidth(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantInt *int64
	}{
		{name: "number", raw: `{"mac":"aa:bb:cc:dd:ee:ff","channel_width":80}`, want: "80", wantInt: ptrInt64(80)},
		{name: "string", raw: `{"mac":"aa:bb:cc:dd:ee:ff","channel_width":"80"}`, want: "80", wantInt: ptrInt64(80)},
		{name: "empty string", raw: `{"mac":"aa:bb:cc:dd:ee:ff","channel_width":""}`, want: ""},
		{name: "null", raw: `{"mac":"aa:bb:cc:dd:ee:ff","channel_width":null}`, want: ""},
		{name: "absent", raw: `{"mac":"aa:bb:cc:dd:ee:ff"}`, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ci ClientInfo
			if err := json.Unmarshal([]byte(tt.raw), &ci); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if ci.ChannelWidth.String() != tt.want {
				t.Errorf("ChannelWidth = %q, want %q", ci.ChannelWidth, tt.want)
			}
			got := ci.ChannelWidth.Int64Pointer()
			switch {
			case tt.wantInt == nil && got != nil:
				t.Errorf("ChannelWidth.Int64Pointer() = %d, want nil", *got)
			case tt.wantInt != nil && got == nil:
				t.Errorf("ChannelWidth.Int64Pointer() = nil, want %d", *tt.wantInt)
			case tt.wantInt != nil && *got != *tt.wantInt:
				t.Errorf("ChannelWidth.Int64Pointer() = %d, want %d", *got, *tt.wantInt)
			}
		})
	}

	t.Run("mixed list", func(t *testing.T) {
		var list []ClientInfo
		raw := `[{"mac":"aa:bb:cc:dd:ee:01","channel_width":"20"},{"mac":"aa:bb:cc:dd:ee:02","channel_width":80}]`
		if err := json.Unmarshal([]byte(raw), &list); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if len(list) != 2 || list[0].ChannelWidth != "20" || list[1].ChannelWidth != "80" {
			t.Errorf("got %+v", list)
		}
	})
}

func ptrInt64(v int64) *int64 { return &v }
