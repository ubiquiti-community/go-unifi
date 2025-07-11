// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingBroadcast(ctx context.Context, site string) (*SettingBroadcast, error) {
	return c.getSettingBroadcast(ctx, site)
}

func (c *Client) UpdateSettingBroadcast(
	ctx context.Context,
	site string,
	d *SettingBroadcast,
) (*SettingBroadcast, error) {
	return c.updateSettingBroadcast(ctx, site, d)
}
