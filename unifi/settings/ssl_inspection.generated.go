// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package settings

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

type SslInspection struct {
	BaseSetting

	State string `json:"state,omitempty"` // off|simple|advanced
}

func (dst *SslInspection) UnmarshalJSON(b []byte) error {
	type Alias SslInspection
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
