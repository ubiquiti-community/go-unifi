// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type DNSRecord struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Enabled    bool   `json:"enabled"`
	Key        string `json:"key,omitempty"` // .{1,128}
	Port       int    `json:"port,omitempty"`
	Priority   string `json:"priority,omitempty"`    // .{1,128}
	RecordType string `json:"record_type,omitempty"` // A|AAAA|CNAME|MX|NS|PTR|SOA|SRV|TXT
	Ttl        int    `json:"ttl,omitempty"`
	Value      string `json:"value,omitempty"` // .{1,256}
	Weight     int    `json:"weight,omitempty"`
}

func (dst *DNSRecord) UnmarshalJSON(b []byte) error {
	type Alias DNSRecord
	aux := &struct {
		Port   emptyStringInt `json:"port"`
		Ttl    emptyStringInt `json:"ttl"`
		Weight emptyStringInt `json:"weight"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Port = int(aux.Port)
	dst.Ttl = int(aux.Ttl)
	dst.Weight = int(aux.Weight)

	return nil
}

func (c *Client) listDNSRecord(ctx context.Context, site string) ([]DNSRecord, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []DNSRecord `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/static-dns", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getDNSRecord(ctx context.Context, site, id string) (*DNSRecord, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []DNSRecord `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("%s/site/%s/static-dns/%s", c.apiV2Path, site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteDNSRecord(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("%s/site/%s/static-dns/%s", c.apiV2Path, site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createDNSRecord(ctx context.Context, site string, d *DNSRecord) (*DNSRecord, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []DNSRecord `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("%s/site/%s/static-dns", c.apiV2Path, site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateDNSRecord(ctx context.Context, site string, d *DNSRecord) (*DNSRecord, error) {
	var respBody struct {
		Meta meta        `json:"meta"`
		Data []DNSRecord `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("%s/site/%s/static-dns/%s", c.apiV2Path, site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}