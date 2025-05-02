// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingSnmp(ctx context.Context, site string) (*SettingSnmp, error) {
	return c.getSettingSnmp(ctx, site)
}

func (c *Client) UpdateSettingSnmp(ctx context.Context, site string, d *SettingSnmp) (*SettingSnmp, error) {
	return c.updateSettingSnmp(ctx, site, d)
}
