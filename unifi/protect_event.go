package unifi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ubiquiti-community/go-unifi/client/protect"
)

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
