// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingIps(ctx context.Context, site string) (*SettingIps, error) {
	return c.getSettingIps(ctx, site)
}

func (c *Client) UpdateSettingIps(
	ctx context.Context,
	site string,
	d *SettingIps,
) (*SettingIps, error) {
	return c.updateSettingIps(ctx, site, d)
}
