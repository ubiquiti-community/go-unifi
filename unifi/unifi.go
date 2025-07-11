package unifi //

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"
	"sync"
)

//go:generate go run ../cmd/fields/ -version-base-dir=../cmd/fields/ -output-dir=../unifi/ -latest

const (
	loginPath    = "/api/login"
	loginPathNew = "/api/auth/login"
)

type NotFoundError struct{}

func (err *NotFoundError) Error() string {
	return "not found"
}

type APIError struct {
	RC      string
	Message string
}

func (err *APIError) Error() string {
	return err.Message
}

type Client struct {
	// single thread client calls for CSRF, etc.
	sync.Mutex

	c       *http.Client
	baseURL *url.URL

	apiKey    string
	loginPath string
	apiPath   string // path to API, e.g. "proxy/network" for new style API, "/api" for old style API

	csrf string

	version string
}

func (c *Client) CSRFToken() string {
	return c.csrf
}

func (c *Client) Version() string {
	return c.version
}

func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *Client) SetBaseURL(base string) error {
	var err error
	c.baseURL, err = url.Parse(base)
	if err != nil {
		return err
	}

	// error for people who are still passing hard coded old paths
	if path := strings.TrimSuffix(c.baseURL.Path, "/"); path == "/api" {
		return fmt.Errorf("expected a base URL without the `/api`, got: %q", c.baseURL)
	}

	return nil
}

func (c *Client) SetHTTPClient(hc *http.Client) error {
	c.c = hc
	return nil
}

func (c *Client) setAPIUrlStyle(ctx context.Context) error {
	// check if new style API
	// this is modified from the unifi-poller (https://github.com/unifi-poller/unifi) implementation.
	// see https://github.com/unifi-poller/unifi/blob/4dc44f11f61a2e08bf7ec5b20c71d5bced837b5d/unifi.go#L101-L104
	// and https://github.com/unifi-poller/unifi/commit/43a6b225031a28f2b358f52d03a7217c7b524143

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL.String(), nil)
	if err != nil {
		return err
	}

	// We can't share these cookies with other requests, so make a new client.
	// Checking the return code on the first request so don't follow a redirect.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: c.c.Transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode == http.StatusOK {
		// the new API returns a 200 for a / request
		c.apiPath = "/proxy/network"
		c.loginPath = loginPathNew
		return nil
	}

	// The old version returns a "302" (to /manage) for a / request
	c.apiPath = "/"
	c.loginPath = loginPath
	return nil
}

func (c *Client) Login(ctx context.Context, user, pass string) error {
	if c.c == nil {
		c.c = &http.Client{}

		jar, _ := cookiejar.New(nil)
		c.c.Jar = jar
	}

	err := c.setAPIUrlStyle(ctx)
	if err != nil {
		return fmt.Errorf("unable to determine API URL style: %w", err)
	}

	var status struct {
		Meta struct {
			ServerVersion string `json:"server_version"`
			UUID          string `json:"uuid"`
		} `json:"meta"`
	}

	if c.apiKey == "" {
		err = c.do(ctx, "POST", c.loginPath, &struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: user,
			Password: pass,
		}, nil)
		if err != nil {
			return err
		}
	}

	err = c.do(ctx, "GET", "status", nil, &status)
	if err != nil {
		return err
	}

	if version := status.Meta.ServerVersion; version != "" {
		c.version = status.Meta.ServerVersion
		return nil
	}

	// newer version of 6.0 controller, use sysinfo to determine version
	// using default site since it must exist
	si, err := c.sysinfo(ctx, "default")
	if err != nil {
		return err
	}

	c.version = si.Version

	if c.version == "" {
		return errors.New("unable to determine controller version")
	}

	return nil
}

func (c *Client) do(
	ctx context.Context,
	method, relativeURL string,
	reqBody any,
	respBody any,
) error {
	// single threading requests, this is mostly to assist in CSRF token propagation
	if c.apiKey == "" {
		c.Lock()
		defer c.Unlock()
	}

	var (
		reqReader io.Reader
		err       error
		reqBytes  []byte
	)
	if reqBody != nil {
		reqBytes, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("unable to marshal JSON: %s %s %w", method, relativeURL, err)
		}
		reqReader = bytes.NewReader(reqBytes)
	}

	reqURL, err := url.Parse(relativeURL)
	if err != nil {
		return fmt.Errorf("unable to parse URL: %s %s %w", method, relativeURL, err)
	}
	if !strings.HasPrefix(relativeURL, "/") && !reqURL.IsAbs() {
		reqURL.Path = path.Join(c.apiPath, reqURL.Path)
	}

	url := c.baseURL.ResolveReference(reqURL)
	req, err := http.NewRequestWithContext(ctx, method, url.String(), reqReader)
	if err != nil {
		return fmt.Errorf("unable to create request: %s %s %w", method, relativeURL, err)
	}

	req.Header.Set("User-Agent", "terraform-provider-unifi/0.1")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	} else if c.csrf != "" {
		req.Header.Set("X-Csrf-Token", c.csrf)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("unable to perform request: %s %s %w", method, relativeURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{}
	}

	if c.apiKey == "" {
		if csrf := resp.Header.Get("X-Csrf-Token"); csrf != "" {
			c.csrf = resp.Header.Get("X-Csrf-Token")
		}
	}

	if resp.StatusCode != http.StatusOK {
		errBody := struct {
			Meta meta `json:"meta"`
			Data []struct {
				Meta meta `json:"meta"`
			} `json:"data"`
		}{}
		if err = json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
			return err
		}
		var apiErr error
		if len(errBody.Data) > 0 && errBody.Data[0].Meta.RC == "error" {
			// check first error in data, should we look for more than one?
			apiErr = errBody.Data[0].Meta.error()
		}
		if apiErr == nil {
			apiErr = errBody.Meta.error()
		}
		return fmt.Errorf("%w (%s) for %s %s", apiErr, resp.Status, method, url.String())
	}

	if respBody == nil || resp.ContentLength == 0 {
		return nil
	}

	// TODO: check rc in addition to status code?

	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		return fmt.Errorf("unable to decode body: %s %s %w", method, relativeURL, err)
	}

	return nil
}

type respData[T any] struct {
	Meta meta `json:"meta"`
	Data T    `json:"data"`
}

type meta struct {
	RC      string `json:"rc"`
	Message string `json:"msg"`
}

func (m *meta) error() error {
	if m.RC != "ok" {
		return &APIError{
			RC:      m.RC,
			Message: m.Message,
		}
	}

	return nil
}
