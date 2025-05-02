// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingBaresip(ctx context.Context, site string) (*SettingBaresip, error) {
	return c.getSettingBaresip(ctx, site)
}

func (c *Client) UpdateSettingBaresip(ctx context.Context, site string, d *SettingBaresip) (*SettingBaresip, error) {
	return c.updateSettingBaresip(ctx, site, d)
}
