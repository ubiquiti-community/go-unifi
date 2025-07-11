// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingTrafficFlow(
	ctx context.Context,
	site string,
) (*SettingTrafficFlow, error) {
	return c.getSettingTrafficFlow(ctx, site)
}

func (c *Client) UpdateSettingTrafficFlow(
	ctx context.Context,
	site string,
	d *SettingTrafficFlow,
) (*SettingTrafficFlow, error) {
	return c.updateSettingTrafficFlow(ctx, site, d)
}
