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

	resp, err := c.protect.client.GetV1NvrsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}
