package settings

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Ad blocking is a site-wide switch plus the networks it applies to, stored on
// the ips setting. The keys are present on Network 10.x and 11 controllers but
// absent from the ace.jar field spec.
func TestIpsAdBlockingRoundTrip(t *testing.T) {
	t.Parallel()

	in := `{"key":"ips","ad_blocking_enabled":true,` +
		`"ad_blocking_configurations":[{"network_id":"5341251b-e57e-4b74-9db3-4933182ef18e"}]}`

	var s Ips
	require.NoError(t, json.Unmarshal([]byte(in), &s))
	require.NotNil(t, s.AdBlockingEnabled)
	assert.True(t, *s.AdBlockingEnabled)
	require.Len(t, s.AdBlockingConfigurations, 1)
	assert.Equal(t, "5341251b-e57e-4b74-9db3-4933182ef18e", s.AdBlockingConfigurations[0].NetworkID)

	out, err := json.Marshal(&s)
	require.NoError(t, err)
	assert.Contains(t, string(out), `"ad_blocking_enabled":true`)
	assert.Contains(t, string(out), `"ad_blocking_configurations":[{"network_id":"5341251b-e57e-4b74-9db3-4933182ef18e"}]`)
}

func TestIpsAdBlockingOmittedWhenUnset(t *testing.T) {
	t.Parallel()

	out, err := json.Marshal(&Ips{})
	require.NoError(t, err)
	assert.NotContains(t, string(out), "ad_blocking")
}
