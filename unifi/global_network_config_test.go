package unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// A global/config/network body as a Network 11 controller returns it, with the
// keys GlobalNetworkConfig does not model.
const globalNetworkConfigJSON = `{
	"default_security_posture": "allow_all",
	"flood_unknown_multicast_for_network_ids": [],
	"igmp_fastleave_for_network_ids": [],
	"igmp_querier_switches": [],
	"igmp_snooping_for": "some",
	"igmp_snooping_for_network_ids": ["net-a", "net-b"],
	"ipv6_pd_interfaces": [],
	"mdns_enabled_for": "all",
	"mdns_enabled_for_network_ids": ["net-a", "net-b", "net-c"]
}`

// globalNetworkConfigTestServer serves getBody at the v2 path and records the
// body of each PUT, which it echoes back as the stored object. A nil getBody
// answers 404, as a controller without the endpoint does.
func globalNetworkConfigTestServer(t *testing.T, getBody *string, put *[]byte) *httptest.Server {
	t.Helper()
	path := "/proxy/network/v2/api/site/default/global/config/network"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		switch {
		case r.URL.Path == path && getBody == nil:
			w.WriteHeader(http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == path:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(*getBody))
		case r.Method == http.MethodPut && r.URL.Path == path:
			b, _ := io.ReadAll(r.Body)
			*put = b
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(b)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newGlobalNetworkConfigTestClient(t *testing.T, srv *httptest.Server) *ApiClient {
	t.Helper()
	c, err := New(context.Background(), &Config{BaseURL: srv.URL, Username: "admin", Password: "admin"})
	if err != nil {
		t.Fatalf("client init: %v", err)
	}
	return c
}

func decodeObject(t *testing.T, b []byte) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode %s: %v", b, err)
	}
	return m
}

// TestUpdateGlobalNetworkConfig_KeepsUnmodelledKeys proves an update sends back
// every key it read, so changing the mDNS selection cannot reset the IGMP,
// multicast or security-posture settings that share the object.
func TestUpdateGlobalNetworkConfig_KeepsUnmodelledKeys(t *testing.T) {
	body := globalNetworkConfigJSON
	var put []byte
	c := newGlobalNetworkConfigTestClient(t, globalNetworkConfigTestServer(t, &body, &put))
	ctx := context.Background()

	cfg, err := c.GetGlobalNetworkConfig(ctx, "default")
	if err != nil {
		t.Fatalf("GetGlobalNetworkConfig: %v", err)
	}
	if !cfg.HasMdns() {
		t.Fatal("HasMdns() = false, want true")
	}
	if cfg.MdnsEnabledFor != "all" || len(cfg.MdnsEnabledForNetworkIDs) != 3 {
		t.Fatalf("decoded mdns = %q %v", cfg.MdnsEnabledFor, cfg.MdnsEnabledForNetworkIDs)
	}

	cfg.MdnsEnabledFor = "some"
	cfg.MdnsEnabledForNetworkIDs = []string{"net-a"}
	got, err := c.UpdateGlobalNetworkConfig(ctx, "default", cfg)
	if err != nil {
		t.Fatalf("UpdateGlobalNetworkConfig: %v", err)
	}
	if got.MdnsEnabledFor != "some" || !reflect.DeepEqual(got.MdnsEnabledForNetworkIDs, []string{"net-a"}) {
		t.Errorf("returned mdns = %q %v", got.MdnsEnabledFor, got.MdnsEnabledForNetworkIDs)
	}

	sent := decodeObject(t, put)
	orig := decodeObject(t, []byte(globalNetworkConfigJSON))
	for k, v := range orig {
		if k == keyMdnsEnabledFor || k == keyMdnsEnabledForIDs {
			continue
		}
		var want bytes.Buffer
		if err := json.Compact(&want, v); err != nil {
			t.Fatal(err)
		}
		if string(sent[k]) != want.String() {
			t.Errorf("key %s sent as %s, want %s unchanged", k, sent[k], want.String())
		}
	}
	if string(sent[keyMdnsEnabledFor]) != `"some"` {
		t.Errorf("mdns_enabled_for sent as %s", sent[keyMdnsEnabledFor])
	}
	if string(sent[keyMdnsEnabledForIDs]) != `["net-a"]` {
		t.Errorf("mdns_enabled_for_network_ids sent as %s", sent[keyMdnsEnabledForIDs])
	}
	if len(sent) != len(orig) {
		t.Errorf("sent %d keys, want %d", len(sent), len(orig))
	}
}

// TestGetGlobalNetworkConfig_NotFound verifies a controller without the v2
// endpoint yields a typed NotFoundError, so callers can fall back to the
// legacy mdns setting.
func TestGetGlobalNetworkConfig_NotFound(t *testing.T) {
	c := newGlobalNetworkConfigTestClient(t, globalNetworkConfigTestServer(t, nil, nil))

	_, err := c.GetGlobalNetworkConfig(context.Background(), "default")
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("err = %v, want *NotFoundError", err)
	}
}

// TestGlobalNetworkConfig_HasMdnsFalseWithoutKey covers a controller that
// serves the object but keeps the mDNS selection elsewhere.
func TestGlobalNetworkConfig_HasMdnsFalseWithoutKey(t *testing.T) {
	var cfg GlobalNetworkConfig
	if err := json.Unmarshal([]byte(`{"igmp_snooping_for":"all"}`), &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.HasMdns() {
		t.Error("HasMdns() = true for an object without mdns_enabled_for")
	}
}

// TestGlobalNetworkConfig_NilNetworkIDsSentAsEmptyList checks a nil selection
// is written as [] rather than null, as the UI sends it.
func TestGlobalNetworkConfig_NilNetworkIDsSentAsEmptyList(t *testing.T) {
	b, err := json.Marshal(GlobalNetworkConfig{MdnsEnabledFor: "none"})
	if err != nil {
		t.Fatal(err)
	}
	m := decodeObject(t, b)
	if string(m[keyMdnsEnabledForIDs]) != `[]` {
		t.Errorf("mdns_enabled_for_network_ids = %s, want []", m[keyMdnsEnabledForIDs])
	}
}
