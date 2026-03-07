package unifi

import "github.com/ubiquiti-community/go-unifi/unifi/types"

func emptyBoolToTrue(b *bool) bool {
	if b == nil {
		return true
	}
	return *b
}

// numberToInt64Ptr converts a *types.Number to a *int64.
// If n is nil (field absent from JSON), returns nil.
// If n is non-nil but empty string (e.g. ""), returns a pointer to zero.
// Otherwise returns a pointer to the parsed int64 value.
func numberToInt64Ptr(n *types.Number) *int64 {
	if n == nil {
		return nil
	}
	if val, err := n.Int64(); err == nil {
		return &val
	}
	// Empty string maps to zero per UniFi API convention.
	return new(int64)
}
