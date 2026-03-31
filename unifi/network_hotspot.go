package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/network"

	"github.com/google/uuid"
)

func (c *ApiClient) ListVouchers(ctx context.Context, site uuid.UUID) ([]network.HotspotVoucherDetails, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	return FetchAll(ctx, func(offset int32) (*network.HotspotVoucherDetailPage, error) {
		resp, err := c.network.client.GetVouchersWithResponse(ctx, site, &network.GetVouchersParams{
			Offset: Ptr(offset),
			Limit:  Ptr[int32](50),
		})
		if err != nil {
			return nil, err
		}
		return resp.JSON200, nil
	})
}

func (c *ApiClient) GetVoucher(ctx context.Context, site, id uuid.UUID) (*network.HotspotVoucherDetails, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.GetVoucherWithResponse(ctx, site, id)
	if err != nil {
		return nil, err
	}
	return resp.JSON200, nil
}

func (c *ApiClient) DeleteVoucher(ctx context.Context, site, id uuid.UUID) error {
	if c.network == nil {
		return fmt.Errorf("Network API is unavailable")
	}

	_, err := c.network.client.DeleteVoucher(ctx, site, id)
	return err
}

func (c *ApiClient) CreateVouchers(ctx context.Context, site uuid.UUID, data network.CreateVouchersJSONRequestBody) (*network.IntegrationVoucherCreationResultDto, error) {
	if c.network == nil {
		return nil, fmt.Errorf("Network API is unavailable")
	}

	resp, err := c.network.client.CreateVouchersWithResponse(ctx, site, data)
	if err != nil {
		return nil, err
	}
	return resp.JSON201, nil
}
