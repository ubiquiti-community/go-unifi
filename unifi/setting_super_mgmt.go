// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSuperMgmt(ctx context.Context, site string) (*SettingSuperMgmt, error) {
	return c.getSettingSuperMgmt(ctx, site)
}

func (c *Client) UpdateSettingSuperMgmt(
	ctx context.Context,
	site string,
	d *SettingSuperMgmt,
) (*SettingSuperMgmt, error) {
	return c.updateSettingSuperMgmt(ctx, site, d)
}
