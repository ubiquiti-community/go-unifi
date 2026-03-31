package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/protect"
)

func (c *ApiClient) ListLights(ctx context.Context) ([]protect.Light, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1LightsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200, nil
}

func (c *ApiClient) GetLight(ctx context.Context, id protect.LightId) (*protect.Light, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1LightsIdWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (c *ApiClient) UpdateLight(ctx context.Context, id protect.LightId, params protect.PatchV1LightsIdJSONRequestBody) (*protect.Light, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.PatchV1LightsIdWithResponse(ctx, id, params)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (c *ApiClient) ListSensors(ctx context.Context) ([]protect.Sensor, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1SensorsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200, nil
}

func (c *ApiClient) GetSensor(ctx context.Context, id protect.SensorId) (*protect.Sensor, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1SensorsIdWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (c *ApiClient) UpdateSensor(ctx context.Context, id protect.SensorId, params protect.PatchV1SensorsIdJSONRequestBody) (*protect.Sensor, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.PatchV1SensorsIdWithResponse(ctx, id, params)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

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
