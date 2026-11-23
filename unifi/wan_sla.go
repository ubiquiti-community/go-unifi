package unifi

import (
	"context"
	"fmt"
	"net/http"
)

// WAN SLA profiles (Network 11: Settings > Profiles > WAN SLA) decide whether a WAN is usable. Each profile has one or
// more verification servers (ICMP or DNS) with packet loss, latency and jitter thresholds; when enough servers breach
// them, the WAN is taken out of load balancing or failover. A WAN network references a profile by ID in its wan_sla
// field (Network.WANSla).
//
// This is a v2 API object (v2/api/site/<site>/wan-slas), so it is coded by hand like APGroup. The controller has a
// per-id GET, PUT and DELETE; DELETE returns an empty body and a missing id is a 404 (api.err.WanSlaNotFound).

// WANSlaMonitorHealthType values: when the profile counts as down.
const (
	WANSlaMonitorHealthAll      = "ALL"      // every verification server breaches
	WANSlaMonitorHealthAny      = "ANY"      // any one server breaches
	WANSlaMonitorHealthMajority = "MAJORITY" // more than half of the servers breach
	WANSlaMonitorHealthCustom   = "CUSTOM"   // MonitorDownThreshold servers breach
)

// WANSlaMonitor types.
const (
	WANSlaMonitorTypeICMP = "ICMP"
	WANSlaMonitorTypeDNS  = "DNS"
)

type WANSla struct {
	ID string `json:"_id,omitempty"`

	Name string `json:"name"`
	// ALL, ANY, MAJORITY or CUSTOM (see the constants above).
	MonitorHealthType string `json:"monitor_health_type"`
	// Only with CUSTOM: how many servers must breach.
	MonitorDownThreshold *int64          `json:"monitor_down_threshold,omitempty"`
	Monitors             []WANSlaMonitor `json:"monitors"`
}

type WANSlaMonitor struct {
	// Assigned by the controller (1, 2, ...). Send it back on update to keep the monitor's identity.
	MonitorID *int64 `json:"monitor_id,omitempty"`
	// ICMP or DNS.
	Type string `json:"type"`
	// ICMP: host name or IPv4 address; DNS: the IPv4 address of the DNS server.
	Target string `json:"target"`
	// DNS only: the name resolved against Target.
	Domain string `json:"domain,omitempty"`
	// Seconds between probes (1-3600).
	Interval int64 `json:"interval"`
	// Seconds over which probe results are averaged (1-65535); must exceed 2 x Interval + latency threshold (s).
	TimePeriod int64       `json:"time_period"`
	Alert      WANSlaAlert `json:"alert"`
}

type WANSlaAlert struct {
	// any, all, any-two, loss-or-all-others, latency-or-all-others, jitter-or-all-others, loss-plus-one,
	// latency-plus-one, jitter-plus-one.
	ThresholdPolicy string `json:"threshold_policy"`
	// Percent (1-100). Nil disables the target.
	LossThreshold *int64 `json:"loss_threshold,omitempty"`
	// Milliseconds (5-10000). Nil disables the target.
	LatencyThreshold *int64 `json:"latency_threshold,omitempty"`
	// Milliseconds (3-1500). Nil disables the target.
	JitterThreshold *int64 `json:"jitter_threshold,omitempty"`
}

func (c *ApiClient) ListWANSla(ctx context.Context, site string) ([]WANSla, error) {
	var respBody []WANSla

	err := c.do(ctx, http.MethodGet, fmt.Sprintf("v2/api/site/%s/wan-slas", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// GetWANSla returns a profile by ID. A missing profile is a *NotFoundError.
func (c *ApiClient) GetWANSla(ctx context.Context, site, id string) (*WANSla, error) {
	var respBody WANSla

	err := c.do(ctx, http.MethodGet, fmt.Sprintf("v2/api/site/%s/wan-slas/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}

func (c *ApiClient) CreateWANSla(ctx context.Context, site string, d *WANSla) (*WANSla, error) {
	var respBody WANSla

	err := c.do(ctx, http.MethodPost, fmt.Sprintf("v2/api/site/%s/wan-slas", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}

func (c *ApiClient) UpdateWANSla(ctx context.Context, site string, d *WANSla) (*WANSla, error) {
	var respBody WANSla

	err := c.do(ctx, http.MethodPut, fmt.Sprintf("v2/api/site/%s/wan-slas/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}

func (c *ApiClient) DeleteWANSla(ctx context.Context, site, id string) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("v2/api/site/%s/wan-slas/%s", site, id), struct{}{}, nil)
}
