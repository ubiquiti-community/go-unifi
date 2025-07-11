// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingMdns(ctx context.Context, site string) (*SettingMdns, error) {
	return c.getSettingMdns(ctx, site)
}

func (c *Client) UpdateSettingMdns(
	ctx context.Context,
	site string,
	d *SettingMdns,
) (*SettingMdns, error) {
	return c.updateSettingMdns(ctx, site, d)
}
