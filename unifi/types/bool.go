package types

import (
	"strconv"
	"strings"
)

// Bool handles JSON values that the UniFi controller encodes inconsistently as
// either a native boolean (true/false) or a quoted string ("true"/"false").
// Some controllers (observed on UniFi Network 10.x) return string-encoded
// booleans for fields such as dhcpd_enabled or a client's blocked flag, which
// makes a plain bool field fail to unmarshal. An empty string is treated as
// false.
type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		*b = false
		return nil
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*b = Bool(v)
	return nil
}

// Bool returns the underlying boolean value.
func (b Bool) Bool() bool {
	return bool(b)
}
