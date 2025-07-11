// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingGlobalSwitch(
	ctx context.Context,
	site string,
) (*SettingGlobalSwitch, error) {
	return c.getSettingGlobalSwitch(ctx, site)
}

func (c *Client) UpdateSettingGlobalSwitch(
	ctx context.Context,
	site string,
	d *SettingGlobalSwitch,
) (*SettingGlobalSwitch, error) {
	return c.updateSettingGlobalSwitch(ctx, site, d)
}
