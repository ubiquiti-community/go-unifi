// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ubiquiti-community/go-unifi/unifi/types"
)

// just to fix compile issues with the import.
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
	_ types.Number
	_ strconv.NumError
	_ strings.Builder
)

type QOSRule struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Destination       *QOSRuleDestination `json:"destination,omitempty"`
	DownloadBurst     string              `json:"download_burst,omitempty"`      // OFF|SHORT|LONG
	DownloadLimitKbps *int64              `json:"download_limit_kbps,omitempty"` // [1-9][0-9]*
	Enabled           bool                `json:"enabled"`
	Index             *int64              `json:"index,omitempty"`     // [1-9][0-9]+
	Name              string              `json:"name,omitempty"`      // .{1,128}
	Objective         string              `json:"objective,omitempty"` // LIMIT|PRIORITIZE|LIMIT_AND_PRIORITIZE
	Schedule          *QOSRuleSchedule    `json:"schedule,omitempty"`
	Source            *QOSRuleSource      `json:"source,omitempty"`
	UploadBurst       string              `json:"upload_burst,omitempty"`      // OFF|SHORT|LONG
	UploadLimitKbps   *int64              `json:"upload_limit_kbps,omitempty"` // [1-9][0-9]*
	WANOrVPNNetwork   string              `json:"wan_or_vpn_network,omitempty"`
}

func (dst *QOSRule) UnmarshalJSON(b []byte) error {
	type Alias QOSRule
	aux := &struct {
		DownloadLimitKbps *types.Number `json:"download_limit_kbps"`
		Index             *types.Number `json:"index"`
		UploadLimitKbps   *types.Number `json:"upload_limit_kbps"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	if aux.DownloadLimitKbps != nil {
		if val, err := aux.DownloadLimitKbps.Int64(); err == nil {
			dst.DownloadLimitKbps = &val
		} else if string(*aux.DownloadLimitKbps) == "" {
			var zero int64
			dst.DownloadLimitKbps = &zero
		}
	}
	if aux.Index != nil {
		if val, err := aux.Index.Int64(); err == nil {
			dst.Index = &val
		} else if string(*aux.Index) == "" {
			var zero int64
			dst.Index = &zero
		}
	}
	if aux.UploadLimitKbps != nil {
		if val, err := aux.UploadLimitKbps.Int64(); err == nil {
			dst.UploadLimitKbps = &val
		} else if string(*aux.UploadLimitKbps) == "" {
			var zero int64
			dst.UploadLimitKbps = &zero
		}
	}

	return nil
}

type QOSRuleDestination struct {
	AppCategoryIDs     []int64  `json:"app_category_ids,omitempty"`
	AppIDs             []int64  `json:"app_ids,omitempty"`
	IPGroupID          string   `json:"ip_group_id,omitempty"`
	IPs                []string `json:"ips,omitempty"`
	Iid                string   `json:"iid,omitempty"`
	MatchingTarget     string   `json:"matching_target,omitempty"`      // ANY|APP|APP_CATEGORY|IP|IID|REGION|WEB
	MatchingTargetType string   `json:"matching_target_type,omitempty"` // SPECIFIC|OBJECT
	Port               string   `json:"port,omitempty"`
	PortGroupID        string   `json:"port_group_id,omitempty"`
	PortMatchingType   string   `json:"port_matching_type,omitempty"` // ANY|SPECIFIC|OBJECT
	Regions            []string `json:"regions,omitempty"`
	WebDomains         []string `json:"web_domains,omitempty"`
	WebGroupID         string   `json:"web_group_id,omitempty"`
}

func (dst *QOSRuleDestination) UnmarshalJSON(b []byte) error {
	type Alias QOSRuleDestination
	aux := &struct {
		AppCategoryIDs []types.Number `json:"app_category_ids"`
		AppIDs         []types.Number `json:"app_ids"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.AppCategoryIDs = make([]int64, len(aux.AppCategoryIDs))
	for i, v := range aux.AppCategoryIDs {
		if val, err := v.Int64(); err == nil {
			dst.AppCategoryIDs[i] = val
		}
	}
	dst.AppIDs = make([]int64, len(aux.AppIDs))
	for i, v := range aux.AppIDs {
		if val, err := v.Int64(); err == nil {
			dst.AppIDs[i] = val
		}
	}

	return nil
}

type QOSRuleSchedule struct {
	Date           string   `json:"date,omitempty"`
	DateEnd        string   `json:"date_end,omitempty"`
	DateStart      string   `json:"date_start,omitempty"`
	Mode           string   `json:"mode,omitempty"`           // ALWAYS|EVERY_DAY|EVERY_WEEK|ONE_TIME_ONLY|CUSTOM
	RepeatOnDays   []string `json:"repeat_on_days,omitempty"` // mon|tue|wed|thu|fri|sat|sun
	TimeAllDay     bool     `json:"time_all_day"`
	TimeRangeEnd   string   `json:"time_range_end,omitempty"`
	TimeRangeStart string   `json:"time_range_start,omitempty"`
}

func (dst *QOSRuleSchedule) UnmarshalJSON(b []byte) error {
	type Alias QOSRuleSchedule
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

type QOSRuleSource struct {
	ClientMACs       []string `json:"client_macs,omitempty"`     // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$
	MatchingTarget   string   `json:"matching_target,omitempty"` // ANY|CLIENT|NETWORK
	NetworkIDs       []string `json:"network_ids,omitempty"`
	PortMatchingType string   `json:"port_matching_type,omitempty"` // ANY
}

func (dst *QOSRuleSource) UnmarshalJSON(b []byte) error {
	type Alias QOSRuleSource
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

func (c *ApiClient) listQOSRule(
	ctx context.Context,
	site string,
	query ...map[string]string,
) ([]QOSRule, error) {
	var respBody []QOSRule

	err := c.do(
		ctx,
		http.MethodGet,
		fmt.Sprintf("v2/api/site/%s/qos-rules", site),
		nil,
		&respBody,
		query...,
	)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func (c *ApiClient) getQOSRule(
	ctx context.Context,
	site string,
	id string,
) (*QOSRule, error) {
	respBody, err := c.listQOSRule(ctx, site)
	if err != nil {
		return nil, err
	}

	if len(respBody) == 0 {
		return nil, &NotFoundError{}
	}

	for _, val := range respBody {
		if val.ID == id {
			return &val, nil
		}
	}

	return nil, &NotFoundError{}
}

func (c *ApiClient) deleteQOSRule(
	ctx context.Context,
	site string,
	id string,
) error {
	err := c.do(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("v2/api/site/%s/qos-rules/%s", site, id),
		struct{}{},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) createQOSRule(
	ctx context.Context,
	site string,
	d *QOSRule,
) (*QOSRule, error) {
	var respBody QOSRule

	err := c.do(
		ctx,
		http.MethodPost,
		fmt.Sprintf("v2/api/site/%s/qos-rules", site),
		d,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}

func (c *ApiClient) updateQOSRule(
	ctx context.Context,
	site string,
	d *QOSRule,
) (*QOSRule, error) {
	var respBody QOSRule
	err := c.do(
		ctx,
		http.MethodPut,
		fmt.Sprintf("v2/api/site/%s/qos-rules/%s", site, d.ID),
		d,
		&respBody,
	)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}
