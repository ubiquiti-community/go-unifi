// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingEvaluationScore(ctx context.Context, site string) (*SettingEvaluationScore, error) {
	return c.getSettingEvaluationScore(ctx, site)
}

func (c *Client) UpdateSettingEvaluationScore(ctx context.Context, site string, d *SettingEvaluationScore) (*SettingEvaluationScore, error) {
	return c.updateSettingEvaluationScore(ctx, site, d)
}
