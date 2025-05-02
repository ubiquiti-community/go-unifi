// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingConnectivity(ctx context.Context, site string) (*SettingConnectivity, error) {
	return c.getSettingConnectivity(ctx, site)
}

func (c *Client) UpdateSettingConnectivity(ctx context.Context, site string, d *SettingConnectivity) (*SettingConnectivity, error) {
	return c.updateSettingConnectivity(ctx, site, d)
}
