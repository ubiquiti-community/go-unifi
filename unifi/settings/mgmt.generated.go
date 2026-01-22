// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
	_ types.Number
	_ strconv.NumError
)

type Mgmt struct {
	BaseSetting

	AdvancedFeatureEnabled  bool                  `json:"advanced_feature_enabled"`
	AlertEnabled            bool                  `json:"alert_enabled"`
	AutoUpgrade             bool                  `json:"auto_upgrade"`
	AutoUpgradeHour         *int64                `json:"auto_upgrade_hour,omitempty"` // [0-9]|1[0-9]|2[0-3]|^$
	BootSound               bool                  `json:"boot_sound"`
	DebugToolsEnabled       bool                  `json:"debug_tools_enabled"`
	DirectConnectEnabled    bool                  `json:"direct_connect_enabled"`
	LedEnabled              bool                  `json:"led_enabled"`
	OutdoorModeEnabled      bool                  `json:"outdoor_mode_enabled"`
	UnifiIDpEnabled         bool                  `json:"unifi_idp_enabled"`
	WifimanEnabled          bool                  `json:"wifiman_enabled"`
	XMgmtKey                string                `json:"x_mgmt_key,omitempty"` // [0-9a-f]{32}
	XSshAuthPasswordEnabled bool                  `json:"x_ssh_auth_password_enabled"`
	XSshBindWildcard        bool                  `json:"x_ssh_bind_wildcard"`
	XSshEnabled             bool                  `json:"x_ssh_enabled"`
	XSshKeys                []SettingMgmtXSshKeys `json:"x_ssh_keys,omitempty"`
	XSshMd5Passwd           string                `json:"x_ssh_md5passwd,omitempty"`
	XSshPassword            string                `json:"x_ssh_password,omitempty"` // .{1,128}
	XSshSha512Passwd        string                `json:"x_ssh_sha512passwd,omitempty"`
	XSshUsername            string                `json:"x_ssh_username,omitempty"` // ^[_A-Za-z0-9][-_.A-Za-z0-9]{0,29}$
}

func (dst *Mgmt) UnmarshalJSON(b []byte) error {
	type Alias Mgmt
	aux := &struct {
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

	return nil
}

type SettingMgmtXSshKeys struct {
	Comment     string `json:"comment"`
	Date        string `json:"date"`
	Fingerprint string `json:"fingerprint"`
	Key         string `json:"key"`
	KeyType     string `json:"type"`
	Name        string `json:"name"`
}

func (dst *SettingMgmtXSshKeys) UnmarshalJSON(b []byte) error {
	type Alias SettingMgmtXSshKeys
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
