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

type IgmpSnooping struct {
	BaseSetting

	AutoUnknownTrafficHandling         bool                                  `json:"auto_unknown_traffic_handling"`
	Enabled                            bool                                  `json:"enabled"`
	FailoverQuerier                    string                                `json:"failover_querier,omitempty"`
	FastleaveForNetworkIDs             []string                              `json:"fastleave_for_network_ids,omitempty"`
	FloodKnownProtocols                bool                                  `json:"flood_known_protocols"`
	FloodUnknownMulticastForNetworkIDs []string                              `json:"flood_unknown_multicast_for_network_ids,omitempty"`
	ForwardUnknownMcastRouterPorts     bool                                  `json:"forward_unknown_mcast_router_ports"`
	NetworkIDs                         []string                              `json:"network_ids,omitempty"`
	PrimaryQuerier                     string                                `json:"primary_querier,omitempty"`
	QuerierAddresses                   []SettingIgmpSnoopingQuerierAddresses `json:"querier_addresses,omitempty"`
	QuerierMode                        string                                `json:"querier_mode,omitempty"`              // PRIMARY_AND_FAILOVER|CUSTOM|OFF
	QuerierSubscriptionMode            string                                `json:"querier_subscription_mode,omitempty"` // ALL|CUSTOM
	QuerierSwitches                    []string                              `json:"querier_switches,omitempty"`
	SubscriptionMode                   string                                `json:"subscription_mode,omitempty"` // ALL|CUSTOM
	Switches                           []string                              `json:"switches,omitempty"`
}

func (dst *IgmpSnooping) UnmarshalJSON(b []byte) error {
	type Alias IgmpSnooping
	aux := &struct {
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

	return nil
}

type SettingIgmpSnoopingQuerierAddresses struct {
	MAC            string `json:"mac,omitempty"`
	NetworkID      string `json:"network_id,omitempty"`
	QuerierAddress string `json:"querier_address,omitempty"`
	QueryInterval  *int64 `json:"query_interval,omitempty"` // [3-9][0-9]|1[0-7][0-9]|180
}

func (dst *SettingIgmpSnoopingQuerierAddresses) UnmarshalJSON(b []byte) error {
	type Alias SettingIgmpSnoopingQuerierAddresses
	aux := &struct {
		QueryInterval *types.Number `json:"query_interval"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if aux.QueryInterval != nil {
		if val, err := aux.QueryInterval.Int64(); err == nil {
			dst.QueryInterval = &val
		} else if string(*aux.QueryInterval) == "" {
			var zero int64
			dst.QueryInterval = &zero
		}
	}

	return nil
}
