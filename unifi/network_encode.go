package unifi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
)

const (
	PurposeCorporate = "corporate"
	PurposeGuest     = "guest"
	PurposeVLANOnly  = "vlan-only"
	PurposeWAN       = "wan"
	PurposeSiteVPN   = "site-vpn"
	PurposeVPNClient = "vpn-client"
	PurposeUserVPN   = "remote-user-vpn"
)

// MarshalJSON implements custom JSON marshaling that only includes fields relevant to the network's Purpose.
func (n *Network) MarshalJSON() ([]byte, error) {
	switch n.Purpose {
	case PurposeWAN:
		return n.marshalWAN()
	case PurposeCorporate:
		return n.marshalCorporate()
	case PurposeGuest:
		return n.marshalGuest()
	case PurposeVLANOnly:
		return n.marshalVLANOnly()
	case PurposeSiteVPN:
		return n.marshalSiteVPN()
	case PurposeVPNClient:
		return n.marshalVPNClient()
	case PurposeUserVPN:
		return n.marshalUserVPN()
	default:
		return nil, fmt.Errorf("unknown network purpose: %s", n.Purpose)
	}
}

// marshalCorporate marshals a Corporate/LAN network using the alias pattern.
func (n *Network) marshalCorporate() ([]byte, error) {
	// Calculate DHCP range defaults if needed
	var defaultStart, defaultEnd string
	if n.IPSubnet != nil {
		var err error
		defaultStart, defaultEnd, err = dhcpRange(*n.IPSubnet)
		if err != nil {
			log.Default().Printf("error calculating DHCP range: %s", err)
		}
	}

	// Use anonymous struct with explicit field selection
	return json.Marshal(&struct {
		ID       string `json:"_id,omitempty"`
		SiteID   string `json:"site_id,omitempty"`
		Hidden   bool   `json:"attr_hidden,omitempty"`
		HiddenID string `json:"attr_hidden_id,omitempty"`
		NoDelete bool   `json:"attr_no_delete,omitempty"`
		NoEdit   bool   `json:"attr_no_edit,omitempty"`

		Name                    *string                         `json:"name,omitempty"`
		Purpose                 string                          `json:"purpose"`
		Enabled                 bool                            `json:"enabled"`
		NetworkGroup            *string                         `json:"networkgroup,omitempty"`
		IPSubnet                *string                         `json:"ip_subnet,omitempty"`
		VLAN                    *int64                          `json:"vlan,omitempty"`
		VLANEnabled             bool                            `json:"vlan_enabled"`
		DomainName              *string                         `json:"domain_name,omitempty"`
		AutoScaleEnabled        bool                            `json:"auto_scale_enabled"`
		GatewayType             *string                         `json:"gateway_type,omitempty"`
		InternetAccessEnabled   bool                            `json:"internet_access_enabled"`
		NetworkIsolationEnabled bool                            `json:"network_isolation_enabled"`
		SettingPreference       *string                         `json:"setting_preference,omitempty"`
		IGMPSnooping            bool                            `json:"igmp_snooping"`
		DHCPguardEnabled        bool                            `json:"dhcpguard_enabled"`
		MdnsEnabled             bool                            `json:"mdns_enabled"`
		LteLanEnabled           bool                            `json:"lte_lan_enabled"`
		IPAliases               []string                        `json:"ip_aliases"`
		NATOutboundIPAddresses  []NetworkNATOutboundIPAddresses `json:"nat_outbound_ip_addresses"`
		MACOverride             string                          `json:"mac_override,omitempty"`

		// DHCP Server
		DHCPDEnabled           bool    `json:"dhcpd_enabled"`
		DHCPDStart             *string `json:"dhcpd_start,omitempty"`
		DHCPDStop              *string `json:"dhcpd_stop,omitempty"`
		DHCPDLeaseTime         *int64  `json:"dhcpd_leasetime,omitempty"`
		DHCPDDNSEnabled        bool    `json:"dhcpd_dns_enabled"`
		DHCPDDNS1              string  `json:"dhcpd_dns_1"`
		DHCPDDNS2              string  `json:"dhcpd_dns_2"`
		DHCPDDNS3              string  `json:"dhcpd_dns_3"`
		DHCPDDNS4              string  `json:"dhcpd_dns_4"`
		DHCPDGatewayEnabled    bool    `json:"dhcpd_gateway_enabled"`
		DHCPDGateway           *string `json:"dhcpd_gateway,omitempty"`
		DHCPDNtpEnabled        bool    `json:"dhcpd_ntp_enabled"`
		DHCPDNtp1              *string `json:"dhcpd_ntp_1,omitempty"`
		DHCPDNtp2              *string `json:"dhcpd_ntp_2,omitempty"`
		DHCPDWinsEnabled       bool    `json:"dhcpd_wins_enabled"`
		DHCPDWins1             *string `json:"dhcpd_wins_1,omitempty"`
		DHCPDWins2             *string `json:"dhcpd_wins_2,omitempty"`
		DHCPDTimeOffsetEnabled bool    `json:"dhcpd_time_offset_enabled"`
		DHCPDConflictChecking  bool    `json:"dhcpd_conflict_checking"`
		DHCPDBootEnabled       bool    `json:"dhcpd_boot_enabled"`
		DHCPDBootServer        string  `json:"dhcpd_boot_server,omitempty"`
		DHCPDBootFilename      string  `json:"dhcpd_boot_filename,omitempty"`
		DHCPDTFTPServer        *string `json:"dhcpd_tftp_server,omitempty"`
		DHCPDWPAdUrl           *string `json:"dhcpd_wpad_url,omitempty"`
		DHCPDUnifiController   *string `json:"dhcpd_unifi_controller,omitempty"`

		// DHCP Relay
		DHCPRelayEnabled bool     `json:"dhcp_relay_enabled"`
		DHCPRelayServers []string `json:"dhcp_relay_servers"`

		// IPv6
		IPV6InterfaceType     *string `json:"ipv6_interface_type,omitempty"`
		IPV6SettingPreference *string `json:"ipv6_setting_preference,omitempty"`
		IPV6RaPriority        *string `json:"ipv6_ra_priority,omitempty"`

		// DHCPv6
		DHCPDV6DNSAuto    bool    `json:"dhcpdv6_dns_auto,omitempty"`
		DHCPDV6AllowSlaac bool    `json:"dhcpdv6_allow_slaac,omitempty"`
		DHCPDV6Start      *string `json:"dhcpdv6_start,omitempty"`
		DHCPDV6Stop       *string `json:"dhcpdv6_stop,omitempty"`
		DHCPDV6LeaseTime  *int64  `json:"dhcpdv6_leasetime,omitempty"`
	}{
		ID:       n.ID,
		SiteID:   n.SiteID,
		Hidden:   n.Hidden,
		HiddenID: n.HiddenID,
		NoDelete: n.NoDelete,
		NoEdit:   n.NoEdit,

		Name:                    n.Name,
		Purpose:                 n.Purpose,
		Enabled:                 n.Enabled,
		NetworkGroup:            valueOrDefault(n.NetworkGroup, "LAN"),
		IPSubnet:                valueOrDefault(n.IPSubnet, ""),
		VLAN:                    n.VLAN,
		VLANEnabled:             n.VLANEnabled,
		DomainName:              valueOrDefault(n.DomainName, ""),
		AutoScaleEnabled:        n.AutoScaleEnabled,
		GatewayType:             valueOrDefault(n.GatewayType, "default"),
		InternetAccessEnabled:   n.InternetAccessEnabled,
		NetworkIsolationEnabled: n.NetworkIsolationEnabled,
		SettingPreference:       valueOrDefault(n.SettingPreference, "auto"),
		IGMPSnooping:            n.IGMPSnooping,
		DHCPguardEnabled:        n.DHCPguardEnabled,
		MdnsEnabled:             n.MdnsEnabled,
		LteLanEnabled:           n.LteLanEnabled,
		IPAliases:               orEmptySlice(n.IPAliases),
		NATOutboundIPAddresses:  orEmptyNATSlice(n.NATOutboundIPAddresses),
		MACOverride:             n.MACOverride,

		// DHCP Server with defaults
		DHCPDEnabled:           n.DHCPDEnabled,
		DHCPDStart:             valueOrDefault(n.DHCPDStart, defaultStart),
		DHCPDStop:              valueOrDefault(n.DHCPDStop, defaultEnd),
		DHCPDLeaseTime:         valueOrDefault(n.DHCPDLeaseTime, 86400),
		DHCPDDNSEnabled:        n.DHCPDDNSEnabled,
		DHCPDDNS1:              n.DHCPDDNS1,
		DHCPDDNS2:              n.DHCPDDNS2,
		DHCPDDNS3:              n.DHCPDDNS3,
		DHCPDDNS4:              n.DHCPDDNS4,
		DHCPDGatewayEnabled:    n.DHCPDGatewayEnabled,
		DHCPDGateway:           n.DHCPDGateway,
		DHCPDNtpEnabled:        n.DHCPDNtpEnabled,
		DHCPDNtp1:              nilIfEmpty(n.DHCPDNtp1),
		DHCPDNtp2:              nilIfEmpty(n.DHCPDNtp2),
		DHCPDWinsEnabled:       n.DHCPDWinsEnabled,
		DHCPDWins1:             valueOrDefault(n.DHCPDWins1, ""),
		DHCPDWins2:             valueOrDefault(n.DHCPDWins2, ""),
		DHCPDTimeOffsetEnabled: n.DHCPDTimeOffsetEnabled,
		DHCPDConflictChecking:  n.DHCPDConflictChecking,
		DHCPDBootEnabled:       n.DHCPDBootEnabled,
		DHCPDBootServer:        n.DHCPDBootServer,
		DHCPDBootFilename:      derefOrEmpty(n.DHCPDBootFilename),
		DHCPDTFTPServer:        n.DHCPDTFTPServer,
		DHCPDWPAdUrl:           n.DHCPDWPAdUrl,
		DHCPDUnifiController:   valueOrDefault(n.DHCPDUnifiController, ""),

		// DHCP Relay
		DHCPRelayEnabled: n.DHCPRelayEnabled,
		DHCPRelayServers: orEmptySlice(n.RemoteVPNSubnets),

		// IPv6
		IPV6InterfaceType:     valueOrDefault(n.IPV6InterfaceType, "none"),
		IPV6SettingPreference: n.IPV6SettingPreference,
		IPV6RaPriority:        n.IPV6RaPriority,

		// DHCPv6
		DHCPDV6DNSAuto:    n.DHCPDV6DNSAuto,
		DHCPDV6AllowSlaac: n.DHCPDV6AllowSlaac,
		DHCPDV6Start:      n.DHCPDV6Start,
		DHCPDV6Stop:       n.DHCPDV6Stop,
		DHCPDV6LeaseTime:  n.DHCPDV6LeaseTime,
	})
}

// marshalVLANOnly marshals a VLAN-only network (Layer 2 only, no routing).
func (n *Network) marshalVLANOnly() ([]byte, error) {
	enabled := true

	vlanEnabled := n.VLANEnabled
	if !vlanEnabled && n.VLAN != nil && *n.VLAN > 0 {
		vlanEnabled = true
	}

	return json.Marshal(&struct {
		ID       string `json:"_id,omitempty"`
		SiteID   string `json:"site_id,omitempty"`
		Hidden   bool   `json:"attr_hidden,omitempty"`
		HiddenID string `json:"attr_hidden_id,omitempty"`
		NoDelete bool   `json:"attr_no_delete,omitempty"`
		NoEdit   bool   `json:"attr_no_edit,omitempty"`

		Name                    *string `json:"name,omitempty"`
		Purpose                 string  `json:"purpose"`
		Enabled                 bool    `json:"enabled"`
		NetworkGroup            *string `json:"networkgroup,omitempty"`
		VLAN                    *int64  `json:"vlan,omitempty"`
		VLANEnabled             bool    `json:"vlan_enabled"`
		IGMPSnooping            bool    `json:"igmp_snooping"`
		NetworkIsolationEnabled bool    `json:"network_isolation_enabled"`
		DHCPguardEnabled        bool    `json:"dhcpguard_enabled"`
		DHCPDIP1                string  `json:"dhcpd_ip_1"`
		DHCPDIP2                string  `json:"dhcpd_ip_2"`
		DHCPDIP3                string  `json:"dhcpd_ip_3"`
	}{
		ID:       n.ID,
		SiteID:   n.SiteID,
		Hidden:   n.Hidden,
		HiddenID: n.HiddenID,
		NoDelete: n.NoDelete,
		NoEdit:   n.NoEdit,

		Name:                    n.Name,
		Purpose:                 n.Purpose,
		Enabled:                 enabled,
		NetworkGroup:            valueOrDefault(n.NetworkGroup, "LAN"),
		VLAN:                    n.VLAN,
		VLANEnabled:             vlanEnabled,
		IGMPSnooping:            n.IGMPSnooping,
		NetworkIsolationEnabled: n.NetworkIsolationEnabled,
		DHCPguardEnabled:        n.DHCPguardEnabled,
		DHCPDIP1:                n.DHCPDIP1,
		DHCPDIP2:                n.DHCPDIP2,
		DHCPDIP3:                n.DHCPDIP3,
	})
}

// marshalGuest marshals a Guest network.
func (n *Network) marshalGuest() ([]byte, error) {
	var defaultStart, defaultEnd string
	if n.IPSubnet != nil {
		var err error
		defaultStart, defaultEnd, err = dhcpRange(*n.IPSubnet)
		if err != nil {
			log.Default().Printf("error calculating DHCP range: %s", err)
		}
	}

	return json.Marshal(&struct {
		ID       string `json:"_id,omitempty"`
		SiteID   string `json:"site_id,omitempty"`
		Hidden   bool   `json:"attr_hidden,omitempty"`
		HiddenID string `json:"attr_hidden_id,omitempty"`
		NoDelete bool   `json:"attr_no_delete,omitempty"`
		NoEdit   bool   `json:"attr_no_edit,omitempty"`

		Name                    *string                         `json:"name,omitempty"`
		Purpose                 string                          `json:"purpose"`
		Enabled                 bool                            `json:"enabled"`
		NetworkGroup            *string                         `json:"networkgroup,omitempty"`
		IPSubnet                *string                         `json:"ip_subnet,omitempty"`
		VLAN                    *int64                          `json:"vlan,omitempty"`
		VLANEnabled             bool                            `json:"vlan_enabled"`
		DomainName              *string                         `json:"domain_name,omitempty"`
		AutoScaleEnabled        bool                            `json:"auto_scale_enabled"`
		GatewayType             *string                         `json:"gateway_type,omitempty"`
		InternetAccessEnabled   bool                            `json:"internet_access_enabled"`
		NetworkIsolationEnabled bool                            `json:"network_isolation_enabled"`
		SettingPreference       *string                         `json:"setting_preference,omitempty"`
		IGMPSnooping            bool                            `json:"igmp_snooping"`
		DHCPguardEnabled        bool                            `json:"dhcpguard_enabled"`
		MdnsEnabled             bool                            `json:"mdns_enabled"`
		LteLanEnabled           bool                            `json:"lte_lan_enabled"`
		IPAliases               []string                        `json:"ip_aliases"`
		NATOutboundIPAddresses  []NetworkNATOutboundIPAddresses `json:"nat_outbound_ip_addresses"`
		MACOverride             string                          `json:"mac_override,omitempty"`

		// DHCP Server
		DHCPDEnabled           bool    `json:"dhcpd_enabled"`
		DHCPDStart             *string `json:"dhcpd_start,omitempty"`
		DHCPDStop              *string `json:"dhcpd_stop,omitempty"`
		DHCPDLeaseTime         *int64  `json:"dhcpd_leasetime,omitempty"`
		DHCPDDNSEnabled        bool    `json:"dhcpd_dns_enabled"`
		DHCPDDNS1              string  `json:"dhcpd_dns_1"`
		DHCPDDNS2              string  `json:"dhcpd_dns_2"`
		DHCPDDNS3              string  `json:"dhcpd_dns_3"`
		DHCPDDNS4              string  `json:"dhcpd_dns_4"`
		DHCPDGatewayEnabled    bool    `json:"dhcpd_gateway_enabled"`
		DHCPDGateway           *string `json:"dhcpd_gateway,omitempty"`
		DHCPDNtpEnabled        bool    `json:"dhcpd_ntp_enabled"`
		DHCPDNtp1              *string `json:"dhcpd_ntp_1,omitempty"`
		DHCPDNtp2              *string `json:"dhcpd_ntp_2,omitempty"`
		DHCPDWinsEnabled       bool    `json:"dhcpd_wins_enabled"`
		DHCPDWins1             *string `json:"dhcpd_wins_1,omitempty"`
		DHCPDWins2             *string `json:"dhcpd_wins_2,omitempty"`
		DHCPDTimeOffsetEnabled bool    `json:"dhcpd_time_offset_enabled"`
		DHCPDConflictChecking  bool    `json:"dhcpd_conflict_checking"`
		DHCPDBootEnabled       bool    `json:"dhcpd_boot_enabled"`
		DHCPDBootServer        string  `json:"dhcpd_boot_server,omitempty"`
		DHCPDBootFilename      string  `json:"dhcpd_boot_filename,omitempty"`
		DHCPDTFTPServer        *string `json:"dhcpd_tftp_server,omitempty"`
		DHCPDWPAdUrl           *string `json:"dhcpd_wpad_url,omitempty"`
		DHCPDUnifiController   *string `json:"dhcpd_unifi_controller,omitempty"`

		// DHCP Relay
		DHCPRelayEnabled bool     `json:"dhcp_relay_enabled"`
		DHCPRelayServers []string `json:"dhcp_relay_servers"`

		// IPv6
		IPV6InterfaceType     *string `json:"ipv6_interface_type,omitempty"`
		IPV6SettingPreference *string `json:"ipv6_setting_preference,omitempty"`
		IPV6RaPriority        *string `json:"ipv6_ra_priority,omitempty"`

		// DHCPv6
		DHCPDV6DNSAuto    bool    `json:"dhcpdv6_dns_auto,omitempty"`
		DHCPDV6AllowSlaac bool    `json:"dhcpdv6_allow_slaac,omitempty"`
		DHCPDV6Start      *string `json:"dhcpdv6_start,omitempty"`
		DHCPDV6Stop       *string `json:"dhcpdv6_stop,omitempty"`
		DHCPDV6LeaseTime  *int64  `json:"dhcpdv6_leasetime,omitempty"`
	}{
		ID:       n.ID,
		SiteID:   n.SiteID,
		Hidden:   n.Hidden,
		HiddenID: n.HiddenID,
		NoDelete: n.NoDelete,
		NoEdit:   n.NoEdit,

		Name:                    n.Name,
		Purpose:                 n.Purpose,
		Enabled:                 n.Enabled,
		NetworkGroup:            valueOrDefault(n.NetworkGroup, "LAN"),
		IPSubnet:                valueOrDefault(n.IPSubnet, ""),
		VLAN:                    n.VLAN,
		VLANEnabled:             n.VLANEnabled,
		DomainName:              valueOrDefault(n.DomainName, ""),
		AutoScaleEnabled:        n.AutoScaleEnabled,
		GatewayType:             valueOrDefault(n.GatewayType, "default"),
		InternetAccessEnabled:   n.InternetAccessEnabled,
		NetworkIsolationEnabled: n.NetworkIsolationEnabled,
		SettingPreference:       valueOrDefault(n.SettingPreference, "auto"),
		IGMPSnooping:            n.IGMPSnooping,
		DHCPguardEnabled:        n.DHCPguardEnabled,
		MdnsEnabled:             n.MdnsEnabled,
		LteLanEnabled:           n.LteLanEnabled,
		IPAliases:               orEmptySlice(n.IPAliases),
		NATOutboundIPAddresses:  orEmptyNATSlice(n.NATOutboundIPAddresses),
		MACOverride:             n.MACOverride,

		// DHCP Server with defaults
		DHCPDEnabled:           n.DHCPDEnabled,
		DHCPDStart:             valueOrDefault(n.DHCPDStart, defaultStart),
		DHCPDStop:              valueOrDefault(n.DHCPDStop, defaultEnd),
		DHCPDLeaseTime:         valueOrDefault(n.DHCPDLeaseTime, 86400),
		DHCPDDNSEnabled:        n.DHCPDDNSEnabled,
		DHCPDDNS1:              n.DHCPDDNS1,
		DHCPDDNS2:              n.DHCPDDNS2,
		DHCPDDNS3:              n.DHCPDDNS3,
		DHCPDDNS4:              n.DHCPDDNS4,
		DHCPDGatewayEnabled:    n.DHCPDGatewayEnabled,
		DHCPDGateway:           n.DHCPDGateway,
		DHCPDNtpEnabled:        n.DHCPDNtpEnabled,
		DHCPDNtp1:              nilIfEmpty(n.DHCPDNtp1),
		DHCPDNtp2:              nilIfEmpty(n.DHCPDNtp2),
		DHCPDWinsEnabled:       n.DHCPDWinsEnabled,
		DHCPDWins1:             valueOrDefault(n.DHCPDWins1, ""),
		DHCPDWins2:             valueOrDefault(n.DHCPDWins2, ""),
		DHCPDTimeOffsetEnabled: n.DHCPDTimeOffsetEnabled,
		DHCPDConflictChecking:  n.DHCPDConflictChecking,
		DHCPDBootEnabled:       n.DHCPDBootEnabled,
		DHCPDBootServer:        n.DHCPDBootServer,
		DHCPDBootFilename:      derefOrEmpty(n.DHCPDBootFilename),
		DHCPDTFTPServer:        n.DHCPDTFTPServer,
		DHCPDWPAdUrl:           n.DHCPDWPAdUrl,
		DHCPDUnifiController:   valueOrDefault(n.DHCPDUnifiController, ""),

		// DHCP Relay
		DHCPRelayEnabled: n.DHCPRelayEnabled,
		DHCPRelayServers: orEmptySlice(n.DHCPRelayServers),

		// IPv6
		IPV6InterfaceType:     valueOrDefault(n.IPV6InterfaceType, "none"),
		IPV6SettingPreference: n.IPV6SettingPreference,
		IPV6RaPriority:        n.IPV6RaPriority,

		// DHCPv6
		DHCPDV6DNSAuto:    n.DHCPDV6DNSAuto,
		DHCPDV6AllowSlaac: n.DHCPDV6AllowSlaac,
		DHCPDV6Start:      n.DHCPDV6Start,
		DHCPDV6Stop:       n.DHCPDV6Stop,
		DHCPDV6LeaseTime:  n.DHCPDV6LeaseTime,
	})
}

// marshalWAN marshals a WAN network.
func (n *Network) marshalWAN() ([]byte, error) {
	return json.Marshal(&struct {
		ID       string `json:"_id,omitempty"`
		SiteID   string `json:"site_id,omitempty"`
		Hidden   bool   `json:"attr_hidden,omitempty"`
		HiddenID string `json:"attr_hidden_id,omitempty"`
		NoDelete bool   `json:"attr_no_delete,omitempty"`
		NoEdit   bool   `json:"attr_no_edit,omitempty"`

		Name    *string `json:"name,omitempty"`
		Purpose string  `json:"purpose"`
		Enabled bool    `json:"enabled"`

		// WAN type fields
		WANType         *string `json:"wan_type,omitempty"`
		WANTypeV6       *string `json:"wan_type_v6,omitempty"`
		WANNetworkGroup *string `json:"wan_networkgroup,omitempty"`

		// VLAN fields
		WANVLANEnabled bool   `json:"wan_vlan_enabled"`
		WANVLAN        *int64 `json:"wan_vlan,omitempty"`

		// DHCP CoS fields
		WANDHCPCos   *int64 `json:"wan_dhcp_cos,omitempty"`
		WANDHCPv6Cos *int64 `json:"wan_dhcpv6_cos,omitempty"`

		// DNS fields
		WANDNS1              *string `json:"wan_dns1,omitempty"`
		WANDNS2              *string `json:"wan_dns2,omitempty"`
		WANDNSPreference     *string `json:"wan_dns_preference,omitempty"`
		WANIPV6DNS1          *string `json:"wan_ipv6_dns1,omitempty"`
		WANIPV6DNS2          *string `json:"wan_ipv6_dns2,omitempty"`
		WANIPV6DNSPreference *string `json:"wan_ipv6_dns_preference,omitempty"`

		// DHCPv6 / IPv6 fields
		WANDHCPv6PDSize       *int64                    `json:"wan_dhcpv6_pd_size,omitempty"`
		WANDHCPv6PDSizeAuto   bool                      `json:"wan_dhcpv6_pd_size_auto"`
		WANDHCPv6Options      []NetworkWANDHCPv6Options `json:"wan_dhcpv6_options,omitempty"`
		IPV6WANDelegationType *string                   `json:"ipv6_wan_delegation_type,omitempty"`
		IPV6Enabled           bool                      `json:"ipv6_enabled"`

		// QoS fields
		WANEgressQOSEnabled *bool  `json:"wan_egress_qos_enabled,omitempty"`
		WANEgressQOS        *int64 `json:"wan_egress_qos,omitempty"`
		WANSmartQEnabled    bool   `json:"wan_smartq_enabled"`
		WANSmartQUpRate     *int64 `json:"wan_smartq_up_rate,omitempty"`
		WANSmartQDownRate   *int64 `json:"wan_smartq_down_rate,omitempty"`

		// UPnP fields
		UPnPEnabled       *bool   `json:"upnp_enabled,omitempty"`
		UPnPWANInterface  *string `json:"upnp_wan_interface,omitempty"`
		UPnPNatPMPEnabled *bool   `json:"upnp_nat_pmp_enabled,omitempty"`
		UPnPSecureMode    *bool   `json:"upnp_secure_mode,omitempty"`

		// Load balance / failover fields
		WANLoadBalanceType   *string `json:"wan_load_balance_type,omitempty"`
		WANLoadBalanceWeight *int64  `json:"wan_load_balance_weight,omitempty"`
		WANFailoverPriority  *int64  `json:"wan_failover_priority,omitempty"`

		// IGMP fields
		IGMPProxyFor      *string `json:"igmp_proxy_for,omitempty"`
		IGMPProxyUpstream bool    `json:"igmp_proxy_upstream"`

		// Event / alias fields
		ReportWANEvent bool                    `json:"report_wan_event"`
		WANIPAliases   []string                `json:"wan_ip_aliases"`
		WANDHCPOptions []NetworkWANDHCPOptions `json:"wan_dhcp_options"`

		// Provider capabilities
		WANProviderCapabilities *NetworkWANProviderCapabilities `json:"wan_provider_capabilities,omitempty"`
	}{
		ID:       n.ID,
		SiteID:   n.SiteID,
		Hidden:   n.Hidden,
		HiddenID: n.HiddenID,
		NoDelete: n.NoDelete,
		NoEdit:   n.NoEdit,

		Name:    n.Name,
		Purpose: n.Purpose,
		Enabled: n.Enabled,

		// WAN type fields
		WANType:         n.WANType,
		WANTypeV6:       n.WANTypeV6,
		WANNetworkGroup: n.WANNetworkGroup,

		// VLAN fields
		WANVLANEnabled: n.WANVLANEnabled,
		WANVLAN:        n.WANVLAN,

		// DHCP CoS fields
		WANDHCPCos:   n.WANDHCPCos,
		WANDHCPv6Cos: n.WANDHCPv6Cos,

		// DNS fields
		WANDNS1:              n.WANDNS1,
		WANDNS2:              n.WANDNS2,
		WANDNSPreference:     n.WANDNSPreference,
		WANIPV6DNS1:          n.WANIPV6DNS1,
		WANIPV6DNS2:          n.WANIPV6DNS2,
		WANIPV6DNSPreference: n.WANIPV6DNSPreference,

		// DHCPv6 / IPv6 fields
		WANDHCPv6PDSize:       n.WANDHCPv6PDSize,
		WANDHCPv6PDSizeAuto:   n.WANDHCPv6PDSizeAuto,
		WANDHCPv6Options:      n.WANDHCPv6Options,
		IPV6WANDelegationType: n.IPV6WANDelegationType,
		IPV6Enabled:           n.WANTypeV6 != nil && *n.WANTypeV6 != "disabled",

		// QoS fields
		WANEgressQOSEnabled: n.WANEgressQOSEnabled,
		WANEgressQOS:        n.WANEgressQOS,
		WANSmartQEnabled:    n.WANSmartQEnabled,
		WANSmartQUpRate:     n.WANSmartQUpRate,
		WANSmartQDownRate:   n.WANSmartQDownRate,

		// UPnP fields
		UPnPEnabled:       n.UPnPEnabled,
		UPnPWANInterface:  n.UPnPWANInterface,
		UPnPNatPMPEnabled: n.UPnPNatPMPEnabled,
		UPnPSecureMode:    n.UPnPSecureMode,

		// Load balance / failover fields
		WANLoadBalanceType:   n.WANLoadBalanceType,
		WANLoadBalanceWeight: n.WANLoadBalanceWeight,
		WANFailoverPriority:  n.WANFailoverPriority,

		// IGMP fields
		IGMPProxyFor:      n.IGMPProxyFor,
		IGMPProxyUpstream: n.IGMPProxyUpstream,

		// Event / alias fields
		ReportWANEvent: n.ReportWANEvent,
		WANIPAliases:   orEmptySlice(n.WANIPAliases),
		WANDHCPOptions: orEmptyWANDHCPOptions(n.WANDHCPOptions),

		// Provider capabilities
		WANProviderCapabilities: n.WANProviderCapabilities,
	})
}

// marshalSiteVPN marshals a site-to-site VPN network.
func (n *Network) marshalSiteVPN() ([]byte, error) {
	return json.Marshal(&struct {
		ID       string `json:"_id,omitempty"`
		SiteID   string `json:"site_id,omitempty"`
		Hidden   bool   `json:"attr_hidden,omitempty"`
		HiddenID string `json:"attr_hidden_id,omitempty"`
		NoDelete bool   `json:"attr_no_delete,omitempty"`
		NoEdit   bool   `json:"attr_no_edit,omitempty"`

		Name    *string `json:"name,omitempty"`
		Purpose string  `json:"purpose"`
		Enabled bool    `json:"enabled"`
	}{
		ID:       n.ID,
		SiteID:   n.SiteID,
		Hidden:   n.Hidden,
		HiddenID: n.HiddenID,
		NoDelete: n.NoDelete,
		NoEdit:   n.NoEdit,

		Name:    n.Name,
		Purpose: n.Purpose,
		Enabled: n.Enabled,
	})
}

// marshalVPNClient marshals a VPN client network (WireGuard client).
func (n *Network) marshalVPNClient() ([]byte, error) {
	return json.Marshal(&struct {
		ID       string `json:"_id,omitempty"`
		SiteID   string `json:"site_id,omitempty"`
		Hidden   bool   `json:"attr_hidden,omitempty"`
		HiddenID string `json:"attr_hidden_id,omitempty"`
		NoDelete bool   `json:"attr_no_delete,omitempty"`
		NoEdit   bool   `json:"attr_no_edit,omitempty"`

		Name     *string `json:"name,omitempty"`
		Purpose  string  `json:"purpose"`
		Enabled  bool    `json:"enabled"`
		IPSubnet *string `json:"ip_subnet,omitempty"`

		// VPN Type
		VPNType *string `json:"vpn_type,omitempty"`

		// VPN Client routing
		VPNClientDefaultRoute bool `json:"vpn_client_default_route"`
		VPNClientPullDNS      bool `json:"vpn_client_pull_dns"`

		// WireGuard Client Configuration
		WireguardClientMode                  *string `json:"wireguard_client_mode,omitempty"`
		WireguardClientConfigurationFile     *string `json:"wireguard_client_configuration_file,omitempty"`
		WireguardClientConfigurationFilename *string `json:"wireguard_client_configuration_filename,omitempty"`
		WireguardClientPeerIP                *string `json:"wireguard_client_peer_ip,omitempty"`
		WireguardClientPeerPort              *int64  `json:"wireguard_client_peer_port,omitempty"`
		WireguardClientPeerPublicKey         *string `json:"wireguard_client_peer_public_key,omitempty"`
		WireguardClientPresharedKeyEnabled   bool    `json:"wireguard_client_preshared_key_enabled"`
		WireguardClientPresharedKey          *string `json:"wireguard_client_preshared_key,omitempty"`
		WireguardInterface                   *string `json:"wireguard_interface,omitempty"`
		WireguardPrivateKey                  *string `json:"x_wireguard_private_key,omitempty"`

		// DNS servers for WireGuard interface
		DHCPDDNS1 string `json:"dhcpd_dns_1"`
		DHCPDDNS2 string `json:"dhcpd_dns_2"`
	}{
		ID:       n.ID,
		SiteID:   n.SiteID,
		Hidden:   n.Hidden,
		HiddenID: n.HiddenID,
		NoDelete: n.NoDelete,
		NoEdit:   n.NoEdit,

		Name:     n.Name,
		Purpose:  n.Purpose,
		Enabled:  n.Enabled,
		IPSubnet: n.IPSubnet,

		// VPN Type
		VPNType: n.VPNType,

		// VPN Client routing
		VPNClientDefaultRoute: n.VPNClientDefaultRoute,
		VPNClientPullDNS:      n.VPNClientPullDNS,

		// WireGuard configuration
		WireguardClientMode:                  n.WireguardClientMode,
		WireguardClientConfigurationFile:     n.WireguardClientConfigurationFile,
		WireguardClientConfigurationFilename: n.WireguardClientConfigurationFilename,
		WireguardClientPeerIP:                n.WireguardClientPeerIP,
		WireguardClientPeerPort:              n.WireguardClientPeerPort,
		WireguardClientPeerPublicKey:         n.WireguardClientPeerPublicKey,
		WireguardClientPresharedKeyEnabled:   n.WireguardClientPresharedKeyEnabled,
		WireguardClientPresharedKey:          n.WireguardClientPresharedKey,
		WireguardInterface:                   n.WireguardInterface,
		WireguardPrivateKey:                  n.WireguardPrivateKey,

		// DNS servers
		DHCPDDNS1: n.DHCPDDNS1,
		DHCPDDNS2: n.DHCPDDNS2,
	})
}

// marshalUserVPN marshals a remote user VPN network.
func (n *Network) marshalUserVPN() ([]byte, error) {
	return json.Marshal(&struct {
		ID       string `json:"_id,omitempty"`
		SiteID   string `json:"site_id,omitempty"`
		Hidden   bool   `json:"attr_hidden,omitempty"`
		HiddenID string `json:"attr_hidden_id,omitempty"`
		NoDelete bool   `json:"attr_no_delete,omitempty"`
		NoEdit   bool   `json:"attr_no_edit,omitempty"`

		Name              *string `json:"name,omitempty"`
		Purpose           string  `json:"purpose"`
		Enabled           bool    `json:"enabled"`
		SettingPreference *string `json:"setting_preference,omitempty"`
		IPSubnet          *string `json:"ip_subnet,omitempty"`

		// VPN Type
		VPNType *string `json:"vpn_type,omitempty"`

		// DNS
		DHCPDDNS1       string `json:"dhcpd_dns_1"`
		DHCPDDNS2       string `json:"dhcpd_dns_2"`
		DHCPDDNSEnabled bool   `json:"dhcpd_dns_enabled"`

		// DHCP Range
		DHCPDStart *string `json:"dhcpd_start,omitempty"`
		DHCPDStop  *string `json:"dhcpd_stop,omitempty"`

		// RADIUS
		RADIUSProfileID *string `json:"radiusprofile_id,omitempty"`

		// WireGuard Server Configuration
		WireguardInterface                     *string `json:"wireguard_interface,omitempty"`
		WireguardPrivateKey                    *string `json:"x_wireguard_private_key,omitempty"`
		WireguardLocalWANIP                    *string `json:"wireguard_local_wan_ip,omitempty"`
		LocalPort                              *int64  `json:"local_port,omitempty"`
		WireguardInterfaceBindingModeIPVersion *string `json:"wireguard_interface_binding_mode_ip_version,omitempty"`
		VPNClientConfigurationRemoteIPOverride *string `json:"vpn_client_configuration_remote_ip_override,omitempty"`

		// L2TP Server Configuration
		L2TpInterface        *string `json:"l2tp_interface,omitempty"`
		L2TpLocalWANIP       *string `json:"l2tp_local_wan_ip,omitempty"`
		L2TpAllowWeakCiphers bool    `json:"l2tp_allow_weak_ciphers"`
		IPSecPreSharedKey    *string `json:"x_ipsec_pre_shared_key,omitempty"`

		// OpenVPN Server Configuration
		OpenVPNInterface        *string `json:"openvpn_interface,omitempty"`
		OpenVPNLocalWANIP       *string `json:"openvpn_local_wan_ip,omitempty"`
		OpenVPNMode             *string `json:"openvpn_mode,omitempty"`
		OpenVPNEncryptionCipher *string `json:"openvpn_encryption_cipher,omitempty"`

		// OpenVPN Certificates and Keys
		ServerCrt       *string `json:"x_server_crt,omitempty"`
		ServerKey       *string `json:"x_server_key,omitempty"`
		DhKey           *string `json:"x_dh_key,omitempty"`
		SharedClientKey *string `json:"x_shared_client_key,omitempty"`
		SharedClientCrt *string `json:"x_shared_client_crt,omitempty"`
		AuthKey         *string `json:"x_auth_key,omitempty"`
		CaCrt           *string `json:"x_ca_crt,omitempty"`
		CaKey           *string `json:"x_ca_key,omitempty"`
	}{
		ID:       n.ID,
		SiteID:   n.SiteID,
		Hidden:   n.Hidden,
		HiddenID: n.HiddenID,
		NoDelete: n.NoDelete,
		NoEdit:   n.NoEdit,

		Name:              n.Name,
		Purpose:           n.Purpose,
		Enabled:           n.Enabled,
		SettingPreference: n.SettingPreference,
		IPSubnet:          n.IPSubnet,

		// VPN Type
		VPNType: n.VPNType,

		// DNS
		DHCPDDNS1:       n.DHCPDDNS1,
		DHCPDDNS2:       n.DHCPDDNS2,
		DHCPDDNSEnabled: n.DHCPDDNSEnabled,

		// DHCP Range
		DHCPDStart: n.DHCPDStart,
		DHCPDStop:  n.DHCPDStop,

		// RADIUS
		RADIUSProfileID: n.RADIUSProfileID,

		// WireGuard Server Configuration
		WireguardInterface:                     n.WireguardInterface,
		WireguardPrivateKey:                    n.WireguardPrivateKey,
		WireguardLocalWANIP:                    n.WireguardLocalWANIP,
		LocalPort:                              n.LocalPort,
		WireguardInterfaceBindingModeIPVersion: n.WireguardInterfaceBindingModeIPVersion,
		VPNClientConfigurationRemoteIPOverride: n.VPNClientConfigurationRemoteIPOverride,

		// L2TP Server Configuration
		L2TpInterface:        n.L2TpInterface,
		L2TpLocalWANIP:       n.L2TpLocalWANIP,
		L2TpAllowWeakCiphers: n.L2TpAllowWeakCiphers,
		IPSecPreSharedKey:    n.IPSecPreSharedKey,

		// OpenVPN Server Configuration
		OpenVPNInterface:        n.OpenVPNInterface,
		OpenVPNLocalWANIP:       n.OpenVPNLocalWANIP,
		OpenVPNMode:             n.OpenVPNMode,
		OpenVPNEncryptionCipher: n.OpenVPNEncryptionCipher,

		// OpenVPN Certificates and Keys
		ServerCrt:       n.ServerCrt,
		ServerKey:       n.ServerKey,
		DhKey:           n.DhKey,
		SharedClientKey: n.SharedClientKey,
		SharedClientCrt: n.SharedClientCrt,
		AuthKey:         n.AuthKey,
		CaCrt:           n.CaCrt,
		CaKey:           n.CaKey,
	})
}

// Helper functions for field transformations

func orEmptySlice(s []string) []string {
	if len(s) > 0 {
		return s
	}
	return []string{}
}

func orEmptyNATSlice(s []NetworkNATOutboundIPAddresses) []NetworkNATOutboundIPAddresses {
	if len(s) > 0 {
		return s
	}
	return []NetworkNATOutboundIPAddresses{}
}

func orEmptyWANDHCPOptions(s []NetworkWANDHCPOptions) []NetworkWANDHCPOptions {
	if len(s) > 0 {
		return s
	}
	return []NetworkWANDHCPOptions{}
}

func nilIfEmpty(s *string) *string {
	if s != nil && *s == "" {
		return nil
	}
	return s
}

func derefOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func dhcpRange(cidr string) (start, end string, err error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return "", "", err
	}

	// Only support IPv4
	if !prefix.Addr().Is4() {
		return "", "", fmt.Errorf("only IPv4 supported")
	}

	networkAddr := prefix.Masked().Addr()
	bits := prefix.Bits()

	// Calculate the number of host addresses
	hostBits := 32 - bits
	numHosts := uint32(1) << hostBits

	// UniFi's rules based on subnet size:
	// /30 or smaller (4 or fewer IPs): No DHCP (too small)
	// /29 (8 IPs): Start at +2, End at -2 (gives 4 usable IPs)
	// /28 to /24: Start at +6, End at -1 (broadcast)
	// /23 and larger: Start at +6, End at -1

	if bits >= 30 {
		return "", "", fmt.Errorf("subnet too small for DHCP (/%d)", bits)
	}

	// Convert network address to uint32 for arithmetic
	ip4 := networkAddr.As4()
	baseIP := uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])

	var startOffset, endOffset uint32

	if bits == 29 {
		// /29: 8 IPs total
		// Network: .0, Gateway: .1, DHCP: .2-.5, Reserved: .6, Broadcast: .7
		startOffset = 2
		endOffset = 2
	} else {
		// /28 and larger
		// Network: .0, Gateway: .1, Reserved: .2-.5, DHCP: .6 to (broadcast-1)
		startOffset = 6
		endOffset = 1
	}

	startIP := baseIP + startOffset
	endIP := baseIP + numHosts - 1 - endOffset

	// Convert back to netip.Addr
	start = netip.AddrFrom4([4]byte{
		byte(startIP >> 24),
		byte(startIP >> 16),
		byte(startIP >> 8),
		byte(startIP),
	}).String()

	end = netip.AddrFrom4([4]byte{
		byte(endIP >> 24),
		byte(endIP >> 16),
		byte(endIP >> 8),
		byte(endIP),
	}).String()

	return start, end, nil
}

func valueOrDefault[T any](in *T, defaultValue T) *T {
	if in == nil {
		return &defaultValue
	}
	return in
}
