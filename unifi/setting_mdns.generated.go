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

type SettingMdns struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	CustomServices     []SettingMdnsCustomServices     `json:"custom_services,omitempty"`
	Mode               string                          `json:"mode,omitempty"` // all|auto|custom
	PredefinedServices []SettingMdnsPredefinedServices `json:"predefined_services,omitempty"`
}

func (dst *SettingMdns) UnmarshalJSON(b []byte) error {
	type Alias SettingMdns
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

type SettingMdnsCustomServices struct {
	Address string `json:"address,omitempty"` // ^_[a-zA-Z0-9._-]+\._(tcp|udp)(\.local)?$
	Name    string `json:"name,omitempty"`
}

func (dst *SettingMdnsCustomServices) UnmarshalJSON(b []byte) error {
	type Alias SettingMdnsCustomServices
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

type SettingMdnsPredefinedServices struct {
	Code string `json:"code,omitempty"` // amazon_devices|android_tv_remote|apple_airDrop|apple_airPlay|apple_file_sharing|apple_iChat|apple_iTunes|aqara|bose|dns_service_discovery|ftp_servers|google_chromecast|homeKit|matter_network|philips_hue|printers|roku|scanners|sonos|spotify_connect|ssh_servers|time_capsule|web_servers|windows_file_sharing_samba
}

func (dst *SettingMdnsPredefinedServices) UnmarshalJSON(b []byte) error {
	type Alias SettingMdnsPredefinedServices
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

func (c *Client) getSettingMdns(ctx context.Context, site string) (*SettingMdns, error) {
	var respBody struct {
		Meta meta          `json:"meta"`
		Data []SettingMdns `json:"data"`
	}
	err := c.do(ctx, "GET", fmt.Sprintf("api/s/%s/get/setting/mdns", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingMdns(ctx context.Context, site string, d *SettingMdns) (*SettingMdns, error) {
	var respBody struct {
		Meta meta          `json:"meta"`
		Data []SettingMdns `json:"data"`
	}

	d.Key = "mdns"
	err := c.do(ctx, "PUT", fmt.Sprintf("api/s/%s/set/setting/mdns", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}
