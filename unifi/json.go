package unifi

import "github.com/ubiquiti-community/go-unifi/unifi/types"

func emptyBoolToTrue(b *bool) bool {
	if b == nil {
		return true
	}
	return *b
}

// boolValue resolves a tolerant *types.Bool (which accepts both native and
// string-encoded JSON booleans) into a plain bool, defaulting to false when the
// field is absent.
func boolValue(b *types.Bool) bool {
	if b == nil {
		return false
	}
	return b.Bool()
}

// boolPtrValue resolves a tolerant *types.Bool into a *bool, preserving the
// distinction between an absent field (nil) and an explicit false.
func boolPtrValue(b *types.Bool) *bool {
	if b == nil {
		return nil
	}
	v := b.Bool()
	return &v
}
