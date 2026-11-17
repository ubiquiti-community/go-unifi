package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
)

const (
	globalNetworkConfigPath = "v2/api/site/%s/global/config/network"
	keyMdnsEnabledFor       = "mdns_enabled_for"
	keyMdnsEnabledForIDs    = "mdns_enabled_for_network_ids"
)

// GlobalNetworkConfig is the site-wide network settings object behind
// Settings > Networks in the UniFi Network UI, served by the v2 endpoint
// global/config/network. It is a v2 API object, so it is coded by hand.
//
// On Network 11 this object holds the mDNS reflector's network selection. The
// legacy mdns setting's enabled_for and enabled_for_network_ids, and each
// network's mdns_enabled flag, only mirror it there: the controller answers
// writes to them with rc "ok" and ignores them.
//
// Only the mDNS fields are modelled. Every other key the controller returns is
// kept and sent back unchanged by UpdateGlobalNetworkConfig, since the UI
// replaces the whole object on save. Update a value returned by
// GetGlobalNetworkConfig rather than one built from scratch.
type GlobalNetworkConfig struct {
	// MdnsEnabledFor is "all", "some" or "none".
	MdnsEnabledFor string `json:"mdns_enabled_for,omitempty"`
	// MdnsEnabledForNetworkIDs lists the networks mDNS covers when
	// MdnsEnabledFor is "some".
	MdnsEnabledForNetworkIDs []string `json:"mdns_enabled_for_network_ids"`

	// raw holds every key as the controller returned it.
	raw map[string]json.RawMessage
}

// HasMdns reports whether the controller returned the mDNS selection in this
// object. Controllers that keep it only in the legacy mdns setting do not.
func (d *GlobalNetworkConfig) HasMdns() bool {
	_, ok := d.raw[keyMdnsEnabledFor]
	return ok
}

// UnmarshalJSON decodes the modelled fields and keeps every key for
// MarshalJSON.
func (d *GlobalNetworkConfig) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type alias GlobalNetworkConfig
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*d = GlobalNetworkConfig(a)
	d.raw = raw
	return nil
}

// MarshalJSON writes back every key the controller returned, with the
// modelled fields overlaid. A nil MdnsEnabledForNetworkIDs is sent as [].
func (d GlobalNetworkConfig) MarshalJSON() ([]byte, error) {
	out := make(map[string]json.RawMessage, len(d.raw)+2)
	maps.Copy(out, d.raw)
	if d.MdnsEnabledFor != "" {
		v, err := json.Marshal(d.MdnsEnabledFor)
		if err != nil {
			return nil, err
		}
		out[keyMdnsEnabledFor] = v
	}
	ids := d.MdnsEnabledForNetworkIDs
	if ids == nil {
		ids = []string{}
	}
	v, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}
	out[keyMdnsEnabledForIDs] = v
	return json.Marshal(out)
}

// GetGlobalNetworkConfig reads the site-wide network settings. Controllers
// without the v2 endpoint return a *NotFoundError.
func (c *ApiClient) GetGlobalNetworkConfig(
	ctx context.Context,
	site string,
) (*GlobalNetworkConfig, error) {
	var respBody GlobalNetworkConfig
	err := c.do(ctx, http.MethodGet, fmt.Sprintf(globalNetworkConfigPath, site), nil, &respBody)
	if err != nil {
		return nil, err
	}
	return &respBody, nil
}

// UpdateGlobalNetworkConfig replaces the site-wide network settings with d and
// returns what the controller stored.
func (c *ApiClient) UpdateGlobalNetworkConfig(
	ctx context.Context,
	site string,
	d *GlobalNetworkConfig,
) (*GlobalNetworkConfig, error) {
	var respBody GlobalNetworkConfig
	err := c.do(ctx, http.MethodPut, fmt.Sprintf(globalNetworkConfigPath, site), d, &respBody)
	if err != nil {
		return nil, err
	}
	return &respBody, nil
}
