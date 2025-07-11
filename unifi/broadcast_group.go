// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListBroadcastGroup(ctx context.Context, site string) ([]BroadcastGroup, error) {
	return c.listBroadcastGroup(ctx, site)
}

func (c *Client) GetBroadcastGroup(ctx context.Context, site, id string) (*BroadcastGroup, error) {
	return c.getBroadcastGroup(ctx, site, id)
}

func (c *Client) DeleteBroadcastGroup(ctx context.Context, site, id string) error {
	return c.deleteBroadcastGroup(ctx, site, id)
}

func (c *Client) CreateBroadcastGroup(
	ctx context.Context,
	site string,
	d *BroadcastGroup,
) (*BroadcastGroup, error) {
	return c.createBroadcastGroup(ctx, site, d)
}

func (c *Client) UpdateBroadcastGroup(
	ctx context.Context,
	site string,
	d *BroadcastGroup,
) (*BroadcastGroup, error) {
	return c.updateBroadcastGroup(ctx, site, d)
}
