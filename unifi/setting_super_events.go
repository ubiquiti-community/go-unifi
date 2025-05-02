// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSuperEvents(ctx context.Context, site string) (*SettingSuperEvents, error) {
	return c.getSettingSuperEvents(ctx, site)
}

func (c *Client) UpdateSettingSuperEvents(ctx context.Context, site string, d *SettingSuperEvents) (*SettingSuperEvents, error) {
	return c.updateSettingSuperEvents(ctx, site, d)
}
