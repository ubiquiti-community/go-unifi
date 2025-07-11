package unifi

import (
	"context"
	"fmt"
)

// just to fix compile issues with the import.
var (
	_ fmt.Formatter
	_ context.Context
)

type Fingerprint struct {
	ComputedDevId  int  `json:"computed_dev_id,omitempty"`
	ComputedEngine int  `json:"computed_engine,omitempty"`
	Confidence     int  `json:"confidence,omitempty"`
	DevCat         int  `json:"dev_cat,omitempty"`
	DevFamily      int  `json:"dev_family,omitempty"`
	DevId          int  `json:"dev_id,omitempty"`
	DevIdOverride  int  `json:"dev_id_override,omitempty"`
	DevVendor      int  `json:"dev_vendor,omitempty"`
	HasOverride    bool `json:"has_override,omitempty"`
	OsName         int  `json:"os_name,omitempty"`
}

type UnifiDeviceInfo struct {
	IconFilename      string  `json:"icon_filename,omitempty"`
	IconResolutions   [][]int `json:"icon_resolutions,omitempty"`
	ViewInApplication bool    `json:"view_in_application,omitempty"`
}

type DetailedStates struct {
	UplinkNearPowerLimit bool `json:"uplink_near_power_limit,omitempty"`
}

type ClientInfo struct {
	Anomalies                           int             `json:"anomalies,omitempty"`
	ApMac                               string          `json:"ap_mac,omitempty"`
	AssocTime                           int             `json:"assoc_time,omitempty"`
	Authorized                          bool            `json:"authorized,omitempty"`
	Blocked                             bool            `json:"blocked,omitempty"`
	Bssid                               string          `json:"bssid,omitempty"`
	Ccq                                 int             `json:"ccq,omitempty"`
	Channel                             int             `json:"channel,omitempty"`
	ChannelWidth                        string          `json:"channel_width,omitempty"`
	DetailedStates                      DetailedStates  `json:"detailed_states,omitempty"`
	DhcpendTime                         int             `json:"dhcpend_time,omitempty"`
	DisplayName                         string          `json:"display_name,omitempty"`
	Essid                               string          `json:"essid,omitempty"`
	Fingerprint                         Fingerprint     `json:"fingerprint,omitempty"`
	FirstSeen                           int             `json:"first_seen,omitempty"`
	FixedApEnabled                      bool            `json:"fixed_ap_enabled,omitempty"`
	FixedIP                             string          `json:"fixed_ip,omitempty"`
	GwMac                               string          `json:"gw_mac,omitempty"`
	Hostname                            string          `json:"hostname,omitempty"`
	Id                                  string          `json:"id,omitempty"`
	Idletime                            int             `json:"idletime,omitempty"`
	IP                                  string          `json:"ip,omitempty"`
	Ipv4LeaseExpirationTimestampSeconds int             `json:"ipv4_lease_expiration_timestamp_seconds,omitempty"`
	Ipv6Address                         []string        `json:"ipv6_address,omitempty"`
	IsAllowedInVisualProgramming        bool            `json:"is_allowed_in_visual_programming,omitempty"`
	IsGuest                             bool            `json:"is_guest,omitempty"`
	IsMlo                               bool            `json:"is_mlo,omitempty"`
	IsWired                             bool            `json:"is_wired,omitempty"`
	LastConnectionNetworkId             string          `json:"last_connection_network_id,omitempty"`
	LastConnectionNetworkName           string          `json:"last_connection_network_name,omitempty"`
	LastIP                              string          `json:"last_ip,omitempty"`
	LastIpv6                            []string        `json:"last_ipv6,omitempty"`
	LastRadio                           string          `json:"last_radio,omitempty"`
	LastSeen                            int             `json:"last_seen,omitempty"`
	LastUplinkMac                       string          `json:"last_uplink_mac,omitempty"`
	LastUplinkName                      string          `json:"last_uplink_name,omitempty"`
	LatestAssocTime                     int             `json:"latest_assoc_time,omitempty"`
	LocalDNSRecord                      string          `json:"local_dns_record,omitempty"`
	LocalDNSRecordEnabled               bool            `json:"local_dns_record_enabled,omitempty"`
	Mac                                 string          `json:"mac,omitempty"`
	Mimo                                string          `json:"mimo,omitempty"`
	ModelName                           string          `json:"model_name,omitempty"`
	Name                                string          `json:"name,omitempty"`
	NetworkId                           string          `json:"network_id,omitempty"`
	NetworkName                         string          `json:"network_name,omitempty"`
	Noise                               int             `json:"noise,omitempty"`
	Noted                               bool            `json:"noted,omitempty"`
	Oui                                 string          `json:"oui,omitempty"`
	PowersaveEnabled                    bool            `json:"powersave_enabled,omitempty"`
	Radio                               string          `json:"radio,omitempty"`
	RadioName                           string          `json:"radio_name,omitempty"`
	RadioProto                          string          `json:"radio_proto,omitempty"`
	RateImbalance                       int             `json:"rate_imbalance,omitempty"`
	Rssi                                int             `json:"rssi,omitempty"`
	RxBytes                             int             `json:"rx_bytes,omitempty"`
	RxBytesR                            int             `json:"rx_bytes-r,omitempty"`
	RxPackets                           int             `json:"rx_packets,omitempty"`
	RxRate                              int             `json:"rx_rate,omitempty"`
	Signal                              int             `json:"signal,omitempty"`
	SiteId                              string          `json:"site_id,omitempty"`
	Status                              string          `json:"status,omitempty"`
	SwPort                              int             `json:"sw_port,omitempty"`
	Tags                                []string        `json:"tags,omitempty"`
	TxBytes                             int             `json:"tx_bytes,omitempty"`
	TxBytesR                            int             `json:"tx_bytes-r,omitempty"`
	TxMcsIndex                          int             `json:"tx_mcs_index,omitempty"`
	TxPackets                           int             `json:"tx_packets,omitempty"`
	TxRate                              int             `json:"tx_rate,omitempty"`
	Type                                string          `json:"type,omitempty"`
	UnifiDevice                         bool            `json:"unifi_device,omitempty"`
	UnifiDeviceInfo                     UnifiDeviceInfo `json:"unifi_device_info,omitempty"`
	UplinkMac                           string          `json:"uplink_mac,omitempty"`
	Uptime                              int             `json:"uptime,omitempty"`
	UseFixedip                          bool            `json:"use_fixedip,omitempty"`
	UsergroupId                         string          `json:"usergroup_id,omitempty"`
	UserId                              string          `json:"user_id,omitempty"`
	VirtualNetworkOverrideEnabled       bool            `json:"virtual_network_override_enabled,omitempty"`
	VirtualNetworkOverrideId            string          `json:"virtual_network_override_id,omitempty"`
	WifiExperienceAverage               int             `json:"wifi_experience_average,omitempty"`
	WifiExperienceScore                 int             `json:"wifi_experience_score,omitempty"`
	WifiTxAttempts                      int             `json:"wifi_tx_attempts,omitempty"`
	WifiTxRetriesPercentage             float64         `json:"wifi_tx_retries_percentage,omitempty"`
	WiredRateMbps                       int             `json:"wired_rate_mbps,omitempty"`
	WlanconfId                          string          `json:"wlanconf_id,omitempty"`
}

type ClientList []ClientInfo

func (c *Client) ListClientsActive(ctx context.Context, site string) (ClientList, error) {
	var respBody []ClientInfo

	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("v2/api/site/%s/clients/active?includeUnifiDevices=true", site),
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *Client) GetClientLocal(ctx context.Context, site string, mac string) (*ClientInfo, error) {
	var respBody ClientInfo

	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("v2/api/site/%s/clients/local/%s?includeUnifiDevices=true", site, mac),
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}
