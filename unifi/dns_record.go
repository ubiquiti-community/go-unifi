// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
)

func (c *Client) ListDNSRecord(ctx context.Context, site string) ([]DNSRecord, error) {
	return c.listDNSRecord(ctx, site)
}

func (c *Client) GetDNSRecord(ctx context.Context, site, id string) (*DNSRecord, error) {
	respBody, err := c.listDNSRecord(ctx, site)
	if err != nil {
		return nil, err
	}

	if len(respBody) == 0 {
		return nil, &NotFoundError{}
	}

	for _, val := range respBody {
		if val.ID == id {
			return &val, nil
		}
	}

	return nil, &NotFoundError{}
}

func (c *Client) DeleteDNSRecord(ctx context.Context, site, id string) error {
	return c.deleteDNSRecord(ctx, site, id)
}

func (c *Client) CreateDNSRecord(ctx context.Context, site string, d *DNSRecord) (*DNSRecord, error) {
	return c.createDNSRecord(ctx, site, d)
}

func (c *Client) UpdateDNSRecord(ctx context.Context, site string, d *DNSRecord) (*DNSRecord, error) {
	return c.updateDNSRecord(ctx, site, d)
}
