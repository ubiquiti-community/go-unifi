package unifi

import (
	"context"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/client/protect"
)

func (c *ApiClient) GetNvr(ctx context.Context) (*protect.Nvr, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1NvrWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}

func (c *ApiClient) GetUnvrStorage(ctx context.Context) (*protect.UnvrStorage, error) {
	if c.protect == nil {
		return nil, fmt.Errorf("Protect API is unavailable")
	}

	resp, err := c.protect.client.GetV1UnvrStorageWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}
