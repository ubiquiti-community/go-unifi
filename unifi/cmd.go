package unifi

import (
	"context"
	"fmt"
)

// just to fix compile issues with the import.
var (
	_ fmt.Formatter
	_ context.Context
)

type Cmd struct {
	Command  string `json:"cmd"`
	MAC      string `json:"mac,omitempty"`
	PortIDX  *int   `json:"port_idx,omitempty"`
	FileName string `json:"filename,omitempty"`
	SiteID   string `json:"site_id,omitempty"`
}

func (c *Client) ExecuteCmd(ctx context.Context, site string, mgr string, cmd Cmd) (any, error) {
	var respBody struct{}

	err := c.do(ctx, "POST", fmt.Sprintf("%s/s/%s/cmd/%s", c.apiPath, site, mgr), &cmd, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
