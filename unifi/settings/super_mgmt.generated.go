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

type SuperMgmt struct {
	BaseSetting

	AnalyticsDisapprovedFor                  string   `json:"analytics_disapproved_for,omitempty"`
	AutoUpgrade                              bool     `json:"auto_upgrade"`
	AutobackupCronExpr                       string   `json:"autobackup_cron_expr,omitempty"`
	AutobackupDays                           int      `json:"autobackup_days,omitempty"`
	AutobackupEnabled                        bool     `json:"autobackup_enabled"`
	AutobackupGcsBucket                      string   `json:"autobackup_gcs_bucket,omitempty"`
	AutobackupGcsCertificatePath             string   `json:"autobackup_gcs_certificate_path,omitempty"`
	AutobackupLocalPath                      string   `json:"autobackup_local_path,omitempty"`
	AutobackupMaxFiles                       int      `json:"autobackup_max_files,omitempty"`
	AutobackupPostActions                    []string `json:"autobackup_post_actions,omitempty"` // copy_local|copy_gcs|copy_cloud
	AutobackupTimezone                       string   `json:"autobackup_timezone,omitempty"`
	BackupToCloudEnabled                     bool     `json:"backup_to_cloud_enabled"`
	ContactInfoCity                          string   `json:"contact_info_city,omitempty"`
	ContactInfoCompanyName                   string   `json:"contact_info_company_name,omitempty"`
	ContactInfoCountry                       string   `json:"contact_info_country,omitempty"`
	ContactInfoFullName                      string   `json:"contact_info_full_name,omitempty"`
	ContactInfoPhoneNumber                   string   `json:"contact_info_phone_number,omitempty"`
	ContactInfoShippingAddress1              string   `json:"contact_info_shipping_address_1,omitempty"`
	ContactInfoShippingAddress2              string   `json:"contact_info_shipping_address_2,omitempty"`
	ContactInfoState                         string   `json:"contact_info_state,omitempty"`
	ContactInfoZip                           string   `json:"contact_info_zip,omitempty"`
	DataRetentionSettingPreference           string   `json:"data_retention_setting_preference,omitempty"` // auto|manual
	DataRetentionTimeInHoursFor5MinutesScale int      `json:"data_retention_time_in_hours_for_5minutes_scale,omitempty"`
	DataRetentionTimeInHoursForDailyScale    int      `json:"data_retention_time_in_hours_for_daily_scale,omitempty"`
	DataRetentionTimeInHoursForHourlyScale   int      `json:"data_retention_time_in_hours_for_hourly_scale,omitempty"`
	DataRetentionTimeInHoursForMonthlyScale  int      `json:"data_retention_time_in_hours_for_monthly_scale,omitempty"`
	DataRetentionTimeInHoursForOthers        int      `json:"data_retention_time_in_hours_for_others,omitempty"`
	DefaultSiteDeviceAuthPasswordAlert       string   `json:"default_site_device_auth_password_alert,omitempty"` // false
	Discoverable                             bool     `json:"discoverable"`
	EnableAnalytics                          bool     `json:"enable_analytics"`
	GoogleMapsApiKey                         string   `json:"google_maps_api_key,omitempty"`
	ImageMapsUseGoogleEngine                 bool     `json:"image_maps_use_google_engine"`
	LedEnabled                               bool     `json:"led_enabled"`
	LiveChat                                 string   `json:"live_chat,omitempty"`    // disabled|super-only|everyone
	LiveUpdates                              string   `json:"live_updates,omitempty"` // disabled|live|auto
	MinimumUsableHdSpace                     int      `json:"minimum_usable_hd_space,omitempty"`
	MinimumUsableSdSpace                     int      `json:"minimum_usable_sd_space,omitempty"`
	MultipleSitesEnabled                     bool     `json:"multiple_sites_enabled"`
	OverrideInformHost                       bool     `json:"override_inform_host"`
	OverrideInformHostLocation               string   `json:"override_inform_host_location,omitempty"`
	StoreEnabled                             string   `json:"store_enabled,omitempty"` // disabled|super-only|everyone
	TimeSeriesPerClientStatsEnabled          bool     `json:"time_series_per_client_stats_enabled"`
	XSshPassword                             string   `json:"x_ssh_password,omitempty"`
	XSshUsername                             string   `json:"x_ssh_username,omitempty"`
}

func (dst *SuperMgmt) UnmarshalJSON(b []byte) error {
	type Alias SuperMgmt
	aux := &struct {
		AutobackupDays                           types.Number `json:"autobackup_days"`
		AutobackupMaxFiles                       types.Number `json:"autobackup_max_files"`
		DataRetentionTimeInHoursFor5MinutesScale types.Number `json:"data_retention_time_in_hours_for_5minutes_scale"`
		DataRetentionTimeInHoursForDailyScale    types.Number `json:"data_retention_time_in_hours_for_daily_scale"`
		DataRetentionTimeInHoursForHourlyScale   types.Number `json:"data_retention_time_in_hours_for_hourly_scale"`
		DataRetentionTimeInHoursForMonthlyScale  types.Number `json:"data_retention_time_in_hours_for_monthly_scale"`
		DataRetentionTimeInHoursForOthers        types.Number `json:"data_retention_time_in_hours_for_others"`
		MinimumUsableHdSpace                     types.Number `json:"minimum_usable_hd_space"`
		MinimumUsableSdSpace                     types.Number `json:"minimum_usable_sd_space"`

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
	if val, err := aux.AutobackupDays.Int64(); err == nil {
		dst.AutobackupDays = int(val)
	}
	if val, err := aux.AutobackupMaxFiles.Int64(); err == nil {
		dst.AutobackupMaxFiles = int(val)
	}
	if val, err := aux.DataRetentionTimeInHoursFor5MinutesScale.Int64(); err == nil {
		dst.DataRetentionTimeInHoursFor5MinutesScale = int(val)
	}
	if val, err := aux.DataRetentionTimeInHoursForDailyScale.Int64(); err == nil {
		dst.DataRetentionTimeInHoursForDailyScale = int(val)
	}
	if val, err := aux.DataRetentionTimeInHoursForHourlyScale.Int64(); err == nil {
		dst.DataRetentionTimeInHoursForHourlyScale = int(val)
	}
	if val, err := aux.DataRetentionTimeInHoursForMonthlyScale.Int64(); err == nil {
		dst.DataRetentionTimeInHoursForMonthlyScale = int(val)
	}
	if val, err := aux.DataRetentionTimeInHoursForOthers.Int64(); err == nil {
		dst.DataRetentionTimeInHoursForOthers = int(val)
	}
	if val, err := aux.MinimumUsableHdSpace.Int64(); err == nil {
		dst.MinimumUsableHdSpace = int(val)
	}
	if val, err := aux.MinimumUsableSdSpace.Int64(); err == nil {
		dst.MinimumUsableSdSpace = int(val)
	}

	return nil
}
