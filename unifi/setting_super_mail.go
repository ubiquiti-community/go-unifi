// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSuperMail(ctx context.Context, site string) (*SettingSuperMail, error) {
	return c.getSettingSuperMail(ctx, site)
}

func (c *Client) UpdateSettingSuperMail(ctx context.Context, site string, d *SettingSuperMail) (*SettingSuperMail, error) {
	return c.updateSettingSuperMail(ctx, site, d)
}
