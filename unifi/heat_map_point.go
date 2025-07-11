// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListHeatMapPoint(ctx context.Context, site string) ([]HeatMapPoint, error) {
	return c.listHeatMapPoint(ctx, site)
}

func (c *Client) GetHeatMapPoint(ctx context.Context, site, id string) (*HeatMapPoint, error) {
	return c.getHeatMapPoint(ctx, site, id)
}

func (c *Client) DeleteHeatMapPoint(ctx context.Context, site, id string) error {
	return c.deleteHeatMapPoint(ctx, site, id)
}

func (c *Client) CreateHeatMapPoint(
	ctx context.Context,
	site string,
	d *HeatMapPoint,
) (*HeatMapPoint, error) {
	return c.createHeatMapPoint(ctx, site, d)
}

func (c *Client) UpdateHeatMapPoint(
	ctx context.Context,
	site string,
	d *HeatMapPoint,
) (*HeatMapPoint, error) {
	return c.updateHeatMapPoint(ctx, site, d)
}
