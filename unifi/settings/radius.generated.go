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

type Radius struct {
	BaseSetting

	AccountingEnabled     bool   `json:"accounting_enabled"`
	AcctPort              int    `json:"acct_port,omitempty"` // [1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5]
	AuthPort              int    `json:"auth_port,omitempty"` // [1-9][0-9]{0,3}|[1-5][0-9]{4}|[6][0-4][0-9]{3}|[6][5][0-4][0-9]{2}|[6][5][5][0-2][0-9]|[6][5][5][3][0-5]
	ConfigureWholeNetwork bool   `json:"configure_whole_network"`
	Enabled               bool   `json:"enabled"`
	InterimUpdateInterval int    `json:"interim_update_interval,omitempty"` // ^([6-9][0-9]|[1-9][0-9]{2,3}|[1-7][0-9]{4}|8[0-5][0-9]{3}|86[0-3][0-9][0-9]|86400)$
	TunneledReply         bool   `json:"tunneled_reply"`
	XSecret               string `json:"x_secret,omitempty"` // ^[^\\"' ]{1,48}$
}

func (dst *Radius) UnmarshalJSON(b []byte) error {
	type Alias Radius
	aux := &struct {
		AcctPort              types.Number `json:"acct_port"`
		AuthPort              types.Number `json:"auth_port"`
		InterimUpdateInterval types.Number `json:"interim_update_interval"`

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
	if val, err := aux.AcctPort.Int64(); err == nil {
		dst.AcctPort = int(val)
	}
	if val, err := aux.AuthPort.Int64(); err == nil {
		dst.AuthPort = int(val)
	}
	if val, err := aux.InterimUpdateInterval.Int64(); err == nil {
		dst.InterimUpdateInterval = int(val)
	}

	return nil
}
