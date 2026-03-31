package unifi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ubiquiti-community/go-unifi/client/protect"
)

// --- Events ---

func (c *ApiClient) ListEvents(ctx context.Context, params *protect.GetV1EventsParams) ([]protect.Event, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1EventsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return *resp.JSON200, nil
}

// --- Recordings & Media ---

func (c *ApiClient) GetCameraRecording(ctx context.Context, id protect.CameraId, params *protect.GetV1CamerasIdRecordingParams) ([]byte, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1CamerasIdRecordingWithResponse(ctx, id, params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.Body, nil
}

func (c *ApiClient) GetCameraSnapshot(ctx context.Context, id protect.CameraId, params *protect.GetV1CamerasIdSnapshotParams) ([]byte, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1CamerasIdSnapshotWithResponse(ctx, id, params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.Body, nil
}

func (c *ApiClient) GetCameraThumbnail(ctx context.Context, id protect.CameraId, params *protect.GetV1CamerasIdThumbnailParams) ([]byte, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1CamerasIdThumbnailWithResponse(ctx, id, params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.Body, nil
}
