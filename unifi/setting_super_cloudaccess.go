// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSuperCloudaccess(
	ctx context.Context,
	site string,
) (*SettingSuperCloudaccess, error) {
	return c.getSettingSuperCloudaccess(ctx, site)
}

func (c *Client) UpdateSettingSuperCloudaccess(
	ctx context.Context,
	site string,
	d *SettingSuperCloudaccess,
) (*SettingSuperCloudaccess, error) {
	return c.updateSettingSuperCloudaccess(ctx, site, d)
}
