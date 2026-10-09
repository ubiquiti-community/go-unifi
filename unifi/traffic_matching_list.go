package unifi

import (
	"context"
	"fmt"
	"net/http"
)

// TrafficMatchingList is a reusable set of match criteria (IPv4 addresses,
// ports, etc.) exposed by the UniFi Network Integration API and referenced by
// firewall and traffic-management rules.
//
// The client methods take an IntegrationSiteID, not the legacy site name.
type TrafficMatchingList struct {
	ID    string                    `json:"id,omitempty"`
	Type  string                    `json:"type,omitempty"` // IPV4_ADDRESSES|PORTS
	Name  string                    `json:"name,omitempty"`
	Items []TrafficMatchingListItem `json:"items,omitempty"`
}

// TrafficMatchingListItem is a single entry within a TrafficMatchingList. Value
// is polymorphic: a string for address items (Type "IP_ADDRESS") and a number
// for port items (Type "PORT_NUMBER").
type TrafficMatchingListItem struct {
	Type  string `json:"type,omitempty"` // IP_ADDRESS|PORT_NUMBER
	Value any    `json:"value,omitempty"`
}

func (c *ApiClient) ListTrafficMatchingLists(
	ctx context.Context,
	siteID IntegrationSiteID,
) ([]TrafficMatchingList, error) {
	var respBody integrationPage[TrafficMatchingList]
	err := c.do(
		ctx,
		http.MethodGet,
		fmt.Sprintf("integration/v1/sites/%s/traffic-matching-lists", siteID),
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	return respBody.Data, nil
}

func (c *ApiClient) GetTrafficMatchingList(
	ctx context.Context,
	siteID IntegrationSiteID,
	id string,
) (*TrafficMatchingList, error) {
	var respBody TrafficMatchingList
	err := c.do(
		ctx,
		http.MethodGet,
		fmt.Sprintf("integration/v1/sites/%s/traffic-matching-lists/%s", siteID, id),
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	return &respBody, nil
}

func (c *ApiClient) CreateTrafficMatchingList(
	ctx context.Context,
	siteID IntegrationSiteID,
	d *TrafficMatchingList,
) (*TrafficMatchingList, error) {
	d.ID = ""
	var respBody TrafficMatchingList
	err := c.do(
		ctx,
		http.MethodPost,
		fmt.Sprintf("integration/v1/sites/%s/traffic-matching-lists", siteID),
		d,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	return &respBody, nil
}

func (c *ApiClient) UpdateTrafficMatchingList(
	ctx context.Context,
	siteID IntegrationSiteID,
	d *TrafficMatchingList,
) (*TrafficMatchingList, error) {
	// The API addresses the list by the URL path and rejects an "id" in the body.
	body := *d
	body.ID = ""

	var respBody TrafficMatchingList
	err := c.do(
		ctx,
		http.MethodPut,
		fmt.Sprintf("integration/v1/sites/%s/traffic-matching-lists/%s", siteID, d.ID),
		&body,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	return &respBody, nil
}

func (c *ApiClient) DeleteTrafficMatchingList(
	ctx context.Context,
	siteID IntegrationSiteID,
	id string,
) error {
	return c.do(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("integration/v1/sites/%s/traffic-matching-lists/%s", siteID, id),
		nil,
		nil,
	)
}
