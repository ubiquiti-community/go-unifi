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

type RadioAi struct {
	BaseSetting

	AutoAdjustChannelsToCountry bool                                `json:"auto_adjust_channels_to_country"`
	AutoChannelPresetsType      string                              `json:"auto_channel_presets_type,omitempty"` // maximum_speed|conservative|custom
	Channels6E                  []int                               `json:"channels_6e,omitempty"`               // [1-9]|[1-2][0-9]|3[3-9]|[4-5][0-9]|6[0-1]|6[5-9]|[7-8][0-9]|9[0-3]|9[7-9]|1[0-1][0-9]|12[0-5]|129|1[3-4][0-9]|15[0-7]|16[1-9]|1[7-8][0-9]|19[3-9]|2[0-1][0-9]|22[0-1]|22[5-9]|233
	ChannelsBlacklist           []SettingRadioAiChannelsBlacklist   `json:"channels_blacklist,omitempty"`
	ChannelsNa                  []int                               `json:"channels_na,omitempty"` // 34|36|38|40|42|44|46|48|52|56|60|64|100|104|108|112|116|120|124|128|132|136|140|144|149|153|157|161|165|169
	ChannelsNg                  []int                               `json:"channels_ng,omitempty"` // 1|2|3|4|5|6|7|8|9|10|11|12|13|14
	CronExpr                    string                              `json:"cron_expr,omitempty"`
	Default                     bool                                `json:"default"`
	Enabled                     bool                                `json:"enabled"`
	ExcludeDevices              []string                            `json:"exclude_devices,omitempty"`       // ([0-9a-z]{2}:){5}[0-9a-z]{2}
	HighPriorityDevices         []string                            `json:"high_priority_devices,omitempty"` // ([0-9a-z]{2}:){5}[0-9a-z]{2}
	HtModesNa                   []int                               `json:"ht_modes_na,omitempty"`           // ^(20|40|80|160)$
	HtModesNg                   []int                               `json:"ht_modes_ng,omitempty"`           // ^(20|40)$
	Optimize                    []string                            `json:"optimize,omitempty"`              // channel|power
	Radios                      []string                            `json:"radios,omitempty"`                // na|ng|6e
	RadiosConfiguration         []SettingRadioAiRadiosConfiguration `json:"radios_configuration,omitempty"`
	SettingPreference           string                              `json:"setting_preference,omitempty"` // auto|manual
	UseXy                       bool                                `json:"useXY"`
}

func (dst *RadioAi) UnmarshalJSON(b []byte) error {
	type Alias RadioAi
	aux := &struct {
		Channels6E []types.Number `json:"channels_6e"`
		ChannelsNa []types.Number `json:"channels_na"`
		ChannelsNg []types.Number `json:"channels_ng"`
		HtModesNa  []types.Number `json:"ht_modes_na"`
		HtModesNg  []types.Number `json:"ht_modes_ng"`

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
	dst.Channels6E = make([]int, len(aux.Channels6E))
	for i, v := range aux.Channels6E {
		if val, err := v.Int64(); err == nil {
			dst.Channels6E[i] = int(val)
		}
	}
	dst.ChannelsNa = make([]int, len(aux.ChannelsNa))
	for i, v := range aux.ChannelsNa {
		if val, err := v.Int64(); err == nil {
			dst.ChannelsNa[i] = int(val)
		}
	}
	dst.ChannelsNg = make([]int, len(aux.ChannelsNg))
	for i, v := range aux.ChannelsNg {
		if val, err := v.Int64(); err == nil {
			dst.ChannelsNg[i] = int(val)
		}
	}
	dst.HtModesNa = make([]int, len(aux.HtModesNa))
	for i, v := range aux.HtModesNa {
		if val, err := v.Int64(); err == nil {
			dst.HtModesNa[i] = int(val)
		}
	}
	dst.HtModesNg = make([]int, len(aux.HtModesNg))
	for i, v := range aux.HtModesNg {
		if val, err := v.Int64(); err == nil {
			dst.HtModesNg[i] = int(val)
		}
	}

	return nil
}

type SettingRadioAiChannelsBlacklist struct {
	Channel      int    `json:"channel,omitempty"`       // [1-9]|[1-9][0-9]|1[0-9][0-9]|2[0-9]|2[0-1][0-9]|22[0-1]|22[5-9]|233
	ChannelWidth int    `json:"channel_width,omitempty"` // 20|40|80|160|240|320
	Radio        string `json:"radio,omitempty"`         // na|ng|6e
}

func (dst *SettingRadioAiChannelsBlacklist) UnmarshalJSON(b []byte) error {
	type Alias SettingRadioAiChannelsBlacklist
	aux := &struct {
		Channel      types.Number `json:"channel"`
		ChannelWidth types.Number `json:"channel_width"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if val, err := aux.Channel.Int64(); err == nil {
		dst.Channel = int(val)
	}
	if val, err := aux.ChannelWidth.Int64(); err == nil {
		dst.ChannelWidth = int(val)
	}

	return nil
}

type SettingRadioAiRadiosConfiguration struct {
	ChannelWidth int    `json:"channel_width,omitempty"` // 20|40|80|160|320
	Dfs          bool   `json:"dfs"`
	Radio        string `json:"radio,omitempty"` // na|ng|6e
}

func (dst *SettingRadioAiRadiosConfiguration) UnmarshalJSON(b []byte) error {
	type Alias SettingRadioAiRadiosConfiguration
	aux := &struct {
		ChannelWidth types.Number `json:"channel_width"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if val, err := aux.ChannelWidth.Int64(); err == nil {
		dst.ChannelWidth = int(val)
	}

	return nil
}
