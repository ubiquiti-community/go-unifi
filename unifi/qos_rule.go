// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *ApiClient) ListQOSRule(ctx context.Context, site string) ([]QOSRule, error) {
	return c.listQOSRule(ctx, site)
}

func (c *ApiClient) GetQOSRule(
	ctx context.Context,
	site,
	id string,
) (*QOSRule, error) {
	return c.getQOSRule(ctx, site, id)
}

func (c *ApiClient) DeleteQOSRule(ctx context.Context, site, id string) error {
	return c.deleteQOSRule(ctx, site, id)
}

func (c *ApiClient) CreateQOSRule(ctx context.Context, site string, d *QOSRule) (*QOSRule, error) {
	return c.createQOSRule(ctx, site, d)
}

func (c *ApiClient) UpdateQOSRule(ctx context.Context, site string, d *QOSRule) (*QOSRule, error) {
	return c.updateQOSRule(ctx, site, d)
}
