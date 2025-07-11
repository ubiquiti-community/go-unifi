// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingGuestAccess(
	ctx context.Context,
	site string,
) (*SettingGuestAccess, error) {
	return c.getSettingGuestAccess(ctx, site)
}

func (c *Client) UpdateSettingGuestAccess(
	ctx context.Context,
	site string,
	d *SettingGuestAccess,
) (*SettingGuestAccess, error) {
	return c.updateSettingGuestAccess(ctx, site, d)
}
