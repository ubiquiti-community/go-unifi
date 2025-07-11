// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingNetflow(ctx context.Context, site string) (*SettingNetflow, error) {
	return c.getSettingNetflow(ctx, site)
}

func (c *Client) UpdateSettingNetflow(
	ctx context.Context,
	site string,
	d *SettingNetflow,
) (*SettingNetflow, error) {
	return c.updateSettingNetflow(ctx, site, d)
}
