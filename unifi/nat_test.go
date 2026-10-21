package unifi

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestNatPortDecodesNumbersStringsAndRanges guards the NAT `port` fields: the
// controller stores them as strings ("443") and accepts ranges ("8080-8090"),
// which an *int64 silently dropped on read.
func TestNatPortDecodesNumbersStringsAndRanges(t *testing.T) {
	cases := []struct {
		name string
		json string
		want string
	}{
		{name: "number", json: `443`, want: "443"},
		{name: "string", json: `"443"`, want: "443"},
		{name: "range", json: `"8080-8090"`, want: "8080-8090"},
		{name: "empty", json: `""`, want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var nat Nat
			body := `{"port":` + tc.json +
				`,"source_filter":{"port":` + tc.json + `}` +
				`,"destination_filter":{"port":` + tc.json + `}}`
			if err := json.Unmarshal([]byte(body), &nat); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if nat.Port != tc.want {
				t.Errorf("Port = %q, want %q", nat.Port, tc.want)
			}
			if nat.SourceFilter.Port != tc.want {
				t.Errorf("SourceFilter.Port = %q, want %q", nat.SourceFilter.Port, tc.want)
			}
			if nat.DestinationFilter.Port != tc.want {
				t.Errorf("DestinationFilter.Port = %q, want %q", nat.DestinationFilter.Port, tc.want)
			}
		})
	}

	t.Run("absent", func(t *testing.T) {
		var nat Nat
		if err := json.Unmarshal([]byte(`{"source_filter":{}}`), &nat); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if nat.Port != "" || nat.SourceFilter.Port != "" {
			t.Errorf("ports = %q/%q, want empty", nat.Port, nat.SourceFilter.Port)
		}
	})
}

// TestNatPortMarshalsAsString checks that a range is sent as a string and an
// unset port is omitted rather than sent as "".
func TestNatPortMarshalsAsString(t *testing.T) {
	b, err := json.Marshal(&Nat{
		Port:              "8080-8090",
		SourceFilter:      &NatSourceFilter{FilterType: "NONE"},
		DestinationFilter: &NatDestinationFilter{FilterType: "ADDRESS_AND_PORT", Port: "8080-8090"},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got struct {
		Port              json.RawMessage            `json:"port"`
		SourceFilter      map[string]json.RawMessage `json:"source_filter"`
		DestinationFilter map[string]json.RawMessage `json:"destination_filter"`
	}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if string(got.Port) != `"8080-8090"` {
		t.Errorf("port = %s, want \"8080-8090\"", got.Port)
	}
	if string(got.DestinationFilter["port"]) != `"8080-8090"` {
		t.Errorf("destination_filter.port = %s, want \"8080-8090\"", got.DestinationFilter["port"])
	}
	if _, ok := got.SourceFilter["port"]; ok || strings.Contains(string(b), `"port":""`) {
		t.Errorf("an unset port must be omitted: %s", b)
	}
}
