// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingElementAdopt(ctx context.Context, site string) (*SettingElementAdopt, error) {
	return c.getSettingElementAdopt(ctx, site)
}

func (c *Client) UpdateSettingElementAdopt(ctx context.Context, site string, d *SettingElementAdopt) (*SettingElementAdopt, error) {
	return c.updateSettingElementAdopt(ctx, site, d)
}
