// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingAutoSpeedtest(ctx context.Context, site string) (*SettingAutoSpeedtest, error) {
	return c.getSettingAutoSpeedtest(ctx, site)
}

func (c *Client) UpdateSettingAutoSpeedtest(ctx context.Context, site string, d *SettingAutoSpeedtest) (*SettingAutoSpeedtest, error) {
	return c.updateSettingAutoSpeedtest(ctx, site, d)
}
