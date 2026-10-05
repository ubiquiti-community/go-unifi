package unifi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNatPortUnmarshal(t *testing.T) {
	cases := []struct {
		name string
		port string
		want string
	}{
		{name: "numeric", port: `8080`, want: "8080"},
		{name: "string", port: `"8080"`, want: "8080"},
		{name: "range", port: `"80-90"`, want: "80-90"},
		{name: "empty string", port: `""`, want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := `{"port":` + tc.port +
				`,"source_filter":{"port":` + tc.port + `}` +
				`,"destination_filter":{"port":` + tc.port + `}}`

			var got Nat
			if err := json.Unmarshal([]byte(body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.Port != tc.want {
				t.Errorf("Port = %q, want %q", got.Port, tc.want)
			}
			if got.SourceFilter.Port != tc.want {
				t.Errorf("SourceFilter.Port = %q, want %q", got.SourceFilter.Port, tc.want)
			}
			if got.DestinationFilter.Port != tc.want {
				t.Errorf("DestinationFilter.Port = %q, want %q", got.DestinationFilter.Port, tc.want)
			}
		})
	}
}

func TestNatPortAbsent(t *testing.T) {
	var got Nat
	if err := json.Unmarshal([]byte(`{"source_filter":{},"destination_filter":{}}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Port != "" || got.SourceFilter.Port != "" || got.DestinationFilter.Port != "" {
		t.Fatalf("expected empty ports, got %q, %q, %q",
			got.Port, got.SourceFilter.Port, got.DestinationFilter.Port)
	}
}

func TestNatPortRangeRoundTrip(t *testing.T) {
	body := `{"_id":"nat1","type":"DNAT","port":"8000-8010",` +
		`"source_filter":{"filter_type":"ADDRESS_AND_PORT","port":"80-90"},` +
		`"destination_filter":{"filter_type":"ADDRESS_AND_PORT","port":"8000-8010"}}`

	var nat Nat
	if err := json.Unmarshal([]byte(body), &nat); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	b, err := json.Marshal(&nat)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(b)
	for _, want := range []string{
		`"port":"8000-8010","pppoe_use_base_interface"`,
		`"source_filter":{"filter_type":"ADDRESS_AND_PORT","invert_address":false,"invert_port":false,"port":"80-90"}`,
		`"destination_filter":{"filter_type":"ADDRESS_AND_PORT","invert_address":false,"invert_port":false,"port":"8000-8010"}`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %s in payload, got: %s", want, got)
		}
	}
}

func TestNatPortOmittedWhenEmpty(t *testing.T) {
	b, err := json.Marshal(&Nat{
		Type:              "DNAT",
		SourceFilter:      &NatSourceFilter{FilterType: "NONE"},
		DestinationFilter: &NatDestinationFilter{FilterType: "NONE"},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), `"port"`) {
		t.Fatalf("expected port to be absent, got: %s", b)
	}
}
