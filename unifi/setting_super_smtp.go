// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSuperSmtp(ctx context.Context, site string) (*SettingSuperSmtp, error) {
	return c.getSettingSuperSmtp(ctx, site)
}

func (c *Client) UpdateSettingSuperSmtp(
	ctx context.Context,
	site string,
	d *SettingSuperSmtp,
) (*SettingSuperSmtp, error) {
	return c.updateSettingSuperSmtp(ctx, site, d)
}
