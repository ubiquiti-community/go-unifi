package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/protect"
)

// --- Viewers ---

func (c *ApiClient) ListViewers(ctx context.Context) ([]protect.Viewer, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1ViewersWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return *resp.JSON200, nil
}

func (c *ApiClient) GetViewer(ctx context.Context, id protect.ViewerId) (*protect.Viewer, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1ViewersIdWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}

func (c *ApiClient) UpdateViewer(ctx context.Context, id protect.ViewerId, params protect.PatchV1ViewersIdJSONRequestBody) (*protect.Viewer, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.PatchV1ViewersIdWithResponse(ctx, id, params)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}

// --- LiveViews ---

func (c *ApiClient) ListLiveviews(ctx context.Context) ([]protect.Liveview, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1LiveviewsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return *resp.JSON200, nil
}

func (c *ApiClient) GetLiveview(ctx context.Context, id protect.LiveviewId) (*protect.Liveview, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1LiveviewsIdWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}

func (c *ApiClient) CreateLiveview(ctx context.Context, params protect.PostV1LiveviewsJSONRequestBody) (*protect.Liveview, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.PostV1LiveviewsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}

func (c *ApiClient) UpdateLiveview(ctx context.Context, id protect.LiveviewId, params protect.PatchV1LiveviewsIdJSONRequestBody) (*protect.Liveview, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.PatchV1LiveviewsIdWithResponse(ctx, id, params)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}

// --- Users ---

func (c *ApiClient) ListUsers(ctx context.Context) ([]protect.User, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1UsersWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return *resp.JSON200, nil
}

func (c *ApiClient) GetUser(ctx context.Context, id protect.UserId) (*protect.User, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1UsersIdWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}
