// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListScheduleTask(ctx context.Context, site string) ([]ScheduleTask, error) {
	return c.listScheduleTask(ctx, site)
}

func (c *Client) GetScheduleTask(ctx context.Context, site, id string) (*ScheduleTask, error) {
	return c.getScheduleTask(ctx, site, id)
}

func (c *Client) DeleteScheduleTask(ctx context.Context, site, id string) error {
	return c.deleteScheduleTask(ctx, site, id)
}

func (c *Client) CreateScheduleTask(ctx context.Context, site string, d *ScheduleTask) (*ScheduleTask, error) {
	return c.createScheduleTask(ctx, site, d)
}

func (c *Client) UpdateScheduleTask(ctx context.Context, site string, d *ScheduleTask) (*ScheduleTask, error) {
	return c.updateScheduleTask(ctx, site, d)
}
