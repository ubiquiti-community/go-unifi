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

type DpiGroup struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	DPIappIDs []string `json:"dpiapp_ids,omitempty"` // [\d\w]+
	Enabled   bool     `json:"enabled"`
	Name      string   `json:"name,omitempty"` // .{1,128}
}

func (dst *DpiGroup) UnmarshalJSON(b []byte) error {
	type Alias DpiGroup
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

func (c *Client) listDpiGroup(ctx context.Context, site string) ([]DpiGroup, error) {
	var respBody struct {
		Meta meta       `json:"meta"`
		Data []DpiGroup `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("api/s/%s/rest/dpigroup", site), nil, &respBody)
	if err != nil {
		return nil, err
	}
	return respBody.Data, nil
}

func (c *Client) getDpiGroup(ctx context.Context, site, id string) (*DpiGroup, error) {
	var respBody struct {
		Meta meta       `json:"meta"`
		Data []DpiGroup `json:"data"`
	}
	err := c.do(ctx, "GET", fmt.Sprintf("api/s/%s/rest/dpigroup/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteDpiGroup(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("api/s/%s/rest/dpigroup/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createDpiGroup(ctx context.Context, site string, d *DpiGroup) (*DpiGroup, error) {
	var respBody struct {
		Meta meta       `json:"meta"`
		Data []DpiGroup `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("api/s/%s/rest/dpigroup", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}

func (c *Client) updateDpiGroup(ctx context.Context, site string, d *DpiGroup) (*DpiGroup, error) {
	var respBody struct {
		Meta meta       `json:"meta"`
		Data []DpiGroup `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("api/s/%s/rest/dpigroup/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}
