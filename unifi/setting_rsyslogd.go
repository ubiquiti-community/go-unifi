// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) GetSettingRsyslogd(ctx context.Context, site string) (*SettingRsyslogd, error) {
	return c.getSettingRsyslogd(ctx, site)
}

func (c *Client) UpdateSettingRsyslogd(ctx context.Context, site string, d *SettingRsyslogd) (*SettingRsyslogd, error) {
	return c.updateSettingRsyslogd(ctx, site, d)
}
