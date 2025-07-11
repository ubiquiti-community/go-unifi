// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListHotspot2Conf(ctx context.Context, site string) ([]Hotspot2Conf, error) {
	return c.listHotspot2Conf(ctx, site)
}

func (c *Client) GetHotspot2Conf(ctx context.Context, site, id string) (*Hotspot2Conf, error) {
	return c.getHotspot2Conf(ctx, site, id)
}

func (c *Client) DeleteHotspot2Conf(ctx context.Context, site, id string) error {
	return c.deleteHotspot2Conf(ctx, site, id)
}

func (c *Client) CreateHotspot2Conf(
	ctx context.Context,
	site string,
	d *Hotspot2Conf,
) (*Hotspot2Conf, error) {
	return c.createHotspot2Conf(ctx, site, d)
}

func (c *Client) UpdateHotspot2Conf(
	ctx context.Context,
	site string,
	d *Hotspot2Conf,
) (*Hotspot2Conf, error) {
	return c.updateHotspot2Conf(ctx, site, d)
}
