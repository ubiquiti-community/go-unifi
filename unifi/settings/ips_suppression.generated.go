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

type IpsSuppression struct {
	BaseSetting

	Alerts    []SettingIpsSuppressionAlerts    `json:"alerts,omitempty"`
	Whitelist []SettingIpsSuppressionWhitelist `json:"whitelist,omitempty"`
}

func (dst *IpsSuppression) UnmarshalJSON(b []byte) error {
	type Alias IpsSuppression
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

type SettingIpsSuppressionAlerts struct {
	Category  string                          `json:"category,omitempty"`
	Gid       *int64                          `json:"gid,omitempty"`
	ID        *int64                          `json:"id,omitempty"`
	Signature string                          `json:"signature,omitempty"`
	Tracking  []SettingIpsSuppressionTracking `json:"tracking,omitempty"`
	Type      string                          `json:"type,omitempty"` // all|track
}

func (dst *SettingIpsSuppressionAlerts) UnmarshalJSON(b []byte) error {
	type Alias SettingIpsSuppressionAlerts
	aux := &struct {
		Gid *types.Number `json:"gid"`
		ID  *types.Number `json:"id"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if aux.Gid != nil {
		if val, err := aux.Gid.Int64(); err == nil {
			dst.Gid = &val
		} else if string(*aux.Gid) == "" {
			var zero int64
			dst.Gid = &zero
		}
	}
	if aux.ID != nil {
		if val, err := aux.ID.Int64(); err == nil {
			dst.ID = &val
		} else if string(*aux.ID) == "" {
			var zero int64
			dst.ID = &zero
		}
	}

	return nil
}

type SettingIpsSuppressionTracking struct {
	Direction string `json:"direction,omitempty"` // both|src|dest
	Mode      string `json:"mode,omitempty"`      // ip|subnet|network
	Value     string `json:"value,omitempty"`
}

func (dst *SettingIpsSuppressionTracking) UnmarshalJSON(b []byte) error {
	type Alias SettingIpsSuppressionTracking
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

type SettingIpsSuppressionWhitelist struct {
	Direction string `json:"direction,omitempty"` // both|src|dest
	Mode      string `json:"mode,omitempty"`      // ip|subnet|network
	Value     string `json:"value,omitempty"`
}

func (dst *SettingIpsSuppressionWhitelist) UnmarshalJSON(b []byte) error {
	type Alias SettingIpsSuppressionWhitelist
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
