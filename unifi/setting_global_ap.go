// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingGlobalAp(ctx context.Context, site string) (*SettingGlobalAp, error) {
	return c.getSettingGlobalAp(ctx, site)
}

func (c *Client) UpdateSettingGlobalAp(
	ctx context.Context,
	site string,
	d *SettingGlobalAp,
) (*SettingGlobalAp, error) {
	return c.updateSettingGlobalAp(ctx, site, d)
}
