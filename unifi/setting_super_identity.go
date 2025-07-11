// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSuperIdentity(
	ctx context.Context,
	site string,
) (*SettingSuperIdentity, error) {
	return c.getSettingSuperIdentity(ctx, site)
}

func (c *Client) UpdateSettingSuperIdentity(
	ctx context.Context,
	site string,
	d *SettingSuperIdentity,
) (*SettingSuperIdentity, error) {
	return c.updateSettingSuperIdentity(ctx, site, d)
}
