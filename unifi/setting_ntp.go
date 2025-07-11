// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingNtp(ctx context.Context, site string) (*SettingNtp, error) {
	return c.getSettingNtp(ctx, site)
}

func (c *Client) UpdateSettingNtp(
	ctx context.Context,
	site string,
	d *SettingNtp,
) (*SettingNtp, error) {
	return c.updateSettingNtp(ctx, site, d)
}
