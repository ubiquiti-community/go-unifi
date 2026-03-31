package unifi

import (
	"context"

	network "github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListACLRules(ctx context.Context, site uuid.UUID) ([]network.ACLRuleObject, error) {
	return FetchAll(ctx, func(offset int32) (*network.IntegrationAclRulePageDto, error) {
		resp, err := c.network.client.GetAclRulePageWithResponse(ctx, site, &network.GetAclRulePageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetACLRule(ctx context.Context, site, id uuid.UUID) (*network.ACLRule, error) {
	resp, err := c.network.client.GetAclRuleWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteACLRule(ctx context.Context, site, id uuid.UUID) error {
	_, err := c.network.client.DeleteAclRule(ctx, site, id)
	return err
}

func (c *ApiClient) CreateACLRule(ctx context.Context, site uuid.UUID, data network.CreateAclRuleJSONRequestBody) (*network.ACLRule, error) {
	resp, err := c.network.client.CreateAclRuleWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}

func (c *ApiClient) UpdateACLRule(ctx context.Context, site, id uuid.UUID, data network.UpdateAclRuleJSONRequestBody) (*network.ACLRule, error) {
	resp, err := c.network.client.UpdateAclRuleWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}
