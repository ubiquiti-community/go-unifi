package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListWifiBroadcasts(ctx context.Context, site uuid.UUID) ([]network.WifiBroadcastOverview, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	return fetchAll(ctx, func(offset int32) (*network.IntegrationWifiBroadcastPageDto, error) {
		resp, err := c.network.client.GetWifiBroadcastPageWithResponse(ctx, site, &network.GetWifiBroadcastPageParams{
			Offset: ptr(offset),
			Limit:  ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetWifiBroadcast(ctx context.Context, site, id uuid.UUID) (*network.WifiBroadcastDetails, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.GetWifiBroadcastDetailsWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteWifiBroadcast(ctx context.Context, site, id uuid.UUID) error {
	if c.network == nil {
		return fmt.Errorf("Network API is unavailable")
	}

	_, err := c.network.client.DeleteWifiBroadcast(ctx, site, id, &network.DeleteWifiBroadcastParams{})
	return err
}

func (c *ApiClient) CreateWifiBroadcast(ctx context.Context, site uuid.UUID, data network.CreateWifiBroadcastJSONRequestBody) (*network.WifiBroadcastDetails, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.CreateWifiBroadcastWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}

func (c *ApiClient) UpdateWifiBroadcast(ctx context.Context, site, id uuid.UUID, data network.UpdateWifiBroadcastJSONRequestBody) (*network.WifiBroadcastDetails, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.UpdateWifiBroadcastWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}
