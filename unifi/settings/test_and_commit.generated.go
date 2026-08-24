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

type TestAndCommit struct {
	BaseSetting

	DeviceVLANTesting bool                       `json:"device_vlan_testing"`
	Enabled           bool                       `json:"enabled"`
	Mode              string                     `json:"mode,omitempty"` // auto|custom
	Pages             *SettingTestAndCommitPages `json:"pages,omitempty"`
}

func (dst *TestAndCommit) UnmarshalJSON(b []byte) error {
	type Alias TestAndCommit
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

type SettingTestAndCommitPages struct {
	Internet        bool `json:"internet"`
	NetworkSettings bool `json:"network_settings"`
	VPN             bool `json:"vpn"`
}

func (dst *SettingTestAndCommitPages) UnmarshalJSON(b []byte) error {
	type Alias SettingTestAndCommitPages
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
