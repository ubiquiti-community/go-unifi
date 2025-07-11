// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSuperSdn(ctx context.Context, site string) (*SettingSuperSdn, error) {
	return c.getSettingSuperSdn(ctx, site)
}

func (c *Client) UpdateSettingSuperSdn(
	ctx context.Context,
	site string,
	d *SettingSuperSdn,
) (*SettingSuperSdn, error) {
	return c.updateSettingSuperSdn(ctx, site, d)
}
