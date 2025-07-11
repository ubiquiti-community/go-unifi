// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListHotspotPackage(ctx context.Context, site string) ([]HotspotPackage, error) {
	return c.listHotspotPackage(ctx, site)
}

func (c *Client) GetHotspotPackage(ctx context.Context, site, id string) (*HotspotPackage, error) {
	return c.getHotspotPackage(ctx, site, id)
}

func (c *Client) DeleteHotspotPackage(ctx context.Context, site, id string) error {
	return c.deleteHotspotPackage(ctx, site, id)
}

func (c *Client) CreateHotspotPackage(
	ctx context.Context,
	site string,
	d *HotspotPackage,
) (*HotspotPackage, error) {
	return c.createHotspotPackage(ctx, site, d)
}

func (c *Client) UpdateHotspotPackage(
	ctx context.Context,
	site string,
	d *HotspotPackage,
) (*HotspotPackage, error) {
	return c.updateHotspotPackage(ctx, site, d)
}
