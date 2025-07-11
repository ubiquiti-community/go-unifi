// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingRoamingAssistant(
	ctx context.Context,
	site string,
) (*SettingRoamingAssistant, error) {
	return c.getSettingRoamingAssistant(ctx, site)
}

func (c *Client) UpdateSettingRoamingAssistant(
	ctx context.Context,
	site string,
	d *SettingRoamingAssistant,
) (*SettingRoamingAssistant, error) {
	return c.updateSettingRoamingAssistant(ctx, site, d)
}
