package unifi

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFirewallPolicyScheduleRoundTripsControllerFields(t *testing.T) {
	raw := []byte(`{
		"date":"2026-07-10",
		"date_start":"2026-07-01",
		"date_end":"2026-07-31",
		"mode":"CUSTOM",
		"repeat_on_days":["mon","wed","fri"],
		"time_all_day":false,
		"time_range_start":"09:00",
		"time_range_end":"17:30"
	}`)

	var got FirewallPolicySchedule
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal schedule: %v", err)
	}

	if got.Date != "2026-07-10" {
		t.Fatalf("Date = %q, want %q", got.Date, "2026-07-10")
	}
	if got.DateStart != "2026-07-01" {
		t.Fatalf("DateStart = %q, want %q", got.DateStart, "2026-07-01")
	}
	if got.DateEnd != "2026-07-31" {
		t.Fatalf("DateEnd = %q, want %q", got.DateEnd, "2026-07-31")
	}
	if got.Mode != "CUSTOM" {
		t.Fatalf("Mode = %q, want %q", got.Mode, "CUSTOM")
	}
	if !reflect.DeepEqual(got.RepeatOnDays, []string{"mon", "wed", "fri"}) {
		t.Fatalf("RepeatOnDays = %#v", got.RepeatOnDays)
	}
	if got.TimeAllDay == nil || *got.TimeAllDay {
		t.Fatalf("TimeAllDay = %v, want pointer to false", got.TimeAllDay)
	}
	if got.TimeRangeStart != "09:00" || got.TimeRangeEnd != "17:30" {
		t.Fatalf("time range = %q-%q", got.TimeRangeStart, got.TimeRangeEnd)
	}

	roundTripped, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal schedule: %v", err)
	}

	var wantObject, gotObject map[string]any
	if err := json.Unmarshal(raw, &wantObject); err != nil {
		t.Fatalf("decode expected JSON: %v", err)
	}
	if err := json.Unmarshal(roundTripped, &gotObject); err != nil {
		t.Fatalf("decode round-tripped JSON: %v", err)
	}
	if !reflect.DeepEqual(gotObject, wantObject) {
		t.Fatalf("round-trip mismatch:\n got: %s\nwant: %s", roundTripped, raw)
	}
}

func TestFirewallPolicySchedulePreservesLegacyFieldsUnderAlways(t *testing.T) {
	raw := []byte(`{
		"date_start":"2025-06-20",
		"date_end":"2025-06-27",
		"mode":"ALWAYS",
		"repeat_on_days":[],
		"time_all_day":false,
		"time_range_start":"09:00",
		"time_range_end":"12:00"
	}`)

	var schedule FirewallPolicySchedule
	if err := json.Unmarshal(raw, &schedule); err != nil {
		t.Fatalf("unmarshal schedule: %v", err)
	}
	roundTripped, err := json.Marshal(schedule)
	if err != nil {
		t.Fatalf("marshal schedule: %v", err)
	}

	var got FirewallPolicySchedule
	if err := json.Unmarshal(roundTripped, &got); err != nil {
		t.Fatalf("decode round-tripped schedule: %v", err)
	}
	if got.Mode != "ALWAYS" ||
		got.DateStart != "2025-06-20" || got.DateEnd != "2025-06-27" ||
		got.TimeAllDay == nil || *got.TimeAllDay ||
		got.TimeRangeStart != "09:00" || got.TimeRangeEnd != "12:00" ||
		len(got.RepeatOnDays) != 0 {
		t.Fatalf("legacy schedule semantics changed: %#v", got)
	}
}

func TestFirewallPolicySchedulePreservesAbsentTimeAllDay(t *testing.T) {
	raw := []byte(`{"mode":"ALWAYS"}`)

	var schedule FirewallPolicySchedule
	if err := json.Unmarshal(raw, &schedule); err != nil {
		t.Fatalf("unmarshal schedule: %v", err)
	}
	if schedule.TimeAllDay != nil {
		t.Fatalf("TimeAllDay = %v, want nil", schedule.TimeAllDay)
	}

	roundTripped, err := json.Marshal(schedule)
	if err != nil {
		t.Fatalf("marshal schedule: %v", err)
	}
	if string(roundTripped) != string(raw) {
		t.Fatalf("round-trip = %s, want %s", roundTripped, raw)
	}
}
