// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingGlobalNat(ctx context.Context, site string) (*SettingGlobalNat, error) {
	return c.getSettingGlobalNat(ctx, site)
}

func (c *Client) UpdateSettingGlobalNat(ctx context.Context, site string, d *SettingGlobalNat) (*SettingGlobalNat, error) {
	return c.updateSettingGlobalNat(ctx, site, d)
}
