// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
	_ types.Number
	_ strconv.NumError
)

type Client struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Anomalies                 *int64   `json:"anomalies,omitempty"`            // non-generated field
	AssocTime                 *int64   `json:"assoc_time,omitempty"`           // non-generated field
	Authorized                bool     `json:"authorized,omitempty"`           // non-generated field
	DevIdOverride             *int64   `json:"dev_id_override,omitempty"`      // non-generated field
	DisconnectTimestamp       *int64   `json:"disconnect_timestamp,omitempty"` // non-generated field
	EagerlyDiscovered         bool     `json:"eagerly_discovered,omitempty"`   // non-generated field
	FingerprintOverride       bool     `json:"fingerprint_override,omitempty"` // non-generated field
	FirstSeen                 *int64   `json:"first_seen,omitempty"`           // non-generated field
	GWMAC                     string   `json:"gw_mac,omitempty"`               // non-generated field
	GWVLAN                    *int64   `json:"gw_vlan,omitempty"`              // non-generated field
	HostnameSource            string   `json:"hostname_source,omitempty"`      // non-generated field
	IP                        string   `json:"ip,omitempty"`                   // non-generated field
	IPv6Addresses             []string `json:"ipv6_addresses,omitempty"`
	IsGuest                   bool     `json:"is_guest,omitempty"`                     // non-generated field
	IsGuestByUGW              bool     `json:"_is_guest_by_ugw,omitempty"`             // non-generated field
	IsGuestByUSW              bool     `json:"_is_guest_by_usw,omitempty"`             // non-generated field
	IsWired                   bool     `json:"is_wired,omitempty"`                     // non-generated field
	Last1xIdentity            string   `json:"last_1x_identity,omitempty"`             // non-generated field
	LastConnectionNetworkID   string   `json:"last_connection_network_id,omitempty"`   // non-generated field
	LastConnectionNetworkName string   `json:"last_connection_network_name,omitempty"` // non-generated field
	LastIP                    string   `json:"last_ip,omitempty"`
	LastIPv6                  []string `json:"last_ipv6,omitempty"`
	LastReachableByGW         *int64   `json:"_last_reachable_by_gw,omitempty"`      // non-generated field
	LastSeenByUGW             *int64   `json:"_last_seen_by_ugw,omitempty"`          // non-generated field
	LastSeenByUSW             *int64   `json:"_last_seen_by_usw,omitempty"`          // non-generated field
	LastUplinkMAC             string   `json:"last_uplink_mac,omitempty"`            // non-generated field
	LastUplinkName            string   `json:"last_uplink_name,omitempty"`           // non-generated field
	LastUplinkRemotePort      *int64   `json:"last_uplink_remote_port,omitempty"`    // non-generated field
	LatestAssocTime           *int64   `json:"latest_assoc_time,omitempty"`          // non-generated field
	Network                   string   `json:"network,omitempty"`                    // non-generated field
	Noted                     bool     `json:"noted,omitempty"`                      // non-generated field
	OUI                       string   `json:"oui,omitempty"`                        // non-generated field
	QOSPolicyApplied          bool     `json:"qos_policy_applied,omitempty"`         // non-generated field
	Satisfaction              *int64   `json:"satisfaction,omitempty"`               // non-generated field
	SwDepth                   *int64   `json:"sw_depth,omitempty"`                   // non-generated field
	SwMAC                     string   `json:"sw_mac,omitempty"`                     // non-generated field
	SwPort                    *int64   `json:"sw_port,omitempty"`                    // non-generated field
	TxRetries                 *int64   `json:"tx_retries,omitempty"`                 // non-generated field
	Uptime                    *int64   `json:"uptime,omitempty"`                     // non-generated field
	UptimeByUGW               *int64   `json:"_uptime_by_ugw,omitempty"`             // non-generated field
	UptimeByUSW               *int64   `json:"_uptime_by_usw,omitempty"`             // non-generated field
	UserID                    string   `json:"user_id,omitempty"`                    // non-generated field
	VLAN                      *int64   `json:"vlan,omitempty"`                       // non-generated field
	WiFiTxAttempts            *int64   `json:"wifi_tx_attempts,omitempty"`           // non-generated field
	WiFiTxDropped             *int64   `json:"wifi_tx_dropped,omitempty"`            // non-generated field
	WiFiTxRetriesPercentage   float64  `json:"wifi_tx_retries_percentage,omitempty"` // non-generated field
	WiredRateMbps             *int64   `json:"wired_rate_mbps,omitempty"`            // non-generated field
	WiredRxBytes              *int64   `json:"wired-rx_bytes,omitempty"`             // non-generated field
	WiredRxBytesR             float64  `json:"wired-rx_bytes-r,omitempty"`           // non-generated field
	WiredRxPackets            *int64   `json:"wired-rx_packets,omitempty"`           // non-generated field
	WiredTxBytes              *int64   `json:"wired-tx_bytes,omitempty"`             // non-generated field
	WiredTxBytesR             float64  `json:"wired-tx_bytes-r,omitempty"`           // non-generated field
	WiredTxPackets            *int64   `json:"wired-tx_packets,omitempty"`           // non-generated field

	Blocked                       bool     `json:"blocked,omitempty"`
	FixedApEnabled                bool     `json:"fixed_ap_enabled"`
	FixedApMAC                    string   `json:"fixed_ap_mac,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	FixedIP                       string   `json:"fixed_ip,omitempty"`
	Hostname                      string   `json:"hostname,omitempty"`
	LastSeen                      *int64   `json:"last_seen,omitempty"`
	LocalDNSRecord                string   `json:"local_dns_record,omitempty"`
	LocalDNSRecordEnabled         bool     `json:"local_dns_record_enabled"`
	MAC                           string   `json:"mac,omitempty"` // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	Name                          string   `json:"name,omitempty"`
	NetworkID                     string   `json:"network_id,omitempty"`
	NetworkMembersGroupIDs        []string `json:"network_members_group_ids,omitempty"`
	Note                          string   `json:"note,omitempty"`
	UseFixedIP                    bool     `json:"use_fixedip"`
	UserGroupID                   string   `json:"usergroup_id,omitempty"`
	VirtualNetworkOverrideEnabled bool     `json:"virtual_network_override_enabled"`
	VirtualNetworkOverrideID      string   `json:"virtual_network_override_id,omitempty"`
}

func (dst *Client) UnmarshalJSON(b []byte) error {
	type Alias Client
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}

	return nil
}

func (c *ApiClient) listClient(ctx context.Context, site string) ([]Client, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Client `json:"data"`
	}

	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("api/s/%s/rest/user", site),
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	return respBody.Data, nil
}

func (c *ApiClient) getClient(
	ctx context.Context,
	site string,
	id string,
) (*Client, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Client `json:"data"`
	}
	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("api/s/%s/rest/user/%s", site, id),
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *ApiClient) deleteClient(
	ctx context.Context,
	site string,
	id string,
) error {
	err := c.do(
		ctx,
		"DELETE",
		fmt.Sprintf("api/s/%s/rest/user/%s", site, id),
		struct{}{},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) createClient(
	ctx context.Context,
	site string,
	d *Client,
) (*Client, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Client `json:"data"`
	}

	err := c.do(
		ctx,
		"POST",
		fmt.Sprintf("api/s/%s/rest/user", site),
		d,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}

func (c *ApiClient) updateClient(
	ctx context.Context,
	site string,
	d *Client,
) (*Client, error) {
	var respBody struct {
		Meta meta     `json:"meta"`
		Data []Client `json:"data"`
	}
	err := c.do(
		ctx,
		"PUT",
		fmt.Sprintf("api/s/%s/rest/user/%s", site, d.ID),
		d,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}
