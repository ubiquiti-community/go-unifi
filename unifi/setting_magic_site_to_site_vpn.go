// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingMagicSiteToSiteVpn(ctx context.Context, site string) (*SettingMagicSiteToSiteVpn, error) {
	return c.getSettingMagicSiteToSiteVpn(ctx, site)
}

func (c *Client) UpdateSettingMagicSiteToSiteVpn(ctx context.Context, site string, d *SettingMagicSiteToSiteVpn) (*SettingMagicSiteToSiteVpn, error) {
	return c.updateSettingMagicSiteToSiteVpn(ctx, site, d)
}
