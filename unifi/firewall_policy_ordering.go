package unifi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// ErrFirewallPolicyOrderingIncomplete is returned (wrapped, together with the
// controller's *APIError) by UpdateFirewallPolicyOrdering when the controller
// rejects an ordering that does not list every user-defined policy of the zone
// pair (code api.firewall.policy.ordering.not-all-firewall-policies-specified).
var ErrFirewallPolicyOrderingIncomplete = errors.New(
	"firewall policy ordering must list every user-defined policy of the zone pair",
)

const firewallPolicyOrderingIncompleteCode = "api.firewall.policy.ordering.not-all-firewall-policies-specified"

// FirewallPolicyOrdering is the user-defined firewall policy order for one
// source/destination zone pair, as exposed by the UniFi Network Integration API
// (GET/PUT /integration/v1/sites/{siteId}/firewall/policies/ordering).
//
// Policy and zone ids are Integration API UUIDs, not the legacy ids returned by
// the v2 firewall-policies endpoints. System-defined policies are not listed:
// they sit between the two segments, so BeforeSystemDefined holds the
// user-defined policies evaluated before them and AfterSystemDefined those
// evaluated after them. Return policies derived from allowReturnTraffic are not
// listed either.
//
// Network 10.6.109 has been seen omitting the id of a user-defined policy from
// the policy listing and reporting it as "" in the ordering; such an ordering
// cannot be sent back as-is.
type FirewallPolicyOrdering struct {
	OrderedFirewallPolicyIDs FirewallPolicyOrderedIDs `json:"orderedFirewallPolicyIds"`
}

// FirewallPolicyOrderedIDs holds the ordered user-defined policy ids on either
// side of the system-defined policies of a zone pair.
type FirewallPolicyOrderedIDs struct {
	BeforeSystemDefined []string `json:"beforeSystemDefined"`
	AfterSystemDefined  []string `json:"afterSystemDefined"`
}

// PolicyIDs returns every policy id in the ordering, BeforeSystemDefined first.
func (o *FirewallPolicyOrdering) PolicyIDs() []string {
	ids := make(
		[]string,
		0,
		len(o.OrderedFirewallPolicyIDs.BeforeSystemDefined)+
			len(o.OrderedFirewallPolicyIDs.AfterSystemDefined),
	)
	ids = append(ids, o.OrderedFirewallPolicyIDs.BeforeSystemDefined...)
	return append(ids, o.OrderedFirewallPolicyIDs.AfterSystemDefined...)
}

func firewallPolicyOrderingQuery(sourceZoneID, destinationZoneID string) map[string]string {
	return map[string]string{
		"sourceFirewallZoneId":      sourceZoneID,
		"destinationFirewallZoneId": destinationZoneID,
	}
}

// GetFirewallPolicyOrdering returns the user-defined policy order of the zone
// pair identified by its Integration API zone UUIDs.
func (c *ApiClient) GetFirewallPolicyOrdering(
	ctx context.Context,
	siteID IntegrationSiteID,
	sourceZoneID, destinationZoneID string,
) (*FirewallPolicyOrdering, error) {
	var respBody FirewallPolicyOrdering
	err := c.do(
		ctx,
		http.MethodGet,
		fmt.Sprintf("integration/v1/sites/%s/firewall/policies/ordering", siteID),
		nil,
		&respBody,
		firewallPolicyOrderingQuery(sourceZoneID, destinationZoneID),
	)
	if err != nil {
		return nil, err
	}
	return &respBody, nil
}

// UpdateFirewallPolicyOrdering replaces the user-defined policy order of the
// zone pair and returns the order the controller stored.
//
// The ordering must list EVERY user-defined policy of the zone pair, split
// between BeforeSystemDefined and AfterSystemDefined; the controller rejects a
// subset with HTTP 400, reported as an error wrapping both
// ErrFirewallPolicyOrderingIncomplete and the controller's *APIError. To
// reorder, read the current ordering with GetFirewallPolicyOrdering and
// permute its ids rather than building the list from scratch. Derived return
// policies (allowReturnTraffic) are never part of the ordering.
func (c *ApiClient) UpdateFirewallPolicyOrdering(
	ctx context.Context,
	siteID IntegrationSiteID,
	sourceZoneID, destinationZoneID string,
	d *FirewallPolicyOrdering,
) (*FirewallPolicyOrdering, error) {
	// Both lists are required by the API: send [] rather than null.
	body := *d
	if body.OrderedFirewallPolicyIDs.BeforeSystemDefined == nil {
		body.OrderedFirewallPolicyIDs.BeforeSystemDefined = []string{}
	}
	if body.OrderedFirewallPolicyIDs.AfterSystemDefined == nil {
		body.OrderedFirewallPolicyIDs.AfterSystemDefined = []string{}
	}

	var respBody FirewallPolicyOrdering
	err := c.do(
		ctx,
		http.MethodPut,
		fmt.Sprintf("integration/v1/sites/%s/firewall/policies/ordering", siteID),
		&body,
		&respBody,
		firewallPolicyOrderingQuery(sourceZoneID, destinationZoneID),
	)
	if err != nil {
		apiErr := &APIError{}
		if errors.As(err, &apiErr) && apiErr.RC == firewallPolicyOrderingIncompleteCode {
			return nil, fmt.Errorf("%w: %w", ErrFirewallPolicyOrderingIncomplete, err)
		}
		return nil, err
	}
	return &respBody, nil
}
