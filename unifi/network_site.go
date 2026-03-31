package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/network"
)

func (c *ApiClient) ListSites(ctx context.Context) ([]network.SiteOverview, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	return fetchAll(ctx, func(offset int32) (*network.SiteOverviewPage, error) {
		resp, err := c.network.client.GetSiteOverviewPageWithResponse(ctx, &network.GetSiteOverviewPageParams{
			Offset: ptr(offset),
			Limit:  ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}
