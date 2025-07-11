// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingRadioAi(ctx context.Context, site string) (*SettingRadioAi, error) {
	return c.getSettingRadioAi(ctx, site)
}

func (c *Client) UpdateSettingRadioAi(
	ctx context.Context,
	site string,
	d *SettingRadioAi,
) (*SettingRadioAi, error) {
	return c.updateSettingRadioAi(ctx, site, d)
}
