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

type BroadcastGroup struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	MemberTable []string `json:"member_table,omitempty"`
	Name        string   `json:"name,omitempty"`
}

func (dst *BroadcastGroup) UnmarshalJSON(b []byte) error {
	type Alias BroadcastGroup
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

func (c *ApiClient) listBroadcastGroup(
	ctx context.Context,
	site string,
	params ...struct {
		key string
		val string
	},
) ([]BroadcastGroup, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []BroadcastGroup `json:"data"`
	}

	// Build URL with query parameters
	url := fmt.Sprintf("api/s/%s/rest/broadcastgroup", site)
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

func (c *ApiClient) getBroadcastGroup(
	ctx context.Context,
	site string,
	id string,
) (*BroadcastGroup, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []BroadcastGroup `json:"data"`
	}
	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("api/s/%s/rest/broadcastgroup/%s", site, id),
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

func (c *ApiClient) deleteBroadcastGroup(
	ctx context.Context,
	site string,
	id string,
) error {
	err := c.do(
		ctx,
		"DELETE",
		fmt.Sprintf("api/s/%s/rest/broadcastgroup/%s", site, id),
		struct{}{},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) createBroadcastGroup(
	ctx context.Context,
	site string,
	d *BroadcastGroup,
) (*BroadcastGroup, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []BroadcastGroup `json:"data"`
	}

	err := c.do(
		ctx,
		"POST",
		fmt.Sprintf("api/s/%s/rest/broadcastgroup", site),
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

func (c *ApiClient) updateBroadcastGroup(
	ctx context.Context,
	site string,
	d *BroadcastGroup,
) (*BroadcastGroup, error) {
	var respBody struct {
		Meta meta             `json:"meta"`
		Data []BroadcastGroup `json:"data"`
	}
	err := c.do(
		ctx,
		"PUT",
		fmt.Sprintf("api/s/%s/rest/broadcastgroup/%s", site, d.ID),
		d,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	// UDM SE API returns empty data array on successful PUT.
	// In that case, fetch the updated resource via GET.
	if len(respBody.Data) == 0 {
		return c.getBroadcastGroup(ctx, site, d.ID)
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}
