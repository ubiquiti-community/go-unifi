// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type Map struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Lat        string  `json:"lat,omitempty"`       // ^([-]?[\d]+[.]?[\d]*([eE][-+]?[\d]+)?)$
	Lng        string  `json:"lng,omitempty"`       // ^([-]?[\d]+[.]?[\d]*([eE][-+]?[\d]+)?)$
	MapTypeID  string  `json:"mapTypeId,omitempty"` // satellite|roadmap|hybrid|terrain
	Name       string  `json:"name,omitempty"`
	OffsetLeft float64 `json:"offset_left,omitempty"`
	OffsetTop  float64 `json:"offset_top,omitempty"`
	Opacity    float64 `json:"opacity,omitempty"` // ^(0(\.[\d]{1,2})?|1)$|^$
	Selected   bool    `json:"selected"`
	Tilt       int     `json:"tilt,omitempty"`
	Type       string  `json:"type,omitempty"` // designerMap|imageMap|googleMap
	Unit       string  `json:"unit,omitempty"` // m|f
	Upp        float64 `json:"upp,omitempty"`
	Zoom       int     `json:"zoom,omitempty"`
}

func (dst *Map) UnmarshalJSON(b []byte) error {
	type Alias Map
	aux := &struct {
		Tilt types.Number `json:"tilt"`
		Zoom types.Number `json:"zoom"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if val, err := aux.Tilt.Int64(); err == nil {
		dst.Tilt = int(val)
	}
	if val, err := aux.Zoom.Int64(); err == nil {
		dst.Zoom = int(val)
	}

	return nil
}

func (c *Client) listMap(ctx context.Context, site string) ([]Map, error) {
	var respBody struct {
		Meta meta  `json:"meta"`
		Data []Map `json:"data"`
	}

	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("api/s/%s/rest/map", site),
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	return respBody.Data, nil
}

func (c *Client) getMap(
	ctx context.Context,
	site string,
	id string,
) (*Map, error) {
	var respBody struct {
		Meta meta  `json:"meta"`
		Data []Map `json:"data"`
	}
	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("api/s/%s/rest/map/%s", site, id),
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

func (c *Client) deleteMap(
	ctx context.Context,
	site string,
	id string,
) error {
	err := c.do(
		ctx,
		"DELETE",
		fmt.Sprintf("api/s/%s/rest/map/%s", site, id),
		struct{}{},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createMap(
	ctx context.Context,
	site string,
	d *Map,
) (*Map, error) {
	var respBody struct {
		Meta meta  `json:"meta"`
		Data []Map `json:"data"`
	}

	err := c.do(
		ctx,
		"POST",
		fmt.Sprintf("api/s/%s/rest/map", site),
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

func (c *Client) updateMap(
	ctx context.Context,
	site string,
	d *Map,
) (*Map, error) {
	var respBody struct {
		Meta meta  `json:"meta"`
		Data []Map `json:"data"`
	}
	err := c.do(
		ctx,
		"PUT",
		fmt.Sprintf("api/s/%s/rest/map/%s", site, d.ID),
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
