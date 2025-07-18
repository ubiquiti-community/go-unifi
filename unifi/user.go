package unifi

import (
	"context"
	"errors"
	"fmt"
	"maps"
)

// GetUserByMAC returns slightly different information than GetUser, as they
// use separate endpoints for their lookups. Specifically IP is only returned
// by this method.
func (c *Client) GetUserByMAC(ctx context.Context, site, mac string) (*User, error) {
	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("api/s/%s/stat/user/%s", site, mac), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) CreateUser(ctx context.Context, site string, d *User) (*User, error) {
	reqBody := struct {
		Objects []struct {
			Data *User `json:"data"`
		} `json:"objects"`
	}{
		Objects: []struct {
			Data *User `json:"data"`
		}{
			{Data: d},
		},
	}

	var respBody struct {
		Meta meta `json:"meta"`
		Data []struct {
			Meta meta   `json:"meta"`
			Data []User `json:"data"`
		} `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("api/s/%s/group/user", site), reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, errors.New("malformed group response")
	}

	if err := respBody.Data[0].Meta.error(); err != nil {
		return nil, err
	}

	if len(respBody.Data[0].Data) != 1 {
		return nil, &NotFoundError{}
	}

	user := respBody.Data[0].Data[0]

	return &user, nil
}

func (c *Client) stamgr(
	ctx context.Context,
	site, cmd string,
	data map[string]any,
) ([]User, error) {
	reqBody := map[string]any{}

	maps.Copy(reqBody, data)

	reqBody["cmd"] = cmd

	var respBody struct {
		Meta meta   `json:"meta"`
		Data []User `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("api/s/%s/cmd/stamgr", site), reqBody, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) BlockUserByMAC(ctx context.Context, site, mac string) error {
	users, err := c.stamgr(ctx, site, "block-sta", map[string]any{
		"mac": mac,
	})
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return &NotFoundError{}
	}
	return nil
}

func (c *Client) UnblockUserByMAC(ctx context.Context, site, mac string) error {
	users, err := c.stamgr(ctx, site, "unblock-sta", map[string]any{
		"mac": mac,
	})
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return &NotFoundError{}
	}
	return nil
}

func (c *Client) DeleteUserByMAC(ctx context.Context, site, mac string) error {
	users, err := c.stamgr(ctx, site, "forget-sta", map[string]any{
		"macs": []string{mac},
	})
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return &NotFoundError{}
	}
	return nil
}

func (c *Client) KickUserByMAC(ctx context.Context, site, mac string) error {
	users, err := c.stamgr(ctx, site, "kick-sta", map[string]any{
		"mac": mac,
	})
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return &NotFoundError{}
	}
	return nil
}

func (c *Client) OverrideUserFingerprint(
	ctx context.Context,
	site, mac string,
	devIdOveride int,
) error {
	reqBody := map[string]any{
		"mac":             mac,
		"dev_id_override": devIdOveride,
		"search_query":    "",
	}

	var reqMethod string
	if devIdOveride == 0 {
		reqMethod = "DELETE"
	} else {
		reqMethod = "PUT"
	}

	var respBody struct {
		Mac           string `json:"mac"`
		DevIdOverride int    `json:"dev_id_override"`
		SearchQuery   string `json:"search_query"`
	}

	err := c.do(
		ctx,
		reqMethod,
		fmt.Sprintf("v2/api/site/%s/station/%s/fingerprint_override", site, mac),
		reqBody,
		&respBody,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ListUser(ctx context.Context, site string) ([]User, error) {
	return c.listUser(ctx, site)
}

// GetUser returns information about a user from the REST endpoint.
// The GetUserByMAC method returns slightly different information (for
// example the IP) as it uses a different endpoint.
func (c *Client) GetUser(ctx context.Context, site, id string) (*User, error) {
	return c.getUser(ctx, site, id)
}

func (c *Client) UpdateUser(ctx context.Context, site string, d *User) (*User, error) {
	return c.updateUser(ctx, site, d)
}
