package unifi

import (
	"context"
	"fmt"
	"net/http"
)

// This is a v2 API object, manually coded.
//
// Power supervisors back the "Device Supervisor" feature (UniFi Network 10.2+):
// heartbeat monitoring of a device (keyed by client_mac) plus automatic
// power-cycling of its upstream PoE source after a silence threshold. The
// controller exposes them as a dedicated per-device collection on the v2 API
// (one entry per supervised device), not as a field on the device object.

type PowerSupervisor struct {
	ID     string `json:"id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	ClientMAC string                  `json:"client_mac"`
	Enabled   bool                    `json:"enabled"`
	Settings  PowerSupervisorSettings `json:"settings"`
	// PowerSources may be sent empty on create — the controller resolves the
	// upstream PoE port automatically and populates them on read.
	PowerSources []PowerSupervisorSource `json:"power_sources"`

	// Computed/read-only.
	ConsecutiveFailures int `json:"consecutive_failures,omitempty"`
}

// PowerSupervisorSettings are the heartbeat/power-cycle timings, all in seconds.
type PowerSupervisorSettings struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
	SilenceThreshold  int `json:"silence_threshold"`
	PowerOffDuration  int `json:"power_off_duration"`
}

// PowerSupervisorSource references the upstream power source used to recover the
// supervised device. The controller resolves and fills these on read.
type PowerSupervisorSource struct {
	ClientPsuIndex   int    `json:"client_psu_index,omitempty"`
	PowerSourceIndex int    `json:"power_source_index,omitempty"`
	PowerSourceMAC   string `json:"power_source_mac,omitempty"`
	PowerSourceType  string `json:"power_source_type,omitempty"`
}

func (c *ApiClient) powerSupervisorsPath(site string) string {
	return fmt.Sprintf("v2/api/site/%s/power-supervisors", site)
}

func (c *ApiClient) ListPowerSupervisors(
	ctx context.Context,
	site string,
) ([]PowerSupervisor, error) {
	var respBody []PowerSupervisor

	err := c.do(ctx, http.MethodGet, c.powerSupervisorsPath(site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// GetPowerSupervisor fetches a single supervisor by ID by listing and filtering.
func (c *ApiClient) GetPowerSupervisor(
	ctx context.Context,
	site, id string,
) (*PowerSupervisor, error) {
	supervisors, err := c.ListPowerSupervisors(ctx, site)
	if err != nil {
		return nil, err
	}

	for i := range supervisors {
		if supervisors[i].ID == id {
			return &supervisors[i], nil
		}
	}

	return nil, &NotFoundError{Type: "power_supervisor", Attr: "id", Value: id}
}

// GetPowerSupervisorByMAC fetches the supervisor for a supervised device by its
// client MAC (the collection is keyed per device), useful for import.
func (c *ApiClient) GetPowerSupervisorByMAC(
	ctx context.Context,
	site, mac string,
) (*PowerSupervisor, error) {
	supervisors, err := c.ListPowerSupervisors(ctx, site)
	if err != nil {
		return nil, err
	}

	for i := range supervisors {
		if supervisors[i].ClientMAC == mac {
			return &supervisors[i], nil
		}
	}

	return nil, &NotFoundError{Type: "power_supervisor", Attr: "client_mac", Value: mac}
}

func (c *ApiClient) CreatePowerSupervisor(
	ctx context.Context,
	site string,
	d *PowerSupervisor,
) (*PowerSupervisor, error) {
	var respBody PowerSupervisor
	d.ID = ""
	if d.PowerSources == nil {
		d.PowerSources = []PowerSupervisorSource{}
	}

	err := c.do(ctx, http.MethodPost, c.powerSupervisorsPath(site), d, &respBody)
	if err != nil {
		return nil, err
	}

	// The create response may not echo the persisted object (and thus its id).
	// Fall back to listing and matching on the device MAC so callers always get
	// the server-assigned id back.
	if respBody.ID == "" {
		return c.GetPowerSupervisorByMAC(ctx, site, d.ClientMAC)
	}

	return &respBody, nil
}

func (c *ApiClient) UpdatePowerSupervisor(
	ctx context.Context,
	site string,
	d *PowerSupervisor,
) (*PowerSupervisor, error) {
	var respBody PowerSupervisor
	if d.PowerSources == nil {
		d.PowerSources = []PowerSupervisorSource{}
	}

	err := c.do(
		ctx,
		http.MethodPut,
		fmt.Sprintf("%s/%s", c.powerSupervisorsPath(site), d.ID),
		d,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	if respBody.ID == "" {
		return c.GetPowerSupervisor(ctx, site, d.ID)
	}

	return &respBody, nil
}

func (c *ApiClient) DeletePowerSupervisor(
	ctx context.Context,
	site, id string,
) error {
	return c.do(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s", c.powerSupervisorsPath(site), id),
		struct{}{},
		nil,
	)
}
