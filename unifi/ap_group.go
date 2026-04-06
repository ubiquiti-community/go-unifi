package unifi

import (
	"context"
	"fmt"
	"net/http"
)

// just to fix compile issues with the import.
var (
	_ fmt.Formatter
	_ context.Context
)

// This is a v2 API object, so manually coded for now, need to figure out generation...

type APGroup struct {
	ID string `json:"_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenId string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Name       string   `json:"name"`
	DeviceMacs []string `json:"device_macs"`
}

func (c *ApiClient) ListAPGroup(ctx context.Context, site string) ([]APGroup, error) {
	var respBody []APGroup

	err := c.do(ctx, http.MethodGet, fmt.Sprintf("v2/api/site/%s/apgroups", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *ApiClient) CreateAPGroup(ctx context.Context, site string, d *APGroup) (*APGroup, error) {
	var respBody APGroup

	err := c.do(ctx, http.MethodPost, fmt.Sprintf("v2/api/site/%s/apgroups", site), d, &respBody)
	if err != nil {
		return nil, err
	}
	return &respBody, nil
}

func (c *ApiClient) GetAPGroup(ctx context.Context, site, id string) (*APGroup, error) {
	groups, err := c.ListAPGroup(ctx, site)
	if err != nil {
		return nil, err
	}
	for _, g := range groups {
		if g.ID == id {
			return &g, nil
		}
	}
	return nil, &NotFoundError{}
}

func (c *ApiClient) UpdateAPGroup(ctx context.Context, site string, d *APGroup) (*APGroup, error) {
	var respBody APGroup
	err := c.do(ctx, http.MethodPut, fmt.Sprintf("v2/api/site/%s/apgroups/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}
	return &respBody, nil
}

func (c *ApiClient) DeleteAPGroup(ctx context.Context, site, id string) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("v2/api/site/%s/apgroups/%s", site, id), struct{}{}, nil)
}
