package types

import (
	"encoding/json"
	"fmt"
)

// NumberOrString handles fields that the UniFi controller encodes as either a
// JSON string or a JSON number, such as the `date` of a device SSH key (a
// number on UniFi Network 11.0.104, where the SDK had expected a string). The
// value keeps the text it was decoded from. On encoding, a value that reads as
// a JSON number is written as a number and anything else as a string, so a
// read-modify-write sends the field back in the form the controller sent it.
type NumberOrString string

func (n *NumberOrString) UnmarshalJSON(b []byte) error {
	// A JSON string, or null (which leaves s empty).
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*n = NumberOrString(s)
		return nil
	}

	var num json.Number
	if err := json.Unmarshal(b, &num); err != nil {
		return fmt.Errorf("unable to unmarshal %s as a number or a string", b)
	}
	*n = NumberOrString(num)
	return nil
}

func (n NumberOrString) MarshalJSON() ([]byte, error) {
	if n.IsNumber() {
		return []byte(n), nil
	}
	return json.Marshal(string(n))
}

// IsNumber reports whether the value is a JSON number, such as 1732000000.
func (n NumberOrString) IsNumber() bool {
	s := string(n)
	if s == "" || (s[0] != '-' && (s[0] < '0' || s[0] > '9')) {
		return false
	}
	var num json.Number
	return json.Unmarshal([]byte(s), &num) == nil && num.String() == s
}

// String returns the value as text.
func (n NumberOrString) String() string {
	return string(n)
}
