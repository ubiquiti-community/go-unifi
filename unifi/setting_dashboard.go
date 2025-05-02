// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingDashboard(ctx context.Context, site string) (*SettingDashboard, error) {
	return c.getSettingDashboard(ctx, site)
}

func (c *Client) UpdateSettingDashboard(ctx context.Context, site string, d *SettingDashboard) (*SettingDashboard, error) {
	return c.updateSettingDashboard(ctx, site, d)
}
