package unifi

import (
	"context"

	network "github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListTrafficMatchingLists(ctx context.Context, site uuid.UUID) ([]network.TrafficMatchingList, error) {
	return FetchAll(ctx, func(offset int32) (*network.TrafficMatchingListsPage, error) {
		resp, err := c.network.client.GetTrafficMatchingListsWithResponse(ctx, site, &network.GetTrafficMatchingListsParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetTrafficMatchingList(ctx context.Context, site, id uuid.UUID) (*network.TrafficMatchingList, error) {
	resp, err := c.network.client.GetTrafficMatchingListWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteTrafficMatchingList(ctx context.Context, site, id uuid.UUID) error {
	_, err := c.network.client.DeleteTrafficMatchingList(ctx, site, id)
	return err
}

func (c *ApiClient) CreateTrafficMatchingList(ctx context.Context, site uuid.UUID, data network.CreateTrafficMatchingListJSONRequestBody) (*network.TrafficMatchingList, error) {
	resp, err := c.network.client.CreateTrafficMatchingListWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}

func (c *ApiClient) UpdateTrafficMatchingList(ctx context.Context, site, id uuid.UUID, data network.UpdateTrafficMatchingListJSONRequestBody) (*network.TrafficMatchingList, error) {
	resp, err := c.network.client.UpdateTrafficMatchingListWithResponse(ctx, site, id, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}
