package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListACLRules(ctx context.Context, site uuid.UUID) ([]network.ACLRuleObject, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	return fetchAll(ctx, func(offset int32) (*network.IntegrationAclRulePageDto, error) {
		resp, err := c.network.client.GetAclRulePageWithResponse(ctx, site, &network.GetAclRulePageParams{
			Offset: ptr(offset),
			Limit:  ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetACLRule(ctx context.Context, site, id uuid.UUID) (*network.ACLRule, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.GetAclRuleWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteACLRule(ctx context.Context, site, id uuid.UUID) error {
	if c.network == nil {
		return fmt.Errorf("Network API is unavailable")
	}

	_, err := c.network.client.DeleteAclRule(ctx, site, id)
	return err
}

func (c *ApiClient) CreateACLRule(ctx context.Context, site uuid.UUID, data network.CreateAclRuleJSONRequestBody) (*network.ACLRule, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.CreateAclRuleWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}

func (c *ApiClient) UpdateACLRule(ctx context.Context, site, id uuid.UUID, data network.UpdateAclRuleJSONRequestBody) (*network.ACLRule, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.UpdateAclRuleWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}
