// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListChannelPlan(ctx context.Context, site string) ([]ChannelPlan, error) {
	return c.listChannelPlan(ctx, site)
}

func (c *Client) GetChannelPlan(ctx context.Context, site, id string) (*ChannelPlan, error) {
	return c.getChannelPlan(ctx, site, id)
}

func (c *Client) DeleteChannelPlan(ctx context.Context, site, id string) error {
	return c.deleteChannelPlan(ctx, site, id)
}

func (c *Client) CreateChannelPlan(
	ctx context.Context,
	site string,
	d *ChannelPlan,
) (*ChannelPlan, error) {
	return c.createChannelPlan(ctx, site, d)
}

func (c *Client) UpdateChannelPlan(
	ctx context.Context,
	site string,
	d *ChannelPlan,
) (*ChannelPlan, error) {
	return c.updateChannelPlan(ctx, site, d)
}
