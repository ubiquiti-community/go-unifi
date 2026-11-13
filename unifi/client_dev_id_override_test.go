package unifi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dev_id_override is the device type a client was manually set to. UniFi
// Network 11 sends it as a numeric string ("3263"); earlier controllers sent a
// number.
func TestClientUnmarshalDevIdOverride(t *testing.T) {
	t.Parallel()

	v := int64(3263)
	zero := int64(0)
	cases := map[string]struct {
		input string
		want  *int64
	}{
		"number":         {`{"mac":"aa:bb:cc:dd:ee:ff","dev_id_override":3263}`, &v},
		"numeric string": {`{"mac":"aa:bb:cc:dd:ee:ff","dev_id_override":"3263"}`, &v},
		"empty string":   {`{"mac":"aa:bb:cc:dd:ee:ff","dev_id_override":""}`, &zero},
		"null":           {`{"mac":"aa:bb:cc:dd:ee:ff","dev_id_override":null}`, nil},
		"absent":         {`{"mac":"aa:bb:cc:dd:ee:ff"}`, nil},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var c Client
			require.NoError(t, json.Unmarshal([]byte(tc.input), &c))
			assert.Equal(t, tc.want, c.DevIdOverride)
			assert.Equal(t, "aa:bb:cc:dd:ee:ff", c.MAC)
		})
	}
}

func TestClientMarshalOmitsUnsetDevIdOverride(t *testing.T) {
	t.Parallel()

	b, err := json.Marshal(&Client{MAC: "aa:bb:cc:dd:ee:ff"})
	require.NoError(t, err)
	assert.NotContains(t, string(b), "dev_id_override")
}
