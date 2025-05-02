// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingDoh(ctx context.Context, site string) (*SettingDoh, error) {
	return c.getSettingDoh(ctx, site)
}

func (c *Client) UpdateSettingDoh(ctx context.Context, site string, d *SettingDoh) (*SettingDoh, error) {
	return c.updateSettingDoh(ctx, site, d)
}
