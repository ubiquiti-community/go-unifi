// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type BGPConfig struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Enabled        bool   `json:"enabled"`
	Config         string `json:"frr_bgpd_config,omitempty"`
	UploadFileName string `json:"uploaded_file_name,omitempty"`
	Description    string `json:"description,omitempty"`
}

func (c *Client) GetBGPConfig(ctx context.Context, site string) (*BGPConfig, error) {
	var respBody []BGPConfig

	err := c.do(ctx, "GET", fmt.Sprintf("v2/api/site/%s/bgp/config", site), nil, &respBody)
	if err != nil {
		return nil, err
	}
	if len(respBody) == 0 {
		return nil, fmt.Errorf("no BGP config found for site %s", site)
	}

	return &respBody[0], nil
}

func (c *Client) deleteBGPConfig(ctx context.Context, site string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("v2/api/site/%s/bgp/config", site), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateBGPConfig(
	ctx context.Context,
	site string,
	d *BGPConfig,
) (*BGPConfig, error) {
	var respBody BGPConfig

	err := c.do(ctx, "POST", fmt.Sprintf("v2/api/site/%s/bgp/config", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}

func (c *Client) UpdateBGPConfig(
	ctx context.Context,
	site string,
	d *BGPConfig,
) (*BGPConfig, error) {
	var respBody BGPConfig

	err := c.do(ctx, "POST", fmt.Sprintf("v2/api/site/%s/bgp/config", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}

func (c *Client) DeleteBGPConfig(ctx context.Context, site string) error {
	return c.deleteBGPConfig(ctx, site)
}
