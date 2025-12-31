// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type PortProfile struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Autoneg                      bool                  `json:"autoneg"`
	Dot1XCtrl                    string                `json:"dot1x_ctrl,omitempty"`             // auto|force_authorized|force_unauthorized|mac_based|multi_host
	Dot1XIDleTimeout             int                   `json:"dot1x_idle_timeout,omitempty"`     // [0-9]|[1-9][0-9]{1,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]
	EgressRateLimitKbps          int                   `json:"egress_rate_limit_kbps,omitempty"` // 6[4-9]|[7-9][0-9]|[1-9][0-9]{2,6}
	EgressRateLimitKbpsEnabled   bool                  `json:"egress_rate_limit_kbps_enabled"`
	ExcludedNetworkIDs           []string              `json:"excluded_networkconf_ids,omitempty"`
	FecMode                      string                `json:"fec_mode,omitempty"` // rs-fec|fc-fec|default|disabled
	Forward                      string                `json:"forward,omitempty"`  // all|native|customize|disabled
	FullDuplex                   bool                  `json:"full_duplex"`
	Isolation                    bool                  `json:"isolation"`
	LldpmedEnabled               bool                  `json:"lldpmed_enabled"`
	LldpmedNotifyEnabled         bool                  `json:"lldpmed_notify_enabled"`
	MulticastRouterNetworkIDs    []string              `json:"multicast_router_networkconf_ids,omitempty"`
	NATiveNetworkID              string                `json:"native_networkconf_id,omitempty"`
	Name                         string                `json:"name,omitempty"`
	OpMode                       string                `json:"op_mode,omitempty"`  // switch
	PoeMode                      string                `json:"poe_mode,omitempty"` // auto|off
	PortKeepaliveEnabled         bool                  `json:"port_keepalive_enabled"`
	PortSecurityEnabled          bool                  `json:"port_security_enabled"`
	PortSecurityMACAddress       []string              `json:"port_security_mac_address,omitempty"` // ^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$
	PriorityQueue1Level          int                   `json:"priority_queue1_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	PriorityQueue2Level          int                   `json:"priority_queue2_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	PriorityQueue3Level          int                   `json:"priority_queue3_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	PriorityQueue4Level          int                   `json:"priority_queue4_level,omitempty"`     // [0-9]|[1-9][0-9]|100
	QOSProfile                   PortProfileQOSProfile `json:"qos_profile,omitempty"`
	SettingPreference            string                `json:"setting_preference,omitempty"` // auto|manual
	Speed                        int                   `json:"speed,omitempty"`              // 10|100|1000|2500|5000|10000|20000|25000|40000|50000|100000
	StormctrlBroadcastastEnabled bool                  `json:"stormctrl_bcast_enabled"`
	StormctrlBroadcastastLevel   int                   `json:"stormctrl_bcast_level,omitempty"` // [0-9]|[1-9][0-9]|100
	StormctrlBroadcastastRate    int                   `json:"stormctrl_bcast_rate,omitempty"`  // [0-9]|[1-9][0-9]{1,6}|1[0-3][0-9]{6}|14[0-7][0-9]{5}|148[0-7][0-9]{4}|14880000
	StormctrlMcastEnabled        bool                  `json:"stormctrl_mcast_enabled"`
	StormctrlMcastLevel          int                   `json:"stormctrl_mcast_level,omitempty"` // [0-9]|[1-9][0-9]|100
	StormctrlMcastRate           int                   `json:"stormctrl_mcast_rate,omitempty"`  // [0-9]|[1-9][0-9]{1,6}|1[0-3][0-9]{6}|14[0-7][0-9]{5}|148[0-7][0-9]{4}|14880000
	StormctrlType                string                `json:"stormctrl_type,omitempty"`        // level|rate
	StormctrlUcastEnabled        bool                  `json:"stormctrl_ucast_enabled"`
	StormctrlUcastLevel          int                   `json:"stormctrl_ucast_level,omitempty"` // [0-9]|[1-9][0-9]|100
	StormctrlUcastRate           int                   `json:"stormctrl_ucast_rate,omitempty"`  // [0-9]|[1-9][0-9]{1,6}|1[0-3][0-9]{6}|14[0-7][0-9]{5}|148[0-7][0-9]{4}|14880000
	StpPortMode                  bool                  `json:"stp_port_mode"`
	TaggedVLANMgmt               string                `json:"tagged_vlan_mgmt,omitempty"` // auto|block_all|custom
	VoiceNetworkID               string                `json:"voice_networkconf_id,omitempty"`
}

func (dst *PortProfile) UnmarshalJSON(b []byte) error {
	type Alias PortProfile
	aux := &struct {
		Dot1XIDleTimeout           types.Number `json:"dot1x_idle_timeout"`
		EgressRateLimitKbps        types.Number `json:"egress_rate_limit_kbps"`
		PriorityQueue1Level        types.Number `json:"priority_queue1_level"`
		PriorityQueue2Level        types.Number `json:"priority_queue2_level"`
		PriorityQueue3Level        types.Number `json:"priority_queue3_level"`
		PriorityQueue4Level        types.Number `json:"priority_queue4_level"`
		Speed                      types.Number `json:"speed"`
		StormctrlBroadcastastLevel types.Number `json:"stormctrl_bcast_level"`
		StormctrlBroadcastastRate  types.Number `json:"stormctrl_bcast_rate"`
		StormctrlMcastLevel        types.Number `json:"stormctrl_mcast_level"`
		StormctrlMcastRate         types.Number `json:"stormctrl_mcast_rate"`
		StormctrlUcastLevel        types.Number `json:"stormctrl_ucast_level"`
		StormctrlUcastRate         types.Number `json:"stormctrl_ucast_rate"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if val, err := aux.Dot1XIDleTimeout.Int64(); err == nil {
		dst.Dot1XIDleTimeout = int(val)
	}
	if val, err := aux.EgressRateLimitKbps.Int64(); err == nil {
		dst.EgressRateLimitKbps = int(val)
	}
	if val, err := aux.PriorityQueue1Level.Int64(); err == nil {
		dst.PriorityQueue1Level = int(val)
	}
	if val, err := aux.PriorityQueue2Level.Int64(); err == nil {
		dst.PriorityQueue2Level = int(val)
	}
	if val, err := aux.PriorityQueue3Level.Int64(); err == nil {
		dst.PriorityQueue3Level = int(val)
	}
	if val, err := aux.PriorityQueue4Level.Int64(); err == nil {
		dst.PriorityQueue4Level = int(val)
	}
	if val, err := aux.Speed.Int64(); err == nil {
		dst.Speed = int(val)
	}
	if val, err := aux.StormctrlBroadcastastLevel.Int64(); err == nil {
		dst.StormctrlBroadcastastLevel = int(val)
	}
	if val, err := aux.StormctrlBroadcastastRate.Int64(); err == nil {
		dst.StormctrlBroadcastastRate = int(val)
	}
	if val, err := aux.StormctrlMcastLevel.Int64(); err == nil {
		dst.StormctrlMcastLevel = int(val)
	}
	if val, err := aux.StormctrlMcastRate.Int64(); err == nil {
		dst.StormctrlMcastRate = int(val)
	}
	if val, err := aux.StormctrlUcastLevel.Int64(); err == nil {
		dst.StormctrlUcastLevel = int(val)
	}
	if val, err := aux.StormctrlUcastRate.Int64(); err == nil {
		dst.StormctrlUcastRate = int(val)
	}

	return nil
}

type PortProfileQOSMarking struct {
	CosCode          int `json:"cos_code,omitempty"`           // [0-7]
	DscpCode         int `json:"dscp_code,omitempty"`          // 0|8|16|24|32|40|48|56|10|12|14|18|20|22|26|28|30|34|36|38|44|46
	IPPrecedenceCode int `json:"ip_precedence_code,omitempty"` // [0-7]
	Queue            int `json:"queue,omitempty"`              // [0-7]
}

func (dst *PortProfileQOSMarking) UnmarshalJSON(b []byte) error {
	type Alias PortProfileQOSMarking
	aux := &struct {
		CosCode          types.Number `json:"cos_code"`
		DscpCode         types.Number `json:"dscp_code"`
		IPPrecedenceCode types.Number `json:"ip_precedence_code"`
		Queue            types.Number `json:"queue"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if val, err := aux.CosCode.Int64(); err == nil {
		dst.CosCode = int(val)
	}
	if val, err := aux.DscpCode.Int64(); err == nil {
		dst.DscpCode = int(val)
	}
	if val, err := aux.IPPrecedenceCode.Int64(); err == nil {
		dst.IPPrecedenceCode = int(val)
	}
	if val, err := aux.Queue.Int64(); err == nil {
		dst.Queue = int(val)
	}

	return nil
}

type PortProfileQOSMatching struct {
	CosCode          int    `json:"cos_code,omitempty"`           // [0-7]
	DscpCode         int    `json:"dscp_code,omitempty"`          // [0-9]|[1-5][0-9]|6[0-3]
	DstPort          int    `json:"dst_port,omitempty"`           // [0-9]|[1-9][0-9]|[1-9][0-9][0-9]|[1-9][0-9][0-9][0-9]|[1-5][0-9][0-9][0-9][0-9]|6[0-4][0-9][0-9][0-9]|65[0-4][0-9][0-9]|655[0-2][0-9]|6553[0-4]|65535
	IPPrecedenceCode int    `json:"ip_precedence_code,omitempty"` // [0-7]
	Protocol         string `json:"protocol,omitempty"`           // ([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])|ah|ax.25|dccp|ddp|egp|eigrp|encap|esp|etherip|fc|ggp|gre|hip|hmp|icmp|idpr-cmtp|idrp|igmp|igp|ip|ipcomp|ipencap|ipip|ipv6|ipv6-frag|ipv6-icmp|ipv6-nonxt|ipv6-opts|ipv6-route|isis|iso-tp4|l2tp|manet|mobility-header|mpls-in-ip|ospf|pim|pup|rdp|rohc|rspf|rsvp|sctp|shim6|skip|st|tcp|udp|udplite|vmtp|vrrp|wesp|xns-idp|xtp
	SrcPort          int    `json:"src_port,omitempty"`           // [0-9]|[1-9][0-9]|[1-9][0-9][0-9]|[1-9][0-9][0-9][0-9]|[1-5][0-9][0-9][0-9][0-9]|6[0-4][0-9][0-9][0-9]|65[0-4][0-9][0-9]|655[0-2][0-9]|6553[0-4]|65535
}

func (dst *PortProfileQOSMatching) UnmarshalJSON(b []byte) error {
	type Alias PortProfileQOSMatching
	aux := &struct {
		CosCode          types.Number `json:"cos_code"`
		DscpCode         types.Number `json:"dscp_code"`
		DstPort          types.Number `json:"dst_port"`
		IPPrecedenceCode types.Number `json:"ip_precedence_code"`
		SrcPort          types.Number `json:"src_port"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if val, err := aux.CosCode.Int64(); err == nil {
		dst.CosCode = int(val)
	}
	if val, err := aux.DscpCode.Int64(); err == nil {
		dst.DscpCode = int(val)
	}
	if val, err := aux.DstPort.Int64(); err == nil {
		dst.DstPort = int(val)
	}
	if val, err := aux.IPPrecedenceCode.Int64(); err == nil {
		dst.IPPrecedenceCode = int(val)
	}
	if val, err := aux.SrcPort.Int64(); err == nil {
		dst.SrcPort = int(val)
	}

	return nil
}

type PortProfileQOSPolicies struct {
	QOSMarking  PortProfileQOSMarking  `json:"qos_marking,omitempty"`
	QOSMatching PortProfileQOSMatching `json:"qos_matching,omitempty"`
}

func (dst *PortProfileQOSPolicies) UnmarshalJSON(b []byte) error {
	type Alias PortProfileQOSPolicies
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

type PortProfileQOSProfile struct {
	QOSPolicies    []PortProfileQOSPolicies `json:"qos_policies,omitempty"`
	QOSProfileMode string                   `json:"qos_profile_mode,omitempty"` // custom|unifi_play|aes67_audio|crestron_audio_video|dante_audio|ndi_aes67_audio|ndi_dante_audio|qsys_audio_video|qsys_video_dante_audio|sdvoe_aes67_audio|sdvoe_dante_audio|shure_audio
}

func (dst *PortProfileQOSProfile) UnmarshalJSON(b []byte) error {
	type Alias PortProfileQOSProfile
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

func (c *Client) listPortProfile(ctx context.Context, site string) ([]PortProfile, error) {
	var respBody struct {
		Meta meta          `json:"meta"`
		Data []PortProfile `json:"data"`
	}

	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("api/s/%s/rest/portconf", site),
		nil,
		&respBody,
	)
	if err != nil {
		return nil, err
	}
	return respBody.Data, nil
}

func (c *Client) getPortProfile(
	ctx context.Context,
	site string,
	id string,
) (*PortProfile, error) {
	var respBody struct {
		Meta meta          `json:"meta"`
		Data []PortProfile `json:"data"`
	}
	err := c.do(
		ctx,
		"GET",
		fmt.Sprintf("api/s/%s/rest/portconf/%s", site, id),
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

func (c *Client) deletePortProfile(
	ctx context.Context,
	site string,
	id string,
) error {
	err := c.do(
		ctx,
		"DELETE",
		fmt.Sprintf("api/s/%s/rest/portconf/%s", site, id),
		struct{}{},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createPortProfile(
	ctx context.Context,
	site string,
	d *PortProfile,
) (*PortProfile, error) {
	var respBody struct {
		Meta meta          `json:"meta"`
		Data []PortProfile `json:"data"`
	}

	err := c.do(
		ctx,
		"POST",
		fmt.Sprintf("api/s/%s/rest/portconf", site),
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

func (c *Client) updatePortProfile(
	ctx context.Context,
	site string,
	d *PortProfile,
) (*PortProfile, error) {
	var respBody struct {
		Meta meta          `json:"meta"`
		Data []PortProfile `json:"data"`
	}
	err := c.do(
		ctx,
		"PUT",
		fmt.Sprintf("api/s/%s/rest/portconf/%s", site, d.ID),
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
