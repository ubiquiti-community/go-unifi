// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingLocale(ctx context.Context, site string) (*SettingLocale, error) {
	return c.getSettingLocale(ctx, site)
}

func (c *Client) UpdateSettingLocale(ctx context.Context, site string, d *SettingLocale) (*SettingLocale, error) {
	return c.updateSettingLocale(ctx, site, d)
}
