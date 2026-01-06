package unifi

import (
	"context"
	"fmt"
	"slices"
)

func (c *Client) DeleteNetwork(ctx context.Context, site, id, name string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("api/s/%s/rest/networkconf/%s", site, id), struct {
		Name string `json:"name"`
	}{
		Name: name,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ListNetwork(ctx context.Context, site string) ([]Network, error) {
	return c.listNetwork(ctx, site)
}

func (c *Client) GetNetwork(ctx context.Context, site, id string) (*Network, error) {
	return c.getNetwork(ctx, site, id)
}

func (c *Client) GetNetworkByName(ctx context.Context, site, name string) (*Network, error) {
	networks, err := c.listNetwork(ctx, site)
	if err != nil {
		return nil, err
	}
	i := slices.IndexFunc(networks, func(n Network) bool {
		return n.Name == name
	})
	if i < 0 {
		return nil, fmt.Errorf("network with name %s not found", name)
	}
	network := networks[i]
	return &network, nil
}

func (c *Client) CreateNetwork(ctx context.Context, site string, d *Network) (*Network, error) {
	return c.createNetwork(ctx, site, d)
}

func (c *Client) UpdateNetwork(ctx context.Context, site string, d *Network) (*Network, error) {
	return c.updateNetwork(ctx, site, d)
}
