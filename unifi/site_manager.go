package unifi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

type Host struct {
	ID         string `json:"id"`
	HardwareID string `json:"hardwareId"`
	Owner      bool   `json:"owner"`
}

type ListHostsResponse struct {
	Data      []Host `json:"data"`
	NextToken string `json:"nextToken"`
}

func FindHost(client *retryablehttp.Client, apiKey string, matchFunc func(Host) bool) (*Host, error) {
	baseURL := "https://api.ui.com/v1/hosts"
	nextToken := ""
	pageSize := "50" // Adjust as needed

	for {
		// Build the URL with query parameters
		reqURL, err := url.Parse(baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %w", err)
		}

		q := reqURL.Query()
		q.Set("pageSize", pageSize)
		if nextToken != "" {
			q.Set("nextToken", nextToken)
		}
		reqURL.RawQuery = q.Encode()

		// Setup the request
		req, err := retryablehttp.NewRequest("GET", reqURL.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-API-KEY", apiKey)

		// Execute the request
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		// Decode the JSON response
		var apiResponse ListHostsResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode json: %w", err)
		}
		resp.Body.Close()

		// Evaluate each host against our custom matchFunc
		for _, host := range apiResponse.Data {
			if matchFunc(host) {
				return &host, nil // Match found! Return the host immediately.
			}
		}

		// Pagination check
		if apiResponse.NextToken == "" {
			break
		}
		nextToken = apiResponse.NextToken
	}

	return nil, fmt.Errorf("no host found matching criteria across all pages")
}

func GetFirstOwnedHostID(client *retryablehttp.Client, apiKey string) (string, error) {
	host, err := FindHost(client, apiKey, func(h Host) bool {
		return h.Owner // Return true if the user is the owner
	})
	if err != nil {
		return "", err
	}
	return host.ID, nil
}

// 4. Implementation B: Get Host by Hardware ID
func GetHostIDByHardwareID(client *retryablehttp.Client, apiKey, targetHardwareID string) (string, error) {
	host, err := FindHost(client, apiKey, func(h Host) bool {
		return h.HardwareID == targetHardwareID // Return true if the IDs match
	})
	if err != nil {
		return "", err
	}
	return host.ID, nil
}
