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

type UsgGeo struct {
	BaseSetting

	IPFiltering *SettingUsgGeoIPFiltering `json:"ip_filtering,omitempty"`
}

func (dst *UsgGeo) UnmarshalJSON(b []byte) error {
	type Alias UsgGeo
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

type SettingUsgGeoIPFiltering struct {
	Action           string `json:"action,omitempty"`    // block|allow
	Countries        string `json:"countries,omitempty"` // ^([A-Z]{2})?(,[A-Z]{2}){0,149}$
	Enabled          bool   `json:"enabled"`
	TrafficDirection string `json:"traffic_direction,omitempty"` // ^(both|ingress|egress)$
}

func (dst *SettingUsgGeoIPFiltering) UnmarshalJSON(b []byte) error {
	type Alias SettingUsgGeoIPFiltering
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
