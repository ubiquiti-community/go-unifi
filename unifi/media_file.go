// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListMediaFile(ctx context.Context, site string) ([]MediaFile, error) {
	return c.listMediaFile(ctx, site)
}

func (c *Client) GetMediaFile(ctx context.Context, site, id string) (*MediaFile, error) {
	return c.getMediaFile(ctx, site, id)
}

func (c *Client) DeleteMediaFile(ctx context.Context, site, id string) error {
	return c.deleteMediaFile(ctx, site, id)
}

func (c *Client) CreateMediaFile(ctx context.Context, site string, d *MediaFile) (*MediaFile, error) {
	return c.createMediaFile(ctx, site, d)
}

func (c *Client) UpdateMediaFile(ctx context.Context, site string, d *MediaFile) (*MediaFile, error) {
	return c.updateMediaFile(ctx, site, d)
}
