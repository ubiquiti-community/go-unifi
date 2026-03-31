package unifi

import (
	"context"

	network "github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListDNSPolicies(ctx context.Context, site uuid.UUID) ([]network.DNSPolicy, error) {
	return FetchAll(ctx, func(offset int32) (*network.IntegrationDnsPolicyPageDto, error) {
		resp, err := c.network.client.GetDnsPolicyPageWithResponse(ctx, site, &network.GetDnsPolicyPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetDNSPolicy(ctx context.Context, site, id uuid.UUID) (*network.DNSPolicy, error) {
	resp, err := c.network.client.GetDnsPolicyWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteDNSPolicy(ctx context.Context, site, id uuid.UUID) error {
	_, err := c.network.client.DeleteDnsPolicy(ctx, site, id)
	return err
}

func (c *ApiClient) CreateDNSPolicy(ctx context.Context, site uuid.UUID, data network.CreateDnsPolicyJSONRequestBody) (*network.DNSPolicy, error) {
	resp, err := c.network.client.CreateDnsPolicyWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}

func (c *ApiClient) UpdateDNSPolicy(ctx context.Context, site, id uuid.UUID, data network.UpdateDnsPolicyJSONRequestBody) (*network.DNSPolicy, error) {
	resp, err := c.network.client.UpdateDnsPolicyWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}
