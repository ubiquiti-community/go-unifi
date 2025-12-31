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
)

type Netflow struct {
	BaseSetting

	AutoEngineIDEnabled bool     `json:"auto_engine_id_enabled"`
	Enabled             bool     `json:"enabled"`
	EngineID            int      `json:"engine_id,omitempty"` // ^$|[1-9][0-9]*
	ExportFrequency     int      `json:"export_frequency,omitempty"`
	NetworkIDs          []string `json:"network_ids,omitempty"`
	Port                int      `json:"port,omitempty"` // 102[4-9]|10[3-9][0-9]|1[1-9][0-9]{2}|[2-9][0-9]{3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5]
	RefreshRate         int      `json:"refresh_rate,omitempty"`
	SamplingMode        string   `json:"sampling_mode,omitempty"` // off|hash|random|deterministic
	SamplingRate        int      `json:"sampling_rate,omitempty"` // [2-9]|[1-9][0-9]{1,3}|1[0-5][0-9]{3}|16[0-2][0-9]{2}|163[0-7][0-9]|1638[0-3]|^$
	Server              string   `json:"server,omitempty"`        // .{0,252}[^\.]$
	Version             int      `json:"version,omitempty"`       // 5|9|10
}

func (dst *Netflow) UnmarshalJSON(b []byte) error {
	type Alias Netflow
	aux := &struct {
		EngineID        types.Number `json:"engine_id"`
		ExportFrequency types.Number `json:"export_frequency"`
		Port            types.Number `json:"port"`
		RefreshRate     types.Number `json:"refresh_rate"`
		SamplingRate    types.Number `json:"sampling_rate"`
		Version         types.Number `json:"version"`

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
	if val, err := aux.EngineID.Int64(); err == nil {
		dst.EngineID = int(val)
	}
	if val, err := aux.ExportFrequency.Int64(); err == nil {
		dst.ExportFrequency = int(val)
	}
	if val, err := aux.Port.Int64(); err == nil {
		dst.Port = int(val)
	}
	if val, err := aux.RefreshRate.Int64(); err == nil {
		dst.RefreshRate = int(val)
	}
	if val, err := aux.SamplingRate.Int64(); err == nil {
		dst.SamplingRate = int(val)
	}
	if val, err := aux.Version.Int64(); err == nil {
		dst.Version = int(val)
	}

	return nil
}
