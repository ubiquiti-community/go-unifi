package unifi

import (
	"context"

	network "github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListNetworks(ctx context.Context, site uuid.UUID) ([]network.NetworkOverview, error) {
	return FetchAll(ctx, func(offset int32) (*network.NetworkOverviewPage, error) {
		resp, err := c.network.client.GetNetworksOverviewPageWithResponse(ctx, site, &network.GetNetworksOverviewPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetNetwork(ctx context.Context, site, id uuid.UUID) (*network.NetworkDetails, error) {
	resp, err := c.network.client.GetNetworkDetailsWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteNetwork(ctx context.Context, site, id uuid.UUID) error {
	_, err := c.network.client.DeleteNetwork(ctx, site, id, &network.DeleteNetworkParams{})
	return err
}

func (c *ApiClient) CreateNetwork(ctx context.Context, site uuid.UUID, data network.CreateNetworkJSONRequestBody) (*network.NetworkDetails, error) {
	resp, err := c.network.client.CreateNetworkWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}

func (c *ApiClient) UpdateNetwork(ctx context.Context, site, id uuid.UUID, data network.UpdateNetworkJSONRequestBody) (*network.NetworkDetails, error) {
	resp, err := c.network.client.UpdateNetworkWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) GetNetworkReferences(ctx context.Context, site, id uuid.UUID) (*network.NetworkReferences, error) {
	resp, err := c.network.client.GetNetworkReferencesWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}
