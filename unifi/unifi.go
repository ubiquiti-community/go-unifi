package unifi //

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	network "github.com/ubiquiti-community/go-unifi/client/network"
	protect "github.com/ubiquiti-community/go-unifi/client/protect"

	"github.com/hashicorp/go-retryablehttp"
)

//go:generate go tool oapi-codegen -generate models,client -o ../client/protect/client.gen.go -package protect ../client/protect/openapi.json
//go:generate go tool oapi-codegen -generate models,client -o ../client/network/client.gen.go -package network ../client/network/openapi.json

// Config holds all configuration for creating a new ApiClient.
type Config struct {
	BaseURL        string
	APIKey         string
	AllowInsecure  bool
	CloudConnector bool
	HardwareID     string
	Logger         any
	TimeoutSeconds *int
	RetryMax       *int
}

// New creates a fully initialized ApiClient from the provided configuration.
// For cloud connector mode, set CloudConnector=true and optionally HardwareID.
// For direct connection, provide BaseURL and either APIKey or Username/Password.
func New(ctx context.Context, cfg *Config) (*ApiClient, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("API key is required")
	}
	if cfg.BaseURL == "" && !cfg.CloudConnector {
		return nil, errors.New("BaseURL is required for direct connection")
	}

	c := retryablehttp.NewClient()
	timeoutSeconds := 30
	if cfg.TimeoutSeconds != nil {
		timeoutSeconds = *cfg.TimeoutSeconds
	}
	c.HTTPClient.Timeout = time.Duration(timeoutSeconds) * time.Second

	if cfg.Logger != nil {
		c.Logger = cfg.Logger
	} else {
		c.Logger = nil
	}

	if cfg.RetryMax != nil {
		c.RetryMax = *cfg.RetryMax
	}

	if cfg.AllowInsecure {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: time.Duration(timeoutSeconds) * time.Second,
			}).DialContext,
		}
		c.HTTPClient.Transport = transport
	}

	jar, _ := cookiejar.New(nil)
	c.HTTPClient.Jar = jar

	apiKeyEditor := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-API-KEY", cfg.APIKey)
		req.Header.Set("User-Agent", "terraform-provider-unifi/0.1")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		return nil
	}

	if cfg.CloudConnector {
		var hostID string
		var err error
		if cfg.HardwareID != "" {
			hostID, err = GetHostIDByHardwareID(c, cfg.APIKey, cfg.HardwareID)
			if err != nil {
				return nil, fmt.Errorf("unable to find host with hardware ID %s: %w", cfg.HardwareID, err)
			}
		} else {
			hostID, err = GetFirstOwnedHostID(c, cfg.APIKey)
			if err != nil {
				return nil, fmt.Errorf("unable to find first owned host: %w", err)
			}
		}
		cfg.BaseURL, err = url.JoinPath("https://api.ui.com/v1/connector/consoles/", hostID)
		if err != nil {
			return nil, fmt.Errorf("unable to construct base URL: %w", err)
		}
	}

	var err error
	cfg.BaseURL, err = url.JoinPath(cfg.BaseURL, "/proxy/")
	if err != nil {
		return nil, fmt.Errorf("unable to construct base URL: %w", err)
	}

	networkURL, err := url.JoinPath(cfg.BaseURL, "/network/integration")
	if err != nil {
		return nil, fmt.Errorf("unable to construct network URL: %w", err)
	}
	networkClient, err := network.NewClientWithResponses(networkURL, network.WithHTTPClient(c.StandardClient()), network.WithRequestEditorFn(apiKeyEditor))
	if err != nil {
		return nil, fmt.Errorf("unable to create API client: %w", err)
	}

	protectURL, err := url.JoinPath(cfg.BaseURL, "/protect/integration")
	if err != nil {
		return nil, fmt.Errorf("unable to construct protect URL: %w", err)
	}
	protectClient, err := protect.NewClientWithResponses(protectURL, protect.WithHTTPClient(c.StandardClient()), protect.WithRequestEditorFn(apiKeyEditor))
	if err != nil {
		return nil, fmt.Errorf("unable to create API client: %w", err)
	}

	client := &ApiClient{
		ctx:     ctx,
		network: nil,
		protect: nil,
	}

	networkResponse, err := networkClient.GetInfoWithResponse(ctx)
	if err != nil || networkResponse.StatusCode() == http.StatusOK {
		client.network = &NetworkAPI{
			client:  *networkClient,
			version: networkResponse.JSON200.ApplicationVersion,
		}
	}

	protectResponse, err := protectClient.GetV1MetaInfoWithResponse(ctx)
	if err != nil || protectResponse.StatusCode() == http.StatusOK {
		client.protect = &ProtectAPI{
			client:  *protectClient,
			version: protectResponse.JSON200.ApplicationVersion,
		}
	}

	return client, nil
}

type NetworkAPI struct {
	client  network.ClientWithResponses
	version string
}

type ProtectAPI struct {
	client  protect.ClientWithResponses
	version string
}

type ApiClient struct {
	ctx     context.Context
	network *NetworkAPI
	protect *ProtectAPI
}
