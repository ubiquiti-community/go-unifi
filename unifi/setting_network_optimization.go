// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingNetworkOptimization(ctx context.Context, site string) (*SettingNetworkOptimization, error) {
	return c.getSettingNetworkOptimization(ctx, site)
}

func (c *Client) UpdateSettingNetworkOptimization(ctx context.Context, site string, d *SettingNetworkOptimization) (*SettingNetworkOptimization, error) {
	return c.updateSettingNetworkOptimization(ctx, site, d)
}
