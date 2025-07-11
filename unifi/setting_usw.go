// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingUsw(ctx context.Context, site string) (*SettingUsw, error) {
	return c.getSettingUsw(ctx, site)
}

func (c *Client) UpdateSettingUsw(
	ctx context.Context,
	site string,
	d *SettingUsw,
) (*SettingUsw, error) {
	return c.updateSettingUsw(ctx, site, d)
}
