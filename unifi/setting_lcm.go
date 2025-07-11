// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingLcm(ctx context.Context, site string) (*SettingLcm, error) {
	return c.getSettingLcm(ctx, site)
}

func (c *Client) UpdateSettingLcm(
	ctx context.Context,
	site string,
	d *SettingLcm,
) (*SettingLcm, error) {
	return c.updateSettingLcm(ctx, site, d)
}
