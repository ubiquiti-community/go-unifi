// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
	_ types.Number
	_ strconv.NumError
)

type GlobalSwitch struct {
	BaseSetting

	AclDeviceIsolation             []string                            `json:"acl_device_isolation,omitempty"`
	AclL3Isolation                 []SettingGlobalSwitchAclL3Isolation `json:"acl_l3_isolation,omitempty"`
	AutoStpEdgeDetectionEnabled    bool                                `json:"auto_stp_edge_detection_enabled"`
	DHCPSnoop                      bool                                `json:"dhcp_snoop"`
	Dot1XFallbackNetworkID         string                              `json:"dot1x_fallback_networkconf_id,omitempty"` // [\d\w-]+|
	Dot1XPortctrlEnabled           bool                                `json:"dot1x_portctrl_enabled"`
	FloodKnownProtocols            bool                                `json:"flood_known_protocols"`
	FlowctrlEnabled                bool                                `json:"flowctrl_enabled"`
	ForwardUnknownMcastRouterPorts bool                                `json:"forward_unknown_mcast_router_ports"`
	JumboframeEnabled              bool                                `json:"jumboframe_enabled"`
	LinkDebounce                   *int64                              `json:"link_debounce,omitempty"`          // 0|[1-9]00|[1-4][0-9]00|5000
	PoeStagingDelayMsec            *int64                              `json:"poe_staging_delay_msec,omitempty"` // 0|200|400|600|800|1000|1200|1400|1600|1800|2000
	RADIUSProfileID                string                              `json:"radiusprofile_id,omitempty"`
	StpVersion                     string                              `json:"stp_version,omitempty"`       // stp|rstp|disabled
	SwitchExclusions               []string                            `json:"switch_exclusions,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
}

func (dst *GlobalSwitch) UnmarshalJSON(b []byte) error {
	type Alias GlobalSwitch
	aux := &struct {
		LinkDebounce        *types.Number `json:"link_debounce"`
		PoeStagingDelayMsec *types.Number `json:"poe_staging_delay_msec"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	// First unmarshal base setting
	if err := json.Unmarshal(b, &dst.BaseSetting); err != nil {
		return fmt.Errorf("unable to unmarshal base setting: %w", err)
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if aux.LinkDebounce != nil {
		if val, err := aux.LinkDebounce.Int64(); err == nil {
			dst.LinkDebounce = &val
		} else if string(*aux.LinkDebounce) == "" {
			var zero int64
			dst.LinkDebounce = &zero
		}
	}
	if aux.PoeStagingDelayMsec != nil {
		if val, err := aux.PoeStagingDelayMsec.Int64(); err == nil {
			dst.PoeStagingDelayMsec = &val
		} else if string(*aux.PoeStagingDelayMsec) == "" {
			var zero int64
			dst.PoeStagingDelayMsec = &zero
		}
	}

	return nil
}

type SettingGlobalSwitchAclL3Isolation struct {
	DestinationNetworks []string `json:"destination_networks,omitempty"`
	SourceNetwork       string   `json:"source_network,omitempty"`
}

func (dst *SettingGlobalSwitchAclL3Isolation) UnmarshalJSON(b []byte) error {
	type Alias SettingGlobalSwitchAclL3Isolation
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}
