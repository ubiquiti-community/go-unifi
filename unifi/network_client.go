package unifi

import (
	"context"
	"fmt"

	network "github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListClients(ctx context.Context, site uuid.UUID) ([]network.ClientOverview, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	return FetchAll(ctx, func(offset int32) (*network.ClientOverviewPage, error) {
		resp, err := c.network.client.GetConnectedClientOverviewPageWithResponse(ctx, site, &network.GetConnectedClientOverviewPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetClient(ctx context.Context, site, id uuid.UUID) (*network.ClientDetails, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.GetConnectedClientDetailsWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) ExecuteClientAction(ctx context.Context, site, id uuid.UUID, data network.ExecuteConnectedClientActionJSONRequestBody) (*network.ClientActionResponse, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.ExecuteConnectedClientActionWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}
