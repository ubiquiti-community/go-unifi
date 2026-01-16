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

type RoamingAssistant struct {
	BaseSetting

	Enabled bool  `json:"enabled"`
	Rssi    int64 `json:"rssi,omitempty"` // ^-([6-7][0-9]|80)$
}

func (dst *RoamingAssistant) UnmarshalJSON(b []byte) error {
	type Alias RoamingAssistant
	aux := &struct {
		Rssi types.Number `json:"rssi"`

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
	if val, err := aux.Rssi.Int64(); err == nil {
		dst.Rssi = int64(val)
	}

	return nil
}
