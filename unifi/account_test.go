package unifi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubiquiti-community/go-unifi/unifi"
)

func TestAccountMarshalJSON(t *testing.T) {
	for n, c := range map[string]struct {
		expectedJSON string
		acc          unifi.Account
	}{
		"empty (nil pointers omitted)": {
			`{}`,
			unifi.Account{},
		},
		"response": {
			`{"vlan":10,"tunnel_type":1,"tunnel_medium_type":1}`,
			unifi.Account{
				VLAN:             unifi.Ptr[int64](10),
				TunnelType:       unifi.Ptr[int64](1),
				TunnelMediumType: unifi.Ptr[int64](1),
			},
		},
	} {
		t.Run(n, func(t *testing.T) {
			actual, err := json.Marshal(&c.acc)
			if err != nil {
				t.Fatal(err)
			}
			assert.JSONEq(t, c.expectedJSON, string(actual))
		})
	}
}
