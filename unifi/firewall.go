package unifi

import (
	"context"

	network "github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

// Firewall Policies

func (c *ApiClient) ListFirewallPolicies(ctx context.Context, site uuid.UUID) ([]network.FirewallPolicy, error) {
	return FetchAll(ctx, func(offset int32) (*network.FirewallPolicyPage, error) {
		resp, err := c.network.client.GetFirewallPoliciesWithResponse(ctx, site, &network.GetFirewallPoliciesParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetFirewallPolicy(ctx context.Context, site, id uuid.UUID) (*network.FirewallPolicy, error) {
	resp, err := c.network.client.GetFirewallPolicyWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteFirewallPolicy(ctx context.Context, site, id uuid.UUID) error {
	_, err := c.network.client.DeleteFirewallPolicy(ctx, site, id)
	return err
}

func (c *ApiClient) CreateFirewallPolicy(ctx context.Context, site uuid.UUID, data network.CreateFirewallPolicyJSONRequestBody) (*network.FirewallPolicy, error) {
	resp, err := c.network.client.CreateFirewallPolicyWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}

func (c *ApiClient) UpdateFirewallPolicy(ctx context.Context, site, id uuid.UUID, data network.UpdateFirewallPolicyJSONRequestBody) (*network.FirewallPolicy, error) {
	resp, err := c.network.client.UpdateFirewallPolicyWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) ListFirewallZones(ctx context.Context, site uuid.UUID) ([]network.FirewallZone, error) {
	return FetchAll(ctx, func(offset int32) (*network.FirewallZonesPage, error) {
		resp, err := c.network.client.GetFirewallZonesWithResponse(ctx, site, &network.GetFirewallZonesParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetFirewallZone(ctx context.Context, site, id uuid.UUID) (*network.FirewallZone, error) {
	resp, err := c.network.client.GetFirewallZoneWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteFirewallZone(ctx context.Context, site, id uuid.UUID) error {
	_, err := c.network.client.DeleteFirewallZone(ctx, site, id)
	return err
}

func (c *ApiClient) CreateFirewallZone(ctx context.Context, site uuid.UUID, data network.CreateFirewallZoneJSONRequestBody) (*network.FirewallZone, error) {
	resp, err := c.network.client.CreateFirewallZoneWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}

func (c *ApiClient) UpdateFirewallZone(ctx context.Context, site, id uuid.UUID, data network.UpdateFirewallZoneJSONRequestBody) (*network.FirewallZone, error) {
	resp, err := c.network.client.UpdateFirewallZoneWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}
