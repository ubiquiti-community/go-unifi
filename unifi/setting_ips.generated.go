// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type SettingIps struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	AdvancedFilteringPreference         string                `json:"advanced_filtering_preference,omitempty"` // |manual|disabled
	ContentFilteringBlockingPageEnabled bool                  `json:"content_filtering_blocking_page_enabled"`
	EnabledCategories                   []string              `json:"enabled_categories,omitempty"` // emerging-activex|emerging-attackresponse|botcc|emerging-chat|ciarmy|compromised|emerging-dns|emerging-dos|dshield|emerging-exploit|emerging-ftp|emerging-games|emerging-icmp|emerging-icmpinfo|emerging-imap|emerging-inappropriate|emerging-info|emerging-malware|emerging-misc|emerging-mobile|emerging-netbios|emerging-p2p|emerging-policy|emerging-pop3|emerging-rpc|emerging-scada|emerging-scan|emerging-shellcode|emerging-smtp|emerging-snmp|emerging-sql|emerging-telnet|emerging-tftp|tor|emerging-useragent|emerging-voip|emerging-webapps|emerging-webclient|emerging-webserver|emerging-worm|exploit-kit|adware-pup|botcc-portgrouped|phishing|threatview-cs-c2|3coresec|chat|coinminer|current-events|drop|hunting|icmp-info|inappropriate|info|ja3|policy|scada|dark-web-blocker-list|malicious-hosts
	EnabledNetworks                     []string              `json:"enabled_networks,omitempty"`
	Honeypot                            []SettingIpsHoneypot  `json:"honeypot,omitempty"`
	HoneypotEnabled                     bool                  `json:"honeypot_enabled"`
	IPsMode                             string                `json:"ips_mode,omitempty"` // ids|ips|ipsInline|disabled
	MemoryOptimized                     bool                  `json:"memory_optimized"`
	RestrictTorrents                    bool                  `json:"restrict_torrents"`
	Suppression                         SettingIpsSuppression `json:"suppression,omitempty"`
}

func (dst *SettingIps) UnmarshalJSON(b []byte) error {
	type Alias SettingIps
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

type SettingIpsAlerts struct {
	Category  string               `json:"category,omitempty"`
	Gid       int                  `json:"gid,omitempty"`
	ID        int                  `json:"id,omitempty"`
	Signature string               `json:"signature,omitempty"`
	Tracking  []SettingIpsTracking `json:"tracking,omitempty"`
	Type      string               `json:"type,omitempty"` // all|track
}

func (dst *SettingIpsAlerts) UnmarshalJSON(b []byte) error {
	type Alias SettingIpsAlerts
	aux := &struct {
		Gid emptyStringInt `json:"gid"`
		ID  emptyStringInt `json:"id"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.Gid = int(aux.Gid)
	dst.ID = int(aux.ID)

	return nil
}

type SettingIpsHoneypot struct {
	IPAddress string `json:"ip_address,omitempty"`
	NetworkID string `json:"network_id,omitempty"`
	Version   string `json:"version,omitempty"` // v4|v6
}

func (dst *SettingIpsHoneypot) UnmarshalJSON(b []byte) error {
	type Alias SettingIpsHoneypot
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

type SettingIpsSuppression struct {
	Alerts    []SettingIpsAlerts    `json:"alerts,omitempty"`
	Whitelist []SettingIpsWhitelist `json:"whitelist,omitempty"`
}

func (dst *SettingIpsSuppression) UnmarshalJSON(b []byte) error {
	type Alias SettingIpsSuppression
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

type SettingIpsTracking struct {
	Direction string `json:"direction,omitempty"` // both|src|dest
	Mode      string `json:"mode,omitempty"`      // ip|subnet|network
	Value     string `json:"value,omitempty"`
}

func (dst *SettingIpsTracking) UnmarshalJSON(b []byte) error {
	type Alias SettingIpsTracking
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

type SettingIpsWhitelist struct {
	Direction string `json:"direction,omitempty"` // both|src|dest
	Mode      string `json:"mode,omitempty"`      // ip|subnet|network
	Value     string `json:"value,omitempty"`
}

func (dst *SettingIpsWhitelist) UnmarshalJSON(b []byte) error {
	type Alias SettingIpsWhitelist
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

func (c *Client) getSettingIps(ctx context.Context, site string) (*SettingIps, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []SettingIps `json:"data"`
	}
	err := c.do(ctx, "GET", fmt.Sprintf("api/s/%s/get/setting/ips", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingIps(ctx context.Context, site string, d *SettingIps) (*SettingIps, error) {
	var respBody struct {
		Meta meta         `json:"meta"`
		Data []SettingIps `json:"data"`
	}

	d.Key = "ips"
	err := c.do(ctx, "PUT", fmt.Sprintf("api/s/%s/set/setting/ips", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	res := respBody.Data[0]

	return &res, nil
}
