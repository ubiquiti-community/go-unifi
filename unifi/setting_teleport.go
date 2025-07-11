// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingTeleport(ctx context.Context, site string) (*SettingTeleport, error) {
	return c.getSettingTeleport(ctx, site)
}

func (c *Client) UpdateSettingTeleport(
	ctx context.Context,
	site string,
	d *SettingTeleport,
) (*SettingTeleport, error) {
	return c.updateSettingTeleport(ctx, site, d)
}
