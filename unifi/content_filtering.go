// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *ApiClient) ListContentFiltering(ctx context.Context, site string) ([]ContentFiltering, error) {
	return c.listContentFiltering(ctx, site)
}

func (c *ApiClient) GetContentFiltering(
	ctx context.Context,
	site,
	id string,
) (*ContentFiltering, error) {
	return c.getContentFiltering(ctx, site, id)
}

func (c *ApiClient) DeleteContentFiltering(ctx context.Context, site, id string) error {
	return c.deleteContentFiltering(ctx, site, id)
}

func (c *ApiClient) CreateContentFiltering(ctx context.Context, site string, d *ContentFiltering) (*ContentFiltering, error) {
	return c.createContentFiltering(ctx, site, d)
}

func (c *ApiClient) UpdateContentFiltering(ctx context.Context, site string, d *ContentFiltering) (*ContentFiltering, error) {
	return c.updateContentFiltering(ctx, site, d)
}
