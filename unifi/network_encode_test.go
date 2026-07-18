package unifi

import (
	"encoding/json"
	"testing"
)

// Helper function to parse JSON and check for expected/unexpected fields.
func checkJSONFields(t *testing.T, data []byte, expectedFields []string, unexpectedFields []string) {
	t.Helper()

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check expected fields are present
	for _, field := range expectedFields {
		if _, ok := result[field]; !ok {
			t.Errorf("Expected field %q not found in JSON", field)
		}
	}

	// Check unexpected fields are absent
	for _, field := range unexpectedFields {
		if _, ok := result[field]; ok {
			t.Errorf("Unexpected field %q found in JSON", field)
		}
	}
}

func TestMarshalNetworkCorporate(t *testing.T) {
	// Create a corporate network with common fields
	vlan := int64(10)
	leasetime := int64(86400)
	dhcpGateway := "192.168.1.1"
	dhcpStart := "192.168.1.100"
	dhcpStop := "192.168.1.200"

	network := &Network{
		ID:                    "507f1f77bcf86cd799439011",
		SiteID:                "default",
		Name:                  strPtr("Corporate LAN"),
		Purpose:               PurposeCorporate,
		Enabled:               true,
		AutoScaleEnabled:      false,
		NetworkGroup:          strPtr("LAN"),
		IPSubnet:              strPtr("192.168.1.0/24"),
		VLAN:                  &vlan,
		VLANEnabled:           true,
		DomainName:            strPtr("example.local"),
		GatewayType:           strPtr("default"),
		DHCPDGateway:          &dhcpGateway,
		DHCPDGatewayEnabled:   true,
		InternetAccessEnabled: true,
		IGMPSnooping:          false,
		DHCPDEnabled:          true,
		DHCPDStart:            &dhcpStart,
		DHCPDStop:             &dhcpStop,
		DHCPDLeaseTime:        &leasetime,
		DHCPDDNS1:             strPtr("8.8.8.8"),
		DHCPDDNS2:             strPtr("8.8.4.4"),
		DHCPDDNSEnabled:       true,
		IPAliases:             []string{},
	}

	// Marshal to JSON
	data, err := json.Marshal(network)
	if err != nil {
		t.Fatalf("Failed to marshal network: %v", err)
	}

	// Expected fields for corporate network
	expectedFields := []string{
		"_id",
		"site_id",
		"name",
		"purpose",
		"enabled",
		"networkgroup",
		"ip_subnet",
		"vlan",
		"vlan_enabled",
		"domain_name",
		"gateway_type",
		"dhcpd_gateway",
		"dhcpd_gateway_enabled",
		"internet_access_enabled",
		"igmp_snooping",
		"dhcpd_enabled",
		"dhcpd_start",
		"dhcpd_stop",
		"dhcpd_leasetime",
		"dhcpd_dns_1",
		"dhcpd_dns_2",
		"dhcpd_dns_enabled",
		"ip_aliases",
		"auto_scale_enabled",
		"setting_preference",
	}

	// Unexpected fields (WAN-specific)
	unexpectedFields := []string{
		"wan_type",
		"wan_ip",
		"wan_networkgroup",
		"ipsec_key_exchange",
		"wireguard_interface",
	}

	checkJSONFields(t, data, expectedFields, unexpectedFields)

	// Verify purpose is correct
	var result map[string]any
	json.Unmarshal(data, &result)
	if result["purpose"] != string(PurposeCorporate) {
		t.Errorf("Expected purpose %q, got %q", PurposeCorporate, result["purpose"])
	}

	// Verify default values are applied
	if result["networkgroup"] != "LAN" {
		t.Errorf("Expected networkgroup 'LAN', got %q", result["networkgroup"])
	}
	if result["gateway_type"] != "default" {
		t.Errorf("Expected gateway_type 'default', got %q", result["gateway_type"])
	}
	if result["setting_preference"] != "auto" {
		t.Errorf("Expected setting_preference 'auto', got %q", result["setting_preference"])
	}
}

func TestMarshalNetworkCorporateDefaults(t *testing.T) {
	// Create a minimal corporate network to test defaults
	network := &Network{
		ID:      "507f1f77bcf86cd799439011",
		Purpose: PurposeCorporate,
		Enabled: true,
	}

	data, err := json.Marshal(network)
	if err != nil {
		t.Fatalf("Failed to marshal network: %v", err)
	}

	var result map[string]any
	json.Unmarshal(data, &result)

	// Verify defaults are applied
	if result["networkgroup"] != "LAN" {
		t.Errorf("Expected default networkgroup 'LAN', got %v", result["networkgroup"])
	}
	if result["gateway_type"] != "default" {
		t.Errorf("Expected default gateway_type 'default', got %v", result["gateway_type"])
	}
	if result["setting_preference"] != "auto" {
		t.Errorf("Expected default setting_preference 'auto', got %v", result["setting_preference"])
	}
	if result["ip_subnet"] != "" {
		t.Errorf("Expected empty ip_subnet, got %v", result["ip_subnet"])
	}

	// Verify empty arrays are empty, not nil
	if aliases, ok := result["ip_aliases"].([]any); !ok || len(aliases) != 0 {
		t.Errorf("Expected empty array for ip_aliases, got %v", result["ip_aliases"])
	}
}

func TestMarshalNetworkWAN(t *testing.T) {
	vlan := int64(20)
	failoverPriority := int64(1)
	loadBalanceWeight := int64(50)
	dhcpv6PDSize := int64(56)

	network := &Network{
		ID:                    "507f1f77bcf86cd799439012",
		SiteID:                "default",
		Name:                  strPtr("WAN"),
		Purpose:               PurposeWAN,
		Enabled:               true,
		WANType:               strPtr("dhcp"),
		WANTypeV6:             strPtr("dhcpv6"),
		WANNetworkGroup:       strPtr("WAN"),
		WANVLANEnabled:        true,
		WANVLAN:               &vlan,
		WANDNSPreference:      strPtr("auto"),
		WANIPV6DNSPreference:  strPtr("auto"),
		WANDHCPv6PDSize:       &dhcpv6PDSize,
		WANDHCPv6PDSizeAuto:   true,
		IPV6WANDelegationType: strPtr("pd"),
		WANLoadBalanceType:    strPtr("failover-only"),
		WANLoadBalanceWeight:  &loadBalanceWeight,
		WANFailoverPriority:   &failoverPriority,
		IGMPProxyFor:          strPtr("none"),
		IGMPProxyUpstream:     false,
		ReportWANEvent:        true,
		WANIPAliases:          []string{},
		WANDHCPOptions:        []NetworkWANDHCPOptions{},
	}

	data, err := json.Marshal(network)
	if err != nil {
		t.Fatalf("Failed to marshal network: %v", err)
	}

	expectedFields := []string{
		"_id",
		"site_id",
		"name",
		"purpose",
		"enabled",
		"wan_type",
		"wan_type_v6",
		"wan_networkgroup",
		"wan_vlan_enabled",
		"wan_vlan",
		"wan_dns_preference",
		"wan_ipv6_dns_preference",
		"wan_dhcpv6_pd_size",
		"wan_dhcpv6_pd_size_auto",
		"ipv6_wan_delegation_type",
		"wan_load_balance_type",
		"wan_load_balance_weight",
		"wan_failover_priority",
		"igmp_proxy_for",
		"igmp_proxy_upstream",
		"report_wan_event",
		"wan_ip_aliases",
		"wan_dhcp_options",
		"ipv6_enabled",
	}

	unexpectedFields := []string{
		"networkgroup",
		"ip_subnet",
		"vlan",
		"dhcpd_enabled",
		"ipsec_interface",
		"wireguard_interface",
	}

	checkJSONFields(t, data, expectedFields, unexpectedFields)

	var result map[string]any
	json.Unmarshal(data, &result)

	// Verify WAN-specific values
	if result["purpose"] != string(PurposeWAN) {
		t.Errorf("Expected purpose %q, got %q", PurposeWAN, result["purpose"])
	}
	if result["ipv6_enabled"] != true {
		t.Errorf("Expected ipv6_enabled true, got %v", result["ipv6_enabled"])
	}

	// Verify empty arrays
	if aliases, ok := result["wan_ip_aliases"].([]any); !ok || len(aliases) != 0 {
		t.Errorf("Expected empty array for wan_ip_aliases, got %v", result["wan_ip_aliases"])
	}
}

func TestMarshalNetworkUnknownPurpose(t *testing.T) {
	network := &Network{
		ID:      "507f1f77bcf86cd799439016",
		Purpose: "unknown-purpose",
		Enabled: true,
	}

	_, err := json.Marshal(network)
	if err == nil {
		t.Error("Expected error for unknown purpose, got nil")
	}
}

func TestMarshalNetworkVLANOnly(t *testing.T) {
	vlan := int64(92)

	network := &Network{
		ID:      "507f1f77bcf86cd799439017",
		SiteID:  "default",
		Name:    strPtr("VLAN_92"),
		Purpose: PurposeVLANOnly,
		VLAN:    &vlan,
	}

	data, err := json.Marshal(network)
	if err != nil {
		t.Fatalf("Failed to marshal vlan-only network: %v", err)
	}

	expectedFields := []string{
		"_id",
		"site_id",
		"name",
		"purpose",
		"enabled",
		"networkgroup",
		"vlan",
		"vlan_enabled",
	}

	unexpectedFields := []string{
		"ip_subnet",
		"dhcpd_enabled",
		"wan_type",
		"wireguard_interface",
	}

	checkJSONFields(t, data, expectedFields, unexpectedFields)

	var result map[string]any
	json.Unmarshal(data, &result)

	if result["purpose"] != "vlan-only" {
		t.Errorf("Expected purpose 'vlan-only', got %q", result["purpose"])
	}
	if result["enabled"] != true {
		t.Errorf("Expected enabled true (default), got %v", result["enabled"])
	}
	if result["vlan_enabled"] != true {
		t.Errorf("Expected vlan_enabled true (auto-set from VLAN ID), got %v", result["vlan_enabled"])
	}
	if result["networkgroup"] != "LAN" {
		t.Errorf("Expected networkgroup 'LAN', got %v", result["networkgroup"])
	}
}

func TestMarshalNetworkVLANOnlyMinimal(t *testing.T) {
	network := &Network{
		ID:      "507f1f77bcf86cd799439018",
		Purpose: PurposeVLANOnly,
	}

	data, err := json.Marshal(network)
	if err != nil {
		t.Fatalf("Failed to marshal minimal vlan-only network: %v", err)
	}

	var result map[string]any
	json.Unmarshal(data, &result)

	if result["purpose"] != "vlan-only" {
		t.Errorf("Expected purpose 'vlan-only', got %q", result["purpose"])
	}
	if result["enabled"] != true {
		t.Errorf("Expected enabled true (default), got %v", result["enabled"])
	}
	if result["vlan_enabled"] != false {
		t.Errorf("Expected vlan_enabled false (no VLAN ID), got %v", result["vlan_enabled"])
	}
}

func TestMarshalNetworkGuest(t *testing.T) {
	vlan := int64(100)
	leasetime := int64(86400)
	dhcpStart := "192.168.100.100"
	dhcpStop := "192.168.100.200"

	network := &Network{
		ID:                    "507f1f77bcf86cd799439019",
		SiteID:                "default",
		Name:                  strPtr("Guest Network"),
		Purpose:               PurposeGuest,
		Enabled:               true,
		NetworkGroup:          strPtr("LAN"),
		IPSubnet:              strPtr("192.168.100.0/24"),
		VLAN:                  &vlan,
		VLANEnabled:           true,
		InternetAccessEnabled: true,
		DHCPDEnabled:          true,
		DHCPDStart:            &dhcpStart,
		DHCPDStop:             &dhcpStop,
		DHCPDLeaseTime:        &leasetime,
		DHCPDDNSEnabled:       true,
		DHCPDDNS1:             strPtr("8.8.8.8"),
	}

	data, err := json.Marshal(network)
	if err != nil {
		t.Fatalf("Failed to marshal guest network: %v", err)
	}

	expectedFields := []string{
		"_id",
		"site_id",
		"name",
		"purpose",
		"enabled",
		"networkgroup",
		"ip_subnet",
		"vlan",
		"vlan_enabled",
		"internet_access_enabled",
		"dhcpd_enabled",
		"dhcpd_start",
		"dhcpd_stop",
		"dhcpd_leasetime",
		"dhcpd_dns_enabled",
		"dhcpd_dns_1",
		"ip_aliases",
		"setting_preference",
	}

	unexpectedFields := []string{
		"wan_type",
		"wan_networkgroup",
		"wireguard_interface",
	}

	checkJSONFields(t, data, expectedFields, unexpectedFields)

	var result map[string]any
	json.Unmarshal(data, &result)

	if result["purpose"] != "guest" {
		t.Errorf("Expected purpose 'guest', got %q", result["purpose"])
	}
	if result["networkgroup"] != "LAN" {
		t.Errorf("Expected networkgroup 'LAN', got %v", result["networkgroup"])
	}
	if result["setting_preference"] != "auto" {
		t.Errorf("Expected setting_preference 'auto', got %v", result["setting_preference"])
	}
}

// TestMarshalNetworkIPv6ClientAddressAssignment guards that the corporate and
// guest marshalers emit ipv6_client_address_assignment when set, and omit it
// when nil. The field lives on the generated Network struct but the marshalers
// only serialize a curated subset, so it would otherwise be silently dropped on
// write (ubiquiti-community/terraform-provider-unifi#232).
func TestMarshalNetworkIPv6ClientAddressAssignment(t *testing.T) {
	for _, purpose := range []string{PurposeCorporate, PurposeGuest} {
		t.Run(purpose, func(t *testing.T) {
			// Set => emitted with the configured value.
			network := &Network{
				ID:                          "507f1f77bcf86cd799439011",
				Purpose:                     purpose,
				Enabled:                     true,
				IPV6InterfaceType:           strPtr("static"),
				IPV6ClientAddressAssignment: strPtr("slaac-dhcpv6"),
			}
			data, err := json.Marshal(network)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var result map[string]any
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got := result["ipv6_client_address_assignment"]; got != "slaac-dhcpv6" {
				t.Errorf("ipv6_client_address_assignment = %v, want slaac-dhcpv6", got)
			}

			// Unset => omitted (omitempty), no perpetual diff against the API.
			data, err = json.Marshal(&Network{ID: "x", Purpose: purpose, Enabled: true})
			if err != nil {
				t.Fatalf("marshal (unset): %v", err)
			}
			result = map[string]any{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("unmarshal (unset): %v", err)
			}
			if _, ok := result["ipv6_client_address_assignment"]; ok {
				t.Errorf("ipv6_client_address_assignment serialized for nil value: %s", data)
			}
		})
	}
}

// TestMarshalNetworkSiteVPN guards that the site-to-site IPsec VPN marshaler
// emits the VPN/IPsec fields (not just name/purpose/enabled). It was previously
// a stub, which silently dropped the whole tunnel configuration on write
// (ubiquiti-community/terraform-provider-unifi#78).
func TestMarshalNetworkSiteVPN(t *testing.T) {
	dh := int64(14)
	network := &Network{
		ID:                "507f1f77bcf86cd799439011",
		Name:              strPtr("HQ-to-Branch"),
		Purpose:           PurposeSiteVPN,
		Enabled:           true,
		VPNType:           strPtr("ipsec-vpn"),
		IPSecInterface:    strPtr("wan"),
		IPSecPeerIP:       strPtr("203.0.113.9"),
		IPSecKeyExchange:  strPtr("ikev2"),
		IPSecPreSharedKey: strPtr("s3cret-psk"),
		IPSecProfile:      strPtr("customized"),
		IPSecEncryption:   strPtr("aes256"),
		IPSecHash:         strPtr("sha256"),
		IPSecDhGroup:      &dh,
		IPSecPfs:          true,
		RemoteVPNSubnets:  []string{"192.0.2.0/24", "198.51.100.0/24"},
	}

	data, err := json.Marshal(network)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	checkJSONFields(t, data, []string{
		"name", "purpose", "enabled", "vpn_type", "ipsec_interface",
		"ipsec_peer_ip", "ipsec_key_exchange", "x_ipsec_pre_shared_key",
		"ipsec_profile", "ipsec_encryption", "ipsec_hash", "ipsec_dh_group",
		"ipsec_pfs", "remote_vpn_subnets",
	}, []string{"ip_subnet", "dhcpd_enabled", "vlan"})

	if result["purpose"] != PurposeSiteVPN {
		t.Errorf("purpose = %v, want %q", result["purpose"], PurposeSiteVPN)
	}
	if result["vpn_type"] != "ipsec-vpn" {
		t.Errorf("vpn_type = %v, want ipsec-vpn", result["vpn_type"])
	}
	if result["x_ipsec_pre_shared_key"] != "s3cret-psk" {
		t.Errorf("x_ipsec_pre_shared_key = %v", result["x_ipsec_pre_shared_key"])
	}
	subnets, ok := result["remote_vpn_subnets"].([]any)
	if !ok || len(subnets) != 2 {
		t.Errorf("remote_vpn_subnets = %v, want 2 entries", result["remote_vpn_subnets"])
	}
}

// Helper function to create string pointers.
func strPtr(s string) *string {
	return &s
}
