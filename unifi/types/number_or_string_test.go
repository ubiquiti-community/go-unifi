package types

import (
	"encoding/json"
	"testing"
)

func TestNumberOrStringRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want NumberOrString
	}{
		{"integer", `1732000000`, "1732000000"},
		{"float", `1732000000.5`, "1732000000.5"},
		{"negative", `-1`, "-1"},
		{"string", `"2024-01-31T12:00:00Z"`, "2024-01-31T12:00:00Z"},
		{"empty string", `""`, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got NumberOrString
			if err := json.Unmarshal([]byte(tc.in), &got); err != nil {
				t.Fatalf("unmarshal %s: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("unmarshal %s = %q, want %q", tc.in, got, tc.want)
			}

			b, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("marshal %q: %v", got, err)
			}
			if string(b) != tc.in {
				t.Errorf("marshal %q = %s, want %s", got, b, tc.in)
			}
		})
	}
}

func TestNumberOrStringNull(t *testing.T) {
	got := NumberOrString("stale")
	if err := json.Unmarshal([]byte(`null`), &got); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if got != "" {
		t.Errorf("unmarshal null = %q, want empty", got)
	}
}

func TestNumberOrStringRejectsOtherTypes(t *testing.T) {
	for _, in := range []string{`true`, `{}`, `[1]`} {
		var got NumberOrString
		if err := json.Unmarshal([]byte(in), &got); err == nil {
			t.Errorf("unmarshal %s = %q, want an error", in, got)
		}
	}
}

func TestNumberOrStringIsNumber(t *testing.T) {
	for in, want := range map[NumberOrString]bool{
		"0":          true,
		"1e5":        true,
		"-12.5":      true,
		"":           false,
		"01":         false,
		"12 ":        false,
		" 12":        false,
		"0x10":       false,
		"auto":       false,
		"2024-01-31": false,
	} {
		if got := in.IsNumber(); got != want {
			t.Errorf("NumberOrString(%q).IsNumber() = %v, want %v", in, got, want)
		}
	}
}
