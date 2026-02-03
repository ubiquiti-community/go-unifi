package types

import (
	"encoding/json"
)

// NumberOrFalse handles fields that can be either `false` (boolean) or a number.
// For example, `aggregated_by` is `false` when not aggregated, or a port index when aggregated.
// The value is stored as int64: 0 for false/not set, >0 for the actual number.
type NumberOrFalse int64

func (n *NumberOrFalse) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		*n = 0
		return nil
	}

	s := string(b)

	// Handle boolean false
	if s == "false" {
		*n = 0
		return nil
	}

	// Handle boolean true (shouldn't happen but be safe)
	if s == "true" {
		*n = 1
		return nil
	}

	// Try to parse as number
	var num int64
	if err := json.Unmarshal(b, &num); err != nil {
		// If it fails, default to 0
		*n = 0
		return nil
	}

	*n = NumberOrFalse(num)
	return nil
}

func (n NumberOrFalse) MarshalJSON() ([]byte, error) {
	if n == 0 {
		return []byte("false"), nil
	}
	return json.Marshal(int64(n))
}

// Int64 returns the value as int64.
func (n NumberOrFalse) Int64() int64 {
	return int64(n)
}

// Bool returns true if the value is non-zero.
func (n NumberOrFalse) Bool() bool {
	return n != 0
}
