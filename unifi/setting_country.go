// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingCountry(ctx context.Context, site string) (*SettingCountry, error) {
	return c.getSettingCountry(ctx, site)
}

func (c *Client) UpdateSettingCountry(ctx context.Context, site string, d *SettingCountry) (*SettingCountry, error) {
	return c.updateSettingCountry(ctx, site, d)
}
