package unifi

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strconv"
)

// GetClientByMAC returns slightly different information than GetClient, as they
// use separate endpoints for their lookups. Specifically IP is only returned
// by this method.
func (c *ApiClient) GetClientByMAC(ctx context.Context, site, mac string) (*Client, error) {
	resp, err := c.ListClient(ctx, site)
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, &NotFoundError{}
	}
	if i := slices.IndexFunc(resp, func(d Client) bool { return d.MAC == mac }); i >= 0 {
		d := resp[i]
		return &d, nil
	} else {
		return nil, &NotFoundError{}
	}
}

// stamgrMeta issues a station-manager command and returns the response rc
// alongside the client data, so callers can distinguish a genuine miss from a
// success that echoed no client object.
func (c *ApiClient) stamgrMeta(
	ctx context.Context,
	site, cmd string,
	data map[string]any,
) (string, []Client, error) {
	reqBody := map[string]any{}

	maps.Copy(reqBody, data)

	reqBody["cmd"] = cmd

	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Client `json:"data"`
	}

	err := c.do(ctx, http.MethodPost, fmt.Sprintf("api/s/%s/cmd/stamgr", site), reqBody, &respBody)
	if err != nil {
		return "", nil, err
	}

	return respBody.Meta.RC, respBody.Data, nil
}

func (c *ApiClient) stamgr(
	ctx context.Context,
	site, cmd string,
	data map[string]any,
) ([]Client, error) {
	_, clients, err := c.stamgrMeta(ctx, site, cmd, data)
	return clients, err
}

// AuthorizeClientByMAC authorizes a guest WiFi session by MAC via the
// authorize-guest station-manager command. minutes must be a decimal integer
// string (e.g. "480"); it is sent to the controller as a JSON number so the
// session expiry is honoured. If minutes is empty or unparseable the field is
// omitted and the controller applies its own default (typically unlimited).
// Some UniFi OS builds answer authorize-guest with rc:"ok" and an empty data
// array even though the client was authorized, so unlike the other station
// commands an empty reply is treated as success.
func (c *ApiClient) AuthorizeClientByMAC(ctx context.Context, site, mac, apMAC, minutes string) error {
	params := map[string]any{
		"ap_mac": apMAC,
		"mac":    mac,
	}
	if m, err := strconv.Atoi(minutes); err == nil && m > 0 {
		params["minutes"] = m
	}

	rc, clients, err := c.stamgrMeta(ctx, site, "authorize-guest", params)
	if err != nil {
		return err
	}
	if rc == "ok" && len(clients) <= 1 {
		return nil
	}

	return &NotFoundError{}
}

func (c *ApiClient) BlockClientByMAC(ctx context.Context, site, mac string) error {
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

func (c *ApiClient) UnblockClientByMAC(ctx context.Context, site, mac string) error {
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

func (c *ApiClient) DeleteClientByMAC(ctx context.Context, site, mac string) error {
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

func (c *ApiClient) KickClientByMAC(ctx context.Context, site, mac string) error {
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

func (c *ApiClient) OverrideClientFingerprint(
	ctx context.Context,
	site, mac string,
	devIdOveride int64,
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
		DevIdOverride int64  `json:"dev_id_override"`
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

// ListClient returns all clients, optionally filtered by query parameters.
// The query parameter can contain any field from the Client struct to filter results.
// For example: map[string]string{"network_id": "abc123", "blocked": "true"}.
func (c *ApiClient) ListClient(
	ctx context.Context,
	site string,
	params ...map[string]string,
) ([]Client, error) {
	return c.listClient(ctx, site, params...)
}

// ListClientFiltered returns clients filtered by the provided key-value parameters.
// This is the map-based variant of ListClient for use outside the unifi package,
// where the anonymous struct parameter type cannot be spread as variadic args.
func (c *ApiClient) ListClientFiltered(ctx context.Context, site string, filters map[string]string) ([]Client, error) {
	return c.listClient(ctx, site, filters)
}

func (c *ApiClient) CreateClient(ctx context.Context, site string, d *Client) (*Client, error) {
	return c.createClient(ctx, site, d)
}

// GetClient returns information about a user from the REST endpoint.
// The GetClientByMAC method returns slightly different information (for
// example the IP) as it uses a different endpoint.
func (c *ApiClient) GetClient(ctx context.Context, site, id string) (*Client, error) {
	return c.getClient(ctx, site, id)
}

func (c *ApiClient) UpdateClient(ctx context.Context, site string, d *Client) (*Client, error) {
	return c.updateClient(ctx, site, d)
}
