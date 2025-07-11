// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type SettingTrafficFlow struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	EnabledAllowedTraffic        bool `json:"enabled_allowed_traffic"`
	GatewayDNSEnabled            bool `json:"gateway_dns_enabled"`
	UnifiDeviceManagementEnabled bool `json:"unifi_device_management_enabled"`
	UnifiServicesEnabled         bool `json:"unifi_services_enabled"`
}

func (dst *SettingTrafficFlow) UnmarshalJSON(b []byte) error {
	type Alias SettingTrafficFlow
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

func (c *Client) getSettingTrafficFlow(ctx context.Context, site string) (*SettingTrafficFlow, error) {
	var respBody struct {
		Meta meta                 `json:"meta"`
		Data []SettingTrafficFlow `json:"data"`
	}
	err := c.do(ctx, "GET", fmt.Sprintf("api/s/%s/get/setting/traffic_flow", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingTrafficFlow(ctx context.Context, site string, d *SettingTrafficFlow) (*SettingTrafficFlow, error) {
	var respBody struct {
		Meta meta                 `json:"meta"`
		Data []SettingTrafficFlow `json:"data"`
	}

	d.Key = "traffic_flow"
	err := c.do(ctx, "PUT", fmt.Sprintf("api/s/%s/set/setting/traffic_flow", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}
