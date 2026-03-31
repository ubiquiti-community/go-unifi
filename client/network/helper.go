package network

func (p TrafficMatchingListsPage) GetTotalCount() int64           { return p.TotalCount }
func (p TrafficMatchingListsPage) GetData() []TrafficMatchingList { return p.Data }

func (p IntegrationDnsPolicyPageDto) GetTotalCount() int64 { return p.TotalCount }
func (p IntegrationDnsPolicyPageDto) GetData() []DNSPolicy { return p.Data }

func (p IntegrationAclRulePageDto) GetTotalCount() int64     { return p.TotalCount }
func (p IntegrationAclRulePageDto) GetData() []ACLRuleObject { return p.Data }

func (p FirewallZonesPage) GetTotalCount() int64    { return p.TotalCount }
func (p FirewallZonesPage) GetData() []FirewallZone { return p.Data }

func (p FirewallPolicyPage) GetTotalCount() int64      { return p.TotalCount }
func (p FirewallPolicyPage) GetData() []FirewallPolicy { return p.Data }

func (p HotspotVoucherDetailPage) GetTotalCount() int64             { return p.TotalCount }
func (p HotspotVoucherDetailPage) GetData() []HotspotVoucherDetails { return p.Data }

func (p IntegrationWifiBroadcastPageDto) GetTotalCount() int64             { return p.TotalCount }
func (p IntegrationWifiBroadcastPageDto) GetData() []WifiBroadcastOverview { return p.Data }

func (p NetworkOverviewPage) GetTotalCount() int64       { return p.TotalCount }
func (p NetworkOverviewPage) GetData() []NetworkOverview { return p.Data }

func (p ClientOverviewPage) GetTotalCount() int64      { return p.TotalCount }
func (p ClientOverviewPage) GetData() []ClientOverview { return p.Data }

func (p AdoptedDeviceOverviewPage) GetTotalCount() int64             { return p.TotalCount }
func (p AdoptedDeviceOverviewPage) GetData() []AdoptedDeviceOverview { return p.Data }

func (p DevicePendingAdoptionPage) GetTotalCount() int64             { return p.TotalCount }
func (p DevicePendingAdoptionPage) GetData() []DevicePendingAdoption { return p.Data }

func (p SiteOverviewPage) GetTotalCount() int64    { return p.TotalCount }
func (p SiteOverviewPage) GetData() []SiteOverview { return p.Data }

func (p IntegrationDeviceTagPageDto) GetTotalCount() int64 { return p.TotalCount }
func (p IntegrationDeviceTagPageDto) GetData() []DeviceTag { return p.Data }

func (p WANOverviewPage) GetTotalCount() int64 { return p.TotalCount }
func (p WANOverviewPage) GetData() []WANOverview { return p.Data }

func (p IntegrationSiteToSiteVpnTunnelOverviewPageDto) GetTotalCount() int64 { return p.TotalCount }
func (p IntegrationSiteToSiteVpnTunnelOverviewPageDto) GetData() []SiteToSiteVPNTunnelOverview { return p.Data }

func (p IntegrationVpnServerOverviewPageDto) GetTotalCount() int64 { return p.TotalCount }
func (p IntegrationVpnServerOverviewPageDto) GetData() []VPNServerOverview { return p.Data }

func (p RadiusProfileOverviewPage) GetTotalCount() int64 { return p.TotalCount }
func (p RadiusProfileOverviewPage) GetData() []RadiusProfileOverview { return p.Data }

func (p DPICategoryPage) GetTotalCount() int64 { return p.TotalCount }
func (p DPICategoryPage) GetData() []DPICategory { return p.Data }

func (p DPIApplicationPage) GetTotalCount() int64 { return p.TotalCount }
func (p DPIApplicationPage) GetData() []DPIApplication { return p.Data }

func (p CountryDefinitionPage) GetTotalCount() int64 { return p.TotalCount }
func (p CountryDefinitionPage) GetData() []CountryDefinition { return p.Data }
