package unifi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

// ErrAmbiguousIntegrationSite is returned (wrapped) by ResolveIntegrationSiteID when
// more than one site matches the given display name. Other failures, such as a
// failed request, are not wrapped with it.
var ErrAmbiguousIntegrationSite = errors.New("ambiguous Integration API site")

// IntegrationSiteID is the site UUID used by the Network Integration API
// (/proxy/network/integration/v1/sites/{id}/...). It is not the legacy site
// name (e.g. "default") taken by the other client methods. Obtain one with
// ResolveIntegrationSiteID.
type IntegrationSiteID string

// IntegrationSite is a site as listed by the Integration API. InternalReference
// is the legacy site name (e.g. "default") used by the rest of this client.
type IntegrationSite struct {
	ID                IntegrationSiteID `json:"id"`
	InternalReference string            `json:"internalReference"`
	Name              string            `json:"name"`
}

// integrationPage is the paginated response envelope returned by the UniFi
// Network Integration API (/proxy/network/integration/v1/...). Unlike the legacy
// REST API, which wraps results in a {meta, data} envelope, the Integration API
// returns offset/limit pagination metadata alongside the data.
type integrationPage[T any] struct {
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	Count      int `json:"count"`
	TotalCount int `json:"totalCount"`
	Data       []T `json:"data"`
}

const integrationPageLimit = 200

var integrationUUIDPattern = regexp.MustCompile(
	`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
)

// ListIntegrationSites returns every site known to the Integration API,
// following offset/limit pagination.
func (c *ApiClient) ListIntegrationSites(ctx context.Context) ([]IntegrationSite, error) {
	var sites []IntegrationSite
	for offset := 0; ; {
		var page integrationPage[IntegrationSite]
		err := c.do(
			ctx,
			http.MethodGet,
			"integration/v1/sites",
			nil,
			&page,
			map[string]string{
				"offset": strconv.Itoa(offset),
				"limit":  strconv.Itoa(integrationPageLimit),
			},
		)
		if err != nil {
			return nil, err
		}

		sites = append(sites, page.Data...)
		offset += len(page.Data)
		if len(page.Data) == 0 || offset >= page.TotalCount {
			return sites, nil
		}
	}
}

// ResolveIntegrationSiteID converts a site identifier into the UUID the
// Integration API requires. A value that is already a UUID is returned as-is
// without a request. Otherwise it is matched against each site's legacy name
// (internalReference, e.g. "default") and then its display name.
func (c *ApiClient) ResolveIntegrationSiteID(
	ctx context.Context,
	site string,
) (IntegrationSiteID, error) {
	if integrationUUIDPattern.MatchString(site) {
		return IntegrationSiteID(site), nil
	}

	sites, err := c.ListIntegrationSites(ctx)
	if err != nil {
		return "", err
	}

	for _, s := range sites {
		if s.InternalReference == site {
			return s.ID, nil
		}
	}

	var byName []IntegrationSite
	for _, s := range sites {
		if s.Name == site {
			byName = append(byName, s)
		}
	}
	switch len(byName) {
	case 1:
		return byName[0].ID, nil
	case 0:
		return "", &NotFoundError{Type: "IntegrationSite", Attr: "name", Value: site}
	default:
		return "", fmt.Errorf(
			"%w: site %q is ambiguous: %d Integration API sites share that name; use the site UUID",
			ErrAmbiguousIntegrationSite,
			site,
			len(byName),
		)
	}
}
