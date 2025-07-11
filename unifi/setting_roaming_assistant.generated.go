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

type SettingRoamingAssistant struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	Enabled bool `json:"enabled"`
	Rssi    int  `json:"rssi,omitempty"` // ^-([6-7][0-9]|80)$
}

func (dst *SettingRoamingAssistant) UnmarshalJSON(b []byte) error {
	type Alias SettingRoamingAssistant
	aux := &struct {
		Rssi emptyStringInt `json:"rssi"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Rssi = int(aux.Rssi)

	return nil
}

func (c *Client) getSettingRoamingAssistant(ctx context.Context, site string) (*SettingRoamingAssistant, error) {
	var respBody struct {
		Meta meta                      `json:"meta"`
		Data []SettingRoamingAssistant `json:"data"`
	}
	err := c.do(ctx, "GET", fmt.Sprintf("api/s/%s/get/setting/roaming_assistant", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingRoamingAssistant(ctx context.Context, site string, d *SettingRoamingAssistant) (*SettingRoamingAssistant, error) {
	var respBody struct {
		Meta meta                      `json:"meta"`
		Data []SettingRoamingAssistant `json:"data"`
	}

	d.Key = "roaming_assistant"
	err := c.do(ctx, "PUT", fmt.Sprintf("api/s/%s/set/setting/roaming_assistant", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}
