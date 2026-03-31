package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/protect"
)

func (c *ApiClient) ListChimes(ctx context.Context) ([]protect.Chime, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1ChimesWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200, nil
}

func (c *ApiClient) GetChime(ctx context.Context, id protect.ChimeId) (*protect.Chime, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1ChimesIdWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (c *ApiClient) UpdateChime(ctx context.Context, id protect.ChimeId, params protect.PatchV1ChimesIdJSONRequestBody) (*protect.Chime, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.PatchV1ChimesIdWithResponse(ctx, id, params)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}
