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

type Country struct {
	BaseSetting

	Code int `json:"code,omitempty"`
}

func (dst *Country) UnmarshalJSON(b []byte) error {
	type Alias Country
	aux := &struct {
		Code types.Number `json:"code"`

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
	if val, err := aux.Code.Int64(); err == nil {
		dst.Code = int(val)
	}

	return nil
}
