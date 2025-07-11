// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSuperFwupdate(
	ctx context.Context,
	site string,
) (*SettingSuperFwupdate, error) {
	return c.getSettingSuperFwupdate(ctx, site)
}

func (c *Client) UpdateSettingSuperFwupdate(
	ctx context.Context,
	site string,
	d *SettingSuperFwupdate,
) (*SettingSuperFwupdate, error) {
	return c.updateSettingSuperFwupdate(ctx, site, d)
}
