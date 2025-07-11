// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListSpatialRecord(ctx context.Context, site string) ([]SpatialRecord, error) {
	return c.listSpatialRecord(ctx, site)
}

func (c *Client) GetSpatialRecord(ctx context.Context, site, id string) (*SpatialRecord, error) {
	return c.getSpatialRecord(ctx, site, id)
}

func (c *Client) DeleteSpatialRecord(ctx context.Context, site, id string) error {
	return c.deleteSpatialRecord(ctx, site, id)
}

func (c *Client) CreateSpatialRecord(
	ctx context.Context,
	site string,
	d *SpatialRecord,
) (*SpatialRecord, error) {
	return c.createSpatialRecord(ctx, site, d)
}

func (c *Client) UpdateSpatialRecord(
	ctx context.Context,
	site string,
	d *SpatialRecord,
) (*SpatialRecord, error) {
	return c.updateSpatialRecord(ctx, site, d)
}
