package unifi

import "encoding/json"

// natFilterJSON is the wire shape of a NAT filter. Unlike the generated structs it
// always sends firewall_group_ids as an array; an empty address must be omitted,
// as the controller rejects "" as an invalid IP address.
type natFilterJSON struct {
	Address          string   `json:"address,omitempty"`
	FilterType       string   `json:"filter_type,omitempty"`
	FirewallGroupIDs []string `json:"firewall_group_ids"`
	InvertAddress    bool     `json:"invert_address"`
	InvertPort       bool     `json:"invert_port"`
	NetworkConfID    string   `json:"network_conf_id,omitempty"`
	Port             *int64   `json:"port,omitempty"`
}

func marshalNatFilter(
	address, filterType string,
	groups []string,
	invertAddress, invertPort bool,
	networkConfID string,
	port *int64,
) ([]byte, error) {
	if groups == nil {
		groups = []string{}
	}
	return json.Marshal(natFilterJSON{
		Address:          address,
		FilterType:       filterType,
		FirewallGroupIDs: groups,
		InvertAddress:    invertAddress,
		InvertPort:       invertPort,
		NetworkConfID:    networkConfID,
		Port:             port,
	})
}

func (f NatSourceFilter) MarshalJSON() ([]byte, error) {
	return marshalNatFilter(
		f.Address, f.FilterType, f.FirewallGroupIDs,
		f.InvertAddress, f.InvertPort, f.NetworkConfID, f.Port,
	)
}

func (f NatDestinationFilter) MarshalJSON() ([]byte, error) {
	return marshalNatFilter(
		f.Address, f.FilterType, f.FirewallGroupIDs,
		f.InvertAddress, f.InvertPort, f.NetworkConfID, f.Port,
	)
}
