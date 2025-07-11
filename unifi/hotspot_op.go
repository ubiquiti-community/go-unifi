// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListHotspotOp(ctx context.Context, site string) ([]HotspotOp, error) {
	return c.listHotspotOp(ctx, site)
}

func (c *Client) GetHotspotOp(ctx context.Context, site, id string) (*HotspotOp, error) {
	return c.getHotspotOp(ctx, site, id)
}

func (c *Client) DeleteHotspotOp(ctx context.Context, site, id string) error {
	return c.deleteHotspotOp(ctx, site, id)
}

func (c *Client) CreateHotspotOp(
	ctx context.Context,
	site string,
	d *HotspotOp,
) (*HotspotOp, error) {
	return c.createHotspotOp(ctx, site, d)
}

func (c *Client) UpdateHotspotOp(
	ctx context.Context,
	site string,
	d *HotspotOp,
) (*HotspotOp, error) {
	return c.updateHotspotOp(ctx, site, d)
}
