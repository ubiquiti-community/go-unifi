package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListDNSPolicies(ctx context.Context, site uuid.UUID) ([]network.DNSPolicy, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	return fetchAll(ctx, func(offset int32) (*network.IntegrationDnsPolicyPageDto, error) {
		resp, err := c.network.client.GetDnsPolicyPageWithResponse(ctx, site, &network.GetDnsPolicyPageParams{
			Offset: ptr(offset),
			Limit:  ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetDNSPolicy(ctx context.Context, site, id uuid.UUID) (*network.DNSPolicy, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.GetDnsPolicyWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteDNSPolicy(ctx context.Context, site, id uuid.UUID) error {
	if c.network == nil {
		return fmt.Errorf("Network API is unavailable")
	}

	_, err := c.network.client.DeleteDnsPolicy(ctx, site, id)
	return err
}

func (c *ApiClient) CreateDNSPolicy(ctx context.Context, site uuid.UUID, data network.CreateDnsPolicyJSONRequestBody) (*network.DNSPolicy, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.CreateDnsPolicyWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}

func (c *ApiClient) UpdateDNSPolicy(ctx context.Context, site, id uuid.UUID, data network.UpdateDnsPolicyJSONRequestBody) (*network.DNSPolicy, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.UpdateDnsPolicyWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}
