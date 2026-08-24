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

type DeviceSupervision struct {
	BaseSetting

	GlobalSupervisionEnabled bool   `json:"global_supervision_enabled"`
	HeartbeatIntervalSeconds *int64 `json:"heartbeat_interval_seconds,omitempty"` // ^([6-9][0-9]|[1-2][0-9][0-9]|300)$
	PowerOffDurationSeconds  *int64 `json:"power_off_duration_seconds,omitempty"` // ^([6-9][0-9]|[1-9][0-9]{2}|[1-8][0-9]{3}|9000)$
	SilenceThresholdSeconds  *int64 `json:"silence_threshold_seconds,omitempty"`  // ^(300|[3-9][0-9][0-9]|[1-8][0-9]{3}|9000)$
}

func (dst *DeviceSupervision) UnmarshalJSON(b []byte) error {
	type Alias DeviceSupervision
	aux := &struct {
		HeartbeatIntervalSeconds *types.Number `json:"heartbeat_interval_seconds"`
		PowerOffDurationSeconds  *types.Number `json:"power_off_duration_seconds"`
		SilenceThresholdSeconds  *types.Number `json:"silence_threshold_seconds"`

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
	if aux.HeartbeatIntervalSeconds != nil {
		if val, err := aux.HeartbeatIntervalSeconds.Int64(); err == nil {
			dst.HeartbeatIntervalSeconds = &val
		} else if string(*aux.HeartbeatIntervalSeconds) == "" {
			var zero int64
			dst.HeartbeatIntervalSeconds = &zero
		}
	}
	if aux.PowerOffDurationSeconds != nil {
		if val, err := aux.PowerOffDurationSeconds.Int64(); err == nil {
			dst.PowerOffDurationSeconds = &val
		} else if string(*aux.PowerOffDurationSeconds) == "" {
			var zero int64
			dst.PowerOffDurationSeconds = &zero
		}
	}
	if aux.SilenceThresholdSeconds != nil {
		if val, err := aux.SilenceThresholdSeconds.Int64(); err == nil {
			dst.SilenceThresholdSeconds = &val
		} else if string(*aux.SilenceThresholdSeconds) == "" {
			var zero int64
			dst.SilenceThresholdSeconds = &zero
		}
	}

	return nil
}
