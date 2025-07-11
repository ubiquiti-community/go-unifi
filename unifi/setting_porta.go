// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingPorta(ctx context.Context, site string) (*SettingPorta, error) {
	return c.getSettingPorta(ctx, site)
}

func (c *Client) UpdateSettingPorta(
	ctx context.Context,
	site string,
	d *SettingPorta,
) (*SettingPorta, error) {
	return c.updateSettingPorta(ctx, site, d)
}
