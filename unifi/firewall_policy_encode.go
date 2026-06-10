package unifi

import (
	"encoding/json"
	"strconv"
)

// The zone-based firewall API (UniFi Network 8.x+) rejects a PUT/POST whose
// source/destination `port` is encoded as a JSON number. The firmware expects
// the port as a quoted string (e.g. "161"); decoding already tolerates either
// form via the generated UnmarshalJSON. These custom marshalers re-encode the
// port as a string while leaving every other field to the generated layout.
//
// The anonymous wrapper places a string `port` at depth 0, which dominates the
// embedded alias' numeric `port` at depth 1 per encoding/json's field rules, so
// the numeric field is suppressed. An empty string (nil port) is dropped by
// omitempty, matching the firmware's expectation when port_matching_type is ANY.

func (s FirewallPolicySource) MarshalJSON() ([]byte, error) {
	type Alias FirewallPolicySource
	return json.Marshal(&struct {
		Port string `json:"port,omitempty"`
		Alias
	}{
		Port:  portToString(s.Port),
		Alias: Alias(s),
	})
}

func (d FirewallPolicyDestination) MarshalJSON() ([]byte, error) {
	type Alias FirewallPolicyDestination
	return json.Marshal(&struct {
		Port string `json:"port,omitempty"`
		Alias
	}{
		Port:  portToString(d.Port),
		Alias: Alias(d),
	})
}

func portToString(p *int64) string {
	if p == nil {
		return ""
	}
	return strconv.FormatInt(*p, 10)
}
