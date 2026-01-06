// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package settings

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
	_ types.Number
)

type Usg struct {
	BaseSetting

	ArpCacheBaseReachable          int                       `json:"arp_cache_base_reachable,omitempty"` // ^$|^[1-9]{1}[0-9]{0,4}$
	ArpCacheTimeout                string                    `json:"arp_cache_timeout,omitempty"`        // normal|min-dhcp-lease|custom
	BroadcastPing                  bool                      `json:"broadcast_ping"`
	DHCPDHostfileUpdate            bool                      `json:"dhcpd_hostfile_update"`
	DHCPDUseDNSmasq                bool                      `json:"dhcpd_use_dnsmasq"`
	DHCPRelayAgentsPackets         string                    `json:"dhcp_relay_agents_packets"`      // append|discard|forward|replace|^$
	DHCPRelayHopCount              int                       `json:"dhcp_relay_hop_count,omitempty"` // ([1-9]|[1-8][0-9]|9[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])|^$
	DHCPRelayMaxSize               int                       `json:"dhcp_relay_max_size,omitempty"`  // (6[4-9]|[7-9][0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|1[0-3][0-9]{2}|1400)|^$
	DHCPRelayPort                  int                       `json:"dhcp_relay_port,omitempty"`      // [1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5]|^$
	DHCPRelayServer1               string                    `json:"dhcp_relay_server_1"`            // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPRelayServer2               string                    `json:"dhcp_relay_server_2"`            // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPRelayServer3               string                    `json:"dhcp_relay_server_3"`            // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPRelayServer4               string                    `json:"dhcp_relay_server_4"`            // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DHCPRelayServer5               string                    `json:"dhcp_relay_server_5"`            // ^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$|^$
	DNSVerification                SettingUsgDNSVerification `json:"dns_verification,omitempty"`
	DNSmasqAllServers              bool                      `json:"dnsmasq_all_servers"`
	EchoServer                     string                    `json:"echo_server,omitempty"` // [^\"\' ]{1,255}
	FtpModule                      bool                      `json:"ftp_module"`
	GeoIPFilteringBlock            string                    `json:"geo_ip_filtering_block,omitempty"`     // block|allow
	GeoIPFilteringCountries        string                    `json:"geo_ip_filtering_countries,omitempty"` // ^([A-Z]{2})?(,[A-Z]{2}){0,149}$
	GeoIPFilteringEnabled          bool                      `json:"geo_ip_filtering_enabled"`
	GeoIPFilteringTrafficDirection string                    `json:"geo_ip_filtering_traffic_direction,omitempty"` // ^(both|ingress|egress)$
	GreModule                      bool                      `json:"gre_module"`
	H323Module                     bool                      `json:"h323_module"`
	ICMPTimeout                    int                       `json:"icmp_timeout,omitempty"`
	LldpEnableAll                  bool                      `json:"lldp_enable_all"`
	MdnsEnabled                    bool                      `json:"mdns_enabled"`
	MssClamp                       string                    `json:"mss_clamp,omitempty"`     // auto|custom|disabled
	MssClampMss                    int                       `json:"mss_clamp_mss,omitempty"` // [1-9][0-9]{2,3}
	OffloadAccounting              bool                      `json:"offload_accounting"`
	OffloadL2Blocking              bool                      `json:"offload_l2_blocking"`
	OffloadSch                     bool                      `json:"offload_sch"`
	OtherTimeout                   int                       `json:"other_timeout,omitempty"`
	PptpModule                     bool                      `json:"pptp_module"`
	ReceiveRedirects               bool                      `json:"receive_redirects"`
	SendRedirects                  bool                      `json:"send_redirects"`
	SipModule                      bool                      `json:"sip_module"`
	SynCookies                     bool                      `json:"syn_cookies"`
	TCPCloseTimeout                int                       `json:"tcp_close_timeout,omitempty"`
	TCPCloseWaitTimeout            int                       `json:"tcp_close_wait_timeout,omitempty"`
	TCPEstablishedTimeout          int                       `json:"tcp_established_timeout,omitempty"`
	TCPFinWaitTimeout              int                       `json:"tcp_fin_wait_timeout,omitempty"`
	TCPLastAckTimeout              int                       `json:"tcp_last_ack_timeout,omitempty"`
	TCPSynRecvTimeout              int                       `json:"tcp_syn_recv_timeout,omitempty"`
	TCPSynSentTimeout              int                       `json:"tcp_syn_sent_timeout,omitempty"`
	TCPTimeWaitTimeout             int                       `json:"tcp_time_wait_timeout,omitempty"`
	TFTPModule                     bool                      `json:"tftp_module"`
	TimeoutSettingPreference       string                    `json:"timeout_setting_preference,omitempty"` // auto|manual
	UDPOtherTimeout                int                       `json:"udp_other_timeout,omitempty"`
	UDPStreamTimeout               int                       `json:"udp_stream_timeout,omitempty"`
	UPnPEnabled                    bool                      `json:"upnp_enabled"`
	UPnPNATPmpEnabled              bool                      `json:"upnp_nat_pmp_enabled"`
	UPnPSecureMode                 bool                      `json:"upnp_secure_mode"`
	UPnPWANInterface               string                    `json:"upnp_wan_interface,omitempty"` // WAN[2-9]?
	UnbindWANMonitors              bool                      `json:"unbind_wan_monitors"`
}

func (dst *Usg) UnmarshalJSON(b []byte) error {
	type Alias Usg
	aux := &struct {
		ArpCacheBaseReachable types.Number `json:"arp_cache_base_reachable"`
		DHCPRelayHopCount     types.Number `json:"dhcp_relay_hop_count"`
		DHCPRelayMaxSize      types.Number `json:"dhcp_relay_max_size"`
		DHCPRelayPort         types.Number `json:"dhcp_relay_port"`
		ICMPTimeout           types.Number `json:"icmp_timeout"`
		MssClampMss           types.Number `json:"mss_clamp_mss"`
		OtherTimeout          types.Number `json:"other_timeout"`
		TCPCloseTimeout       types.Number `json:"tcp_close_timeout"`
		TCPCloseWaitTimeout   types.Number `json:"tcp_close_wait_timeout"`
		TCPEstablishedTimeout types.Number `json:"tcp_established_timeout"`
		TCPFinWaitTimeout     types.Number `json:"tcp_fin_wait_timeout"`
		TCPLastAckTimeout     types.Number `json:"tcp_last_ack_timeout"`
		TCPSynRecvTimeout     types.Number `json:"tcp_syn_recv_timeout"`
		TCPSynSentTimeout     types.Number `json:"tcp_syn_sent_timeout"`
		TCPTimeWaitTimeout    types.Number `json:"tcp_time_wait_timeout"`
		UDPOtherTimeout       types.Number `json:"udp_other_timeout"`
		UDPStreamTimeout      types.Number `json:"udp_stream_timeout"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	// First unmarshal base setting
	if err := json.Unmarshal(b, &dst.BaseSetting); err != nil {
		return fmt.Errorf("unable to unmarshal base setting: %w", err)
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if val, err := aux.ArpCacheBaseReachable.Int64(); err == nil {
		dst.ArpCacheBaseReachable = int(val)
	}
	if val, err := aux.DHCPRelayHopCount.Int64(); err == nil {
		dst.DHCPRelayHopCount = int(val)
	}
	if val, err := aux.DHCPRelayMaxSize.Int64(); err == nil {
		dst.DHCPRelayMaxSize = int(val)
	}
	if val, err := aux.DHCPRelayPort.Int64(); err == nil {
		dst.DHCPRelayPort = int(val)
	}
	if val, err := aux.ICMPTimeout.Int64(); err == nil {
		dst.ICMPTimeout = int(val)
	}
	if val, err := aux.MssClampMss.Int64(); err == nil {
		dst.MssClampMss = int(val)
	}
	if val, err := aux.OtherTimeout.Int64(); err == nil {
		dst.OtherTimeout = int(val)
	}
	if val, err := aux.TCPCloseTimeout.Int64(); err == nil {
		dst.TCPCloseTimeout = int(val)
	}
	if val, err := aux.TCPCloseWaitTimeout.Int64(); err == nil {
		dst.TCPCloseWaitTimeout = int(val)
	}
	if val, err := aux.TCPEstablishedTimeout.Int64(); err == nil {
		dst.TCPEstablishedTimeout = int(val)
	}
	if val, err := aux.TCPFinWaitTimeout.Int64(); err == nil {
		dst.TCPFinWaitTimeout = int(val)
	}
	if val, err := aux.TCPLastAckTimeout.Int64(); err == nil {
		dst.TCPLastAckTimeout = int(val)
	}
	if val, err := aux.TCPSynRecvTimeout.Int64(); err == nil {
		dst.TCPSynRecvTimeout = int(val)
	}
	if val, err := aux.TCPSynSentTimeout.Int64(); err == nil {
		dst.TCPSynSentTimeout = int(val)
	}
	if val, err := aux.TCPTimeWaitTimeout.Int64(); err == nil {
		dst.TCPTimeWaitTimeout = int(val)
	}
	if val, err := aux.UDPOtherTimeout.Int64(); err == nil {
		dst.UDPOtherTimeout = int(val)
	}
	if val, err := aux.UDPStreamTimeout.Int64(); err == nil {
		dst.UDPStreamTimeout = int(val)
	}

	return nil
}

type SettingUsgDNSVerification struct {
	Domain             string `json:"domain,omitempty"`
	PrimaryDNSServer   string `json:"primary_dns_server,omitempty"`
	SecondaryDNSServer string `json:"secondary_dns_server,omitempty"`
	SettingPreference  string `json:"setting_preference,omitempty"` // auto|manual
}

func (dst *SettingUsgDNSVerification) UnmarshalJSON(b []byte) error {
	type Alias SettingUsgDNSVerification
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
