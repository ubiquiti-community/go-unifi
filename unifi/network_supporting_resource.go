package unifi

import (
	"context"

	network "github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListWANs(ctx context.Context, site uuid.UUID) ([]network.WANOverview, error) {
	return FetchAll(ctx, func(offset int32) (*network.WANOverviewPage, error) {
		resp, err := c.network.client.GetWansOverviewPageWithResponse(ctx, site, &network.GetWansOverviewPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) ListSiteToSiteVPNTunnels(ctx context.Context, site uuid.UUID) ([]network.SiteToSiteVPNTunnelOverview, error) {
	return FetchAll(ctx, func(offset int32) (*network.IntegrationSiteToSiteVpnTunnelOverviewPageDto, error) {
		resp, err := c.network.client.GetSiteToSiteVpnTunnelPageWithResponse(ctx, site, &network.GetSiteToSiteVpnTunnelPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) ListVPNServers(ctx context.Context, site uuid.UUID) ([]network.VPNServerOverview, error) {
	return FetchAll(ctx, func(offset int32) (*network.IntegrationVpnServerOverviewPageDto, error) {
		resp, err := c.network.client.GetVpnServerPageWithResponse(ctx, site, &network.GetVpnServerPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) ListRADIUSProfiles(ctx context.Context, site uuid.UUID) ([]network.RadiusProfileOverview, error) {
	return FetchAll(ctx, func(offset int32) (*network.RadiusProfileOverviewPage, error) {
		resp, err := c.network.client.GetRadiusProfileOverviewPageWithResponse(ctx, site, &network.GetRadiusProfileOverviewPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) ListDPIApplicationCategories(ctx context.Context) ([]network.DPICategory, error) {
	return FetchAll(ctx, func(offset int32) (*network.DPICategoryPage, error) {
		resp, err := c.network.client.GetDpiApplicationCategoriesWithResponse(ctx, &network.GetDpiApplicationCategoriesParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) ListDPIApplications(ctx context.Context) ([]network.DPIApplication, error) {
	return FetchAll(ctx, func(offset int32) (*network.DPIApplicationPage, error) {
		resp, err := c.network.client.GetDpiApplicationsWithResponse(ctx, &network.GetDpiApplicationsParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) ListCountries(ctx context.Context) ([]network.CountryDefinition, error) {
	return FetchAll(ctx, func(offset int32) (*network.CountryDefinitionPage, error) {
		resp, err := c.network.client.GetCountriesWithResponse(ctx, &network.GetCountriesParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}
