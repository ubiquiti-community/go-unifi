package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/protect"
)

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
