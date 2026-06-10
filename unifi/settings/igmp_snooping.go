package settings

// IgmpSnooping is the site-level IGMP snooping setting (key "igmp_snooping").
//
// Hand-maintained: the ace.jar field spec does not define this setting, but the
// controller exposes it at /api/s/<site>/{get,set}/setting/igmp_snooping. On
// UniFi Network 10.3.x the effective IGMP snooping toggle moved here from the
// per-network object (see ubiquiti-community/terraform-provider-unifi#164).
//
// All fields are simple scalars/slices, so the default JSON (un)marshalling of
// the embedded BaseSetting plus these fields is sufficient — no custom
// UnmarshalJSON is needed.
type IgmpSnooping struct {
	BaseSetting

	Enabled                            bool     `json:"enabled"`
	NetworkIDs                         []string `json:"network_ids,omitempty"`
	FloodKnownProtocols                bool     `json:"flood_known_protocols"`
	ForwardUnknownMcastRouterPorts     bool     `json:"forward_unknown_mcast_router_ports"`
	FastleaveForNetworkIDs             []string `json:"fastleave_for_network_ids,omitempty"`
	FloodUnknownMulticastForNetworkIDs []string `json:"flood_unknown_multicast_for_network_ids,omitempty"`
	SubscriptionMode                   string   `json:"subscription_mode,omitempty"`
	QuerierMode                        string   `json:"querier_mode,omitempty"`
	QuerierSubscriptionMode            string   `json:"querier_subscription_mode,omitempty"`
	QuerierSwitches                    []string `json:"querier_switches,omitempty"`
	QuerierAddresses                   []string `json:"querier_addresses,omitempty"`
	Switches                           []string `json:"switches,omitempty"`
	PrimaryQuerier                     string   `json:"primary_querier,omitempty"`
	FailoverQuerier                    string   `json:"failover_querier,omitempty"`
}
