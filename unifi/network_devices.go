package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListDevices(ctx context.Context, site uuid.UUID) ([]network.AdoptedDeviceOverview, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	return FetchAll(ctx, func(offset int32) (*network.AdoptedDeviceOverviewPage, error) {
		resp, err := c.network.client.GetAdoptedDeviceOverviewPageWithResponse(ctx, site, &network.GetAdoptedDeviceOverviewPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetDevice(ctx context.Context, site, id uuid.UUID) (*network.AdoptedDeviceDetails, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.GetAdoptedDeviceDetailsWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) AdoptDevice(ctx context.Context, site uuid.UUID, data network.AdoptDeviceJSONRequestBody) (*network.AdoptedDeviceDetails, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.AdoptDeviceWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) RemoveDevice(ctx context.Context, site, id uuid.UUID) error {
	if c.network == nil {
		return fmt.Errorf("Network API is unavailable")
	}

	_, err := c.network.client.RemoveDevice(ctx, site, id)
	return err
}

func (c *ApiClient) ExecuteDeviceAction(ctx context.Context, site, id uuid.UUID, data network.ExecuteAdoptedDeviceActionJSONRequestBody) error {
	if c.network == nil {
		return fmt.Errorf("Network API is unavailable")
	}

	_, err := c.network.client.ExecuteAdoptedDeviceActionWithResponse(ctx, site, id, data)
	return err
}

func (c *ApiClient) GetDeviceStatistics(ctx context.Context, site, id uuid.UUID) (*network.LatestStatisticsForADevice, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.GetAdoptedDeviceLatestStatisticsWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) ListPendingDevices(ctx context.Context) ([]network.DevicePendingAdoption, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	return FetchAll(ctx, func(offset int32) (*network.DevicePendingAdoptionPage, error) {
		resp, err := c.network.client.GetPendingDevicePageWithResponse(ctx, &network.GetPendingDevicePageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) ListDeviceTags(ctx context.Context, site uuid.UUID) ([]network.DeviceTag, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	return FetchAll(ctx, func(offset int32) (*network.IntegrationDeviceTagPageDto, error) {
		resp, err := c.network.client.GetDeviceTagPageWithResponse(ctx, site, &network.GetDeviceTagPageParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}
