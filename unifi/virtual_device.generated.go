// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
	_ types.Number
	_ strconv.NumError
	_ strings.Builder
)

type VirtualDevice struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	HeightInMeters float64 `json:"heightInMeters,omitempty"`
	Locked         bool    `json:"locked"`
	MapID          string  `json:"map_id,omitempty"`
	Type           string  `json:"type,omitempty"` // uap|usg|usw
	X              string  `json:"x,omitempty"`
	Y              string  `json:"y,omitempty"`
}

func (dst *VirtualDevice) UnmarshalJSON(b []byte) error {
	type Alias VirtualDevice
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

func (c *ApiClient) listVirtualDevice(
	ctx context.Context,
	site string,
	params ...struct {
		key string
		val string
	},
) ([]VirtualDevice, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []VirtualDevice `json:"data"`
	}

	// Build URL with query parameters
	url := fmt.Sprintf("api/s/%s/rest/virtualdevice", site)
	if len(params) > 0 {
		// Build query string manually to avoid URL-encoding colons in MAC addresses
		var parts []string
		for _, p := range params {
			parts = append(parts, p.key+"="+p.val)
		}
		url = fmt.Sprintf("%s?%s", url, strings.Join(parts, "&"))
	}

	err := c.do(
		ctx,
		"GET",
		url,
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	return respBody.Data, nil
}

func (c *ApiClient) getVirtualDevice(
	ctx context.Context,
	site string,
	id string,
) (*VirtualDevice, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []VirtualDevice `json:"data"`
	}
	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("api/s/%s/rest/virtualdevice/%s", site, id),
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *ApiClient) deleteVirtualDevice(
	ctx context.Context,
	site string,
	id string,
) error {
	err := c.do(
		ctx,
		"DELETE",
		fmt.Sprintf("api/s/%s/rest/virtualdevice/%s", site, id),
		struct{}{},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) createVirtualDevice(
	ctx context.Context,
	site string,
	d *VirtualDevice,
) (*VirtualDevice, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []VirtualDevice `json:"data"`
	}

	err := c.do(
		ctx,
		"POST",
		fmt.Sprintf("api/s/%s/rest/virtualdevice", site),
		d,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}

func (c *ApiClient) updateVirtualDevice(
	ctx context.Context,
	site string,
	d *VirtualDevice,
) (*VirtualDevice, error) {
	var respBody struct {
		Meta meta            `json:"meta"`
		Data []VirtualDevice `json:"data"`
	}
	err := c.do(
		ctx,
		"PUT",
		fmt.Sprintf("api/s/%s/rest/virtualdevice/%s", site, d.ID),
		d,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	// UDM SE API returns empty data array on successful PUT.
	// In that case, fetch the updated resource via GET.
	if len(respBody.Data) == 0 {
		return c.getVirtualDevice(ctx, site, d.ID)
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}
