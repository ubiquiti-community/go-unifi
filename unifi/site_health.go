package unifi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
)

// SiteHealthGwSystemStats holds gateway system statistics reported by the
// health endpoint.
type SiteHealthGwSystemStats struct {
	CPU    types.Number `json:"cpu"`
	Mem    types.Number `json:"mem"`
	Uptime types.Number `json:"uptime"`
}

// SiteHealth is one per-subsystem entry from the site health endpoint
// (subsystems include wan, lan, wlan, www, vpn).
type SiteHealth struct {
	Subsystem string `json:"subsystem"`
	Status    string `json:"status"`

	NumUser  types.Number `json:"num_user,omitempty"`
	NumGuest types.Number `json:"num_guest,omitempty"`
	NumIot   types.Number `json:"num_iot,omitempty"`
	TxBytesR types.Number `json:"tx_bytes-r,omitempty"`
	RxBytesR types.Number `json:"rx_bytes-r,omitempty"`

	NumAp           types.Number `json:"num_ap,omitempty"`
	NumAdopted      types.Number `json:"num_adopted,omitempty"`
	NumDisabled     types.Number `json:"num_disabled,omitempty"`
	NumDisconnected types.Number `json:"num_disconnected,omitempty"`
	NumPending      types.Number `json:"num_pending,omitempty"`
	NumGw           types.Number `json:"num_gw,omitempty"`
	NumSw           types.Number `json:"num_sw,omitempty"`
	NumSta          types.Number `json:"num_sta,omitempty"`

	WanIP         string                  `json:"wan_ip,omitempty"`
	Gateways      []string                `json:"gateways,omitempty"`
	Netmask       string                  `json:"netmask,omitempty"`
	Nameservers   []string                `json:"nameservers,omitempty"`
	GwMac         string                  `json:"gw_mac,omitempty"`
	GwName        string                  `json:"gw_name,omitempty"`
	GwSystemStats SiteHealthGwSystemStats `json:"gw_system-stats,omitempty"`
	GwVersion     string                  `json:"gw_version,omitempty"`
	LanIP         string                  `json:"lan_ip,omitempty"`

	Latency          types.Number `json:"latency,omitempty"`
	Uptime           types.Number `json:"uptime,omitempty"`
	Drops            types.Number `json:"drops,omitempty"`
	XputUp           types.Number `json:"xput_up,omitempty"`
	XputDown         types.Number `json:"xput_down,omitempty"`
	SpeedtestStatus  string       `json:"speedtest_status,omitempty"`
	SpeedtestLastrun types.Number `json:"speedtest_lastrun,omitempty"`
	SpeedtestPing    types.Number `json:"speedtest_ping,omitempty"`

	RemoteUserEnabled     bool         `json:"remote_user_enabled,omitempty"`
	RemoteUserNumActive   types.Number `json:"remote_user_num_active,omitempty"`
	RemoteUserNumInactive types.Number `json:"remote_user_num_inactive,omitempty"`
	RemoteUserRxBytes     types.Number `json:"remote_user_rx_bytes,omitempty"`
	RemoteUserTxBytes     types.Number `json:"remote_user_tx_bytes,omitempty"`
	RemoteUserRxPackets   types.Number `json:"remote_user_rx_packets,omitempty"`
	RemoteUserTxPackets   types.Number `json:"remote_user_tx_packets,omitempty"`

	SiteToSiteEnabled     bool         `json:"site_to_site_enabled,omitempty"`
	SiteToSiteNumActive   types.Number `json:"site_to_site_num_active,omitempty"`
	SiteToSiteNumInactive types.Number `json:"site_to_site_num_inactive,omitempty"`
	SiteToSiteRxBytes     types.Number `json:"site_to_site_rx_bytes,omitempty"`
	SiteToSiteTxBytes     types.Number `json:"site_to_site_tx_bytes,omitempty"`
	SiteToSiteRxPackets   types.Number `json:"site_to_site_rx_packets,omitempty"`
	SiteToSiteTxPackets   types.Number `json:"site_to_site_tx_packets,omitempty"`
}

// GetHealth returns the per-subsystem health metrics for a site.
func (c *ApiClient) GetHealth(ctx context.Context, site string) ([]SiteHealth, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []SiteHealth `json:"data"`
	}

	err := c.do(ctx, http.MethodGet, fmt.Sprintf("api/s/%s/stat/health", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if err := respBody.Meta.error(); err != nil {
		return nil, err
	}

	return respBody.Data, nil
}
