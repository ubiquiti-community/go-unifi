// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingDpi(ctx context.Context, site string) (*SettingDpi, error) {
	return c.getSettingDpi(ctx, site)
}

func (c *Client) UpdateSettingDpi(ctx context.Context, site string, d *SettingDpi) (*SettingDpi, error) {
	return c.updateSettingDpi(ctx, site, d)
}
