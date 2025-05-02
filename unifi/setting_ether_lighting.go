// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingEtherLighting(ctx context.Context, site string) (*SettingEtherLighting, error) {
	return c.getSettingEtherLighting(ctx, site)
}

func (c *Client) UpdateSettingEtherLighting(ctx context.Context, site string, d *SettingEtherLighting) (*SettingEtherLighting, error) {
	return c.updateSettingEtherLighting(ctx, site, d)
}
