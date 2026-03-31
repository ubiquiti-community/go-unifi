package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/protect"
)

func (c *ApiClient) ListCameras(ctx context.Context) ([]protect.Camera, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1CamerasWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200, nil
}

func (c *ApiClient) GetCamera(ctx context.Context, id protect.CameraId) (*protect.Camera, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1CamerasIdWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (c *ApiClient) UpdateCamera(ctx context.Context, id protect.CameraId, params protect.PatchV1CamerasIdJSONRequestBody) (*protect.Camera, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.PatchV1CamerasIdWithResponse(ctx, id, params)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (c *ApiClient) StartCameraPatrol(ctx context.Context, id protect.CameraId, slot protect.ActivePatrolSlotString) error {
	if c.protect == nil {
		return fmt.Errorf("Protect API is unavailable")
	}

	_, err := c.protect.client.PostV1CamerasIdPtzPatrolStartSlotWithResponse(ctx, id, slot)
	if err != nil {
		return err
	}

	return nil
}

func (c *ApiClient) StopCameraPatrol(ctx context.Context, id protect.CameraId) error {
	if c.protect == nil {
		return fmt.Errorf("Protect API is unavailable")
	}

	_, err := c.protect.client.PostV1CamerasIdPtzPatrolStopWithResponse(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *ApiClient) GoToCameraPreset(ctx context.Context, id protect.CameraId, slot string) error {
	if c.protect == nil {
		return fmt.Errorf("Protect API is unavailable")
	}

	_, err := c.protect.client.PostV1CamerasIdPtzGotoSlotWithResponse(ctx, id, slot)
	if err != nil {
		return err
	}

	return nil
}
