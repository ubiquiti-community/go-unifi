package unifi

import "encoding/json"

// MarshalJSON customizes Account JSON serialization.
// Fields VLAN, TunnelType, and TunnelMediumType are serialized as empty strings
// when nil, to match the UniFi API expectation that these fields are always present.
func (a Account) MarshalJSON() ([]byte, error) {
	type Alias Account

	marshalOptionalInt64 := func(v *int64) any {
		if v == nil {
			return ""
		}
		return *v
	}

	return json.Marshal(struct {
		VLAN             any `json:"vlan"`
		TunnelType       any `json:"tunnel_type"`
		TunnelMediumType any `json:"tunnel_medium_type"`
		*Alias
	}{
		VLAN:             marshalOptionalInt64(a.VLAN),
		TunnelType:       marshalOptionalInt64(a.TunnelType),
		TunnelMediumType: marshalOptionalInt64(a.TunnelMediumType),
		Alias:            (*Alias)(&a),
	})
}
