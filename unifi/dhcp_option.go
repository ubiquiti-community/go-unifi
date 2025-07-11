// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListDHCPOption(ctx context.Context, site string) ([]DHCPOption, error) {
	return c.listDHCPOption(ctx, site)
}

func (c *Client) GetDHCPOption(ctx context.Context, site, id string) (*DHCPOption, error) {
	return c.getDHCPOption(ctx, site, id)
}

func (c *Client) DeleteDHCPOption(ctx context.Context, site, id string) error {
	return c.deleteDHCPOption(ctx, site, id)
}

func (c *Client) CreateDHCPOption(
	ctx context.Context,
	site string,
	d *DHCPOption,
) (*DHCPOption, error) {
	return c.createDHCPOption(ctx, site, d)
}

func (c *Client) UpdateDHCPOption(
	ctx context.Context,
	site string,
	d *DHCPOption,
) (*DHCPOption, error) {
	return c.updateDHCPOption(ctx, site, d)
}
