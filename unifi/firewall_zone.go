// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *ApiClient) ListFirewallZone(ctx context.Context, site string) ([]FirewallZone, error) {
	return c.listFirewallZone(ctx, site)
}

func (c *ApiClient) GetFirewallZone(ctx context.Context, site, id string) (*FirewallZone, error) {
	return c.getFirewallZone(ctx, site, id)
}

func (c *ApiClient) DeleteFirewallZone(ctx context.Context, site, id string) error {
	return c.deleteFirewallZone(ctx, site, id)
}

func (c *ApiClient) CreateFirewallZone(ctx context.Context, site string, d *FirewallZone) (*FirewallZone, error) {
	return c.createFirewallZone(ctx, site, d)
}

func (c *ApiClient) UpdateFirewallZone(ctx context.Context, site string, d *FirewallZone) (*FirewallZone, error) {
	return c.updateFirewallZone(ctx, site, d)
}
