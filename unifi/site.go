package unifi

import (
	"context"

	network "github.com/ubiquiti-community/go-unifi/client/network"
)

func (c *ApiClient) ListSites(ctx context.Context) ([]network.SiteOverview, error) {
	return FetchAll(ctx, func(offset int32) (*network.SiteOverviewPage, error) {
		resp, err := c.network.client.GetSiteOverviewPageWithResponse(ctx, &network.GetSiteOverviewPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}
