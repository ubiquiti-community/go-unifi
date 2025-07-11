// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSslInspection(
	ctx context.Context,
	site string,
) (*SettingSslInspection, error) {
	return c.getSettingSslInspection(ctx, site)
}

func (c *Client) UpdateSettingSslInspection(
	ctx context.Context,
	site string,
	d *SettingSslInspection,
) (*SettingSslInspection, error) {
	return c.updateSettingSslInspection(ctx, site, d)
}
