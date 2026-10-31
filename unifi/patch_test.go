package unifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
)

func diffTestWLAN() WLAN {
	return WLAN{
		ID:                    "wlan1",
		SiteID:                "site1",
		Name:                  "guest",
		Enabled:               true,
		NetworkID:             "net1",
		MACFilterList:         []string{"02:00:00:00:00:01"},
		MinrateNaEnabled:      true,
		MinrateNaDataRateKbps: Ptr[int64](6000),
		ScheduleWithDuration: []WLANScheduleWithDuration{
			{Name: "weekday", StartDaysOfWeek: []string{"mon"}, StartHour: Ptr[int64](8)},
		},
	}
}

func TestDiff(t *testing.T) {
	minrate := []string{"minrate_na_enabled", "minrate_na_data_rate_kbps"}

	tests := []struct {
		name    string
		prior   func(*WLAN)
		target  func(*WLAN)
		opts    DiffOptions
		want    string
		wantErr string
	}{
		{name: "unchanged", target: func(*WLAN) {}, want: `{}`},
		{name: "changed", target: func(w *WLAN) { w.Name = "guest-2" }, want: `{"name":"guest-2"}`},
		{name: "added", target: func(w *WLAN) { w.VLAN = Ptr[int64](20) }, want: `{"vlan":20}`},
		{
			name: "nested value sent whole",
			target: func(w *WLAN) {
				w.ScheduleWithDuration = []WLANScheduleWithDuration{
					{Name: "weekday", StartDaysOfWeek: []string{"mon"}, StartHour: Ptr[int64](9)},
				}
			},
			want: `{"schedule_with_duration":[{"name":"weekday","start_days_of_week":["mon"],"start_hour":9}]}`,
		},
		{name: "removed without Clear", target: func(w *WLAN) { w.MACFilterList = nil }, want: `{}`},
		{
			name:   "Clear string",
			target: func(w *WLAN) { w.NetworkID = "" },
			opts:   DiffOptions{Clear: []string{"networkconf_id"}},
			want:   `{"networkconf_id":""}`,
		},
		{
			name:   "Clear list",
			target: func(w *WLAN) { w.MACFilterList = nil },
			opts:   DiffOptions{Clear: []string{"mac_filter_list"}},
			want:   `{"mac_filter_list":[]}`,
		},
		{
			name:   "Clear with nothing stored",
			target: func(*WLAN) {},
			opts:   DiffOptions{Clear: []string{"ap_group_ids"}},
			want:   `{}`,
		},
		{
			name:    "Clear number",
			prior:   func(w *WLAN) { w.VLAN = Ptr[int64](20) },
			target:  func(w *WLAN) { w.VLAN = nil },
			opts:    DiffOptions{Clear: []string{"vlan"}},
			wantErr: `unable to clear "vlan"`,
		},
		{
			name:   "group sends unchanged companion",
			target: func(w *WLAN) { w.MinrateNaDataRateKbps = Ptr[int64](12000) },
			opts:   DiffOptions{Groups: [][]string{minrate}},
			want:   `{"minrate_na_data_rate_kbps":12000,"minrate_na_enabled":true}`,
		},
		{
			name:   "group without change",
			target: func(w *WLAN) { w.Name = "guest-2" },
			opts:   DiffOptions{Groups: [][]string{minrate}},
			want:   `{"name":"guest-2"}`,
		},
		{name: "ids never sent", target: func(w *WLAN) { w.ID, w.SiteID = "wlan2", "site2" }, want: `{}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prior, target := diffTestWLAN(), diffTestWLAN()
			if tt.prior != nil {
				tt.prior(&prior)
			}
			tt.target(&target)

			patch, err := Diff(&prior, &target, tt.opts)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want one containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Diff: %v", err)
			}
			got, err := json.Marshal(patch)
			if err != nil {
				t.Fatalf("marshal patch: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("patch = %s, want %s", got, tt.want)
			}
		})
	}
}

// A LAN patch never holds wan_ip, because the corporate encoder never sends it.
func TestDiff_NetworkUsesPurposeEncoder(t *testing.T) {
	prior := Network{
		ID:         "lan1",
		Name:       Ptr("lan"),
		Purpose:    PurposeCorporate,
		IPSubnet:   Ptr("198.51.100.1/24"),
		DHCPDStart: Ptr("198.51.100.10"),
		DHCPDStop:  Ptr("198.51.100.200"),
	}
	target := prior
	target.DHCPDStart = Ptr("198.51.100.20")
	target.WANIP = Ptr("203.0.113.2")

	patch, err := Diff(&prior, &target, DiffOptions{})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if want := map[string]any{"dhcpd_start": "198.51.100.20"}; !reflect.DeepEqual(patch, want) {
		t.Errorf("patch = %v, want %v", patch, want)
	}
}

// TestDiff_WANInterfaceMTU checks that WAN MTU and IGMP proxy changes reach a
// sparse PUT: Diff encodes through marshalWAN, which used to drop these keys.
func TestDiff_WANInterfaceMTU(t *testing.T) {
	prior := Network{
		ID:                            "wan1",
		Name:                          Ptr("WAN"),
		Purpose:                       PurposeWAN,
		WANType:                       Ptr("dhcp"),
		InterfaceMtu:                  Ptr[int64](1500),
		IGMPProxyFor:                  Ptr("some"),
		IGMPProxyDownstreamNetworkIDs: []string{"lan1"},
	}

	target := prior
	target.InterfaceMtu = Ptr[int64](1492)
	target.InterfaceMtuEnabled = true
	patch, err := Diff(&prior, &target, DiffOptions{})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	want := map[string]any{"interface_mtu": json.Number("1492"), "interface_mtu_enabled": true}
	if !reflect.DeepEqual(patch, want) {
		t.Errorf("patch = %v, want %v", patch, want)
	}

	target = prior
	target.IGMPProxyDownstreamNetworkIDs = []string{"lan1", "lan2"}
	patch, err = Diff(&prior, &target, DiffOptions{})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	want = map[string]any{"igmp_proxy_downstream_networkconf_ids": []any{"lan1", "lan2"}}
	if !reflect.DeepEqual(patch, want) {
		t.Errorf("patch = %v, want %v", patch, want)
	}
}

const (
	patchTestWLANPath = "/proxy/network/api/s/default/rest/wlanconf/wlan1"
	patchTestLANPath  = "/proxy/network/api/s/default/rest/networkconf/lan1"

	patchTestWLAN = `{
		"_id": "wlan1", "site_id": "site1", "name": "guest", "enabled": true,
		"security": "wpapsk", "wpa_mode": "wpa2", "x_passphrase": "fake-passphrase",
		"networkconf_id": "net1", "mac_filter_enabled": true, "mac_filter_policy": "allow",
		"mac_filter_list": ["02:00:00:00:00:01", "02:00:00:00:00:02"],
		"schedule_enabled": true,
		"schedule_with_duration": [
			{"name": "weekday", "start_days_of_week": ["mon"], "start_hour": 8, "start_minute": 0, "duration_minutes": 600},
			{"name": "weekend", "start_days_of_week": ["sat"], "start_hour": 10, "start_minute": 0, "duration_minutes": 300}
		],
		"controller_only_key": "kept"
	}`

	patchTestLAN = `{
		"_id": "lan1", "site_id": "site1", "name": "lan", "purpose": "corporate", "enabled": true,
		"ip_subnet": "198.51.100.129/25", "vlan": 52, "vlan_enabled": true,
		"dhcpd_enabled": true, "dhcpd_start": "198.51.100.130", "dhcpd_stop": "198.51.100.250",
		"dhcpd_ntp_enabled": false
	}`
)

// patchTestController answers like the 10.6.106 controller: each top-level key
// of a PUT replaces the stored value whole, a no-op PUT returns empty data, and
// a merged document that fails validation returns 400 and stores nothing.
type patchTestController struct {
	mu       sync.Mutex
	docs     map[string]map[string]any
	validate func(doc map[string]any) string
	puts     []string
	gets     int
}

func newPatchTestController(t *testing.T, docs map[string]string) (*patchTestController, *ApiClient) {
	t.Helper()
	f := &patchTestController{docs: map[string]map[string]any{}}
	for path, doc := range docs {
		var m map[string]any
		if err := json.Unmarshal([]byte(doc), &m); err != nil {
			t.Fatalf("fixture %s: %v", path, err)
		}
		f.docs[path] = m
	}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	c, err := New(context.Background(), &Config{BaseURL: srv.URL, APIKey: "test-key"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return f, c
}

func (f *patchTestController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handleNewStyleSetup(w, r) {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	doc, ok := f.docs[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		f.gets++
		_ = json.NewEncoder(w).Encode(map[string]any{"meta": map[string]any{"rc": "ok"}, "data": []any{doc}})
	case http.MethodPut:
		body, _ := io.ReadAll(r.Body)
		f.puts = append(f.puts, string(body))
		var patch map[string]any
		if err := json.Unmarshal(body, &patch); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		merged := maps.Clone(doc)
		maps.Copy(merged, patch)
		if f.validate != nil {
			if code := f.validate(merged); code != "" {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = fmt.Fprintf(w, `{"meta":{"rc":"error","msg":%q},"data":[]}`, code)
				return
			}
		}
		if reflect.DeepEqual(merged, doc) {
			_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
			return
		}
		f.docs[r.URL.Path] = merged
		_ = json.NewEncoder(w).Encode(map[string]any{"meta": map[string]any{"rc": "ok"}, "data": []any{merged}})
	}
}

func (f *patchTestController) stored(path, key string) any {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.docs[path][key]
}

func (f *patchTestController) requests() ([]string, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.puts), f.gets
}

func TestUpdateWLANPatch_SendsOnlyChangedKeys(t *testing.T) {
	ctx := context.Background()
	f, c := newPatchTestController(t, map[string]string{patchTestWLANPath: patchTestWLAN})

	prior, err := c.GetWLAN(ctx, "default", "wlan1")
	if err != nil {
		t.Fatalf("GetWLAN: %v", err)
	}
	target := *prior
	target.Name = "guest-renamed"

	patch, err := Diff(prior, &target, DiffOptions{})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	got, err := c.UpdateWLANPatch(ctx, "default", "wlan1", patch)
	if err != nil {
		t.Fatalf("UpdateWLANPatch: %v", err)
	}

	puts, _ := f.requests()
	if want := []string{`{"name":"guest-renamed"}`}; !slices.Equal(puts, want) {
		t.Errorf("PUT bodies = %q, want %q", puts, want)
	}
	if got.Name != "guest-renamed" {
		t.Errorf("returned name = %q, want guest-renamed", got.Name)
	}
	if v := f.stored(patchTestWLANPath, "controller_only_key"); v != "kept" {
		t.Errorf("controller_only_key = %v, want kept", v)
	}
	if v, _ := f.stored(patchTestWLANPath, "mac_filter_list").([]any); len(v) != 2 {
		t.Errorf("mac_filter_list = %v, want both entries", v)
	}
}

func TestUpdateWLANPatch_ClearsListAndReplacesNestedValue(t *testing.T) {
	ctx := context.Background()
	f, c := newPatchTestController(t, map[string]string{patchTestWLANPath: patchTestWLAN})

	prior, err := c.GetWLAN(ctx, "default", "wlan1")
	if err != nil {
		t.Fatalf("GetWLAN: %v", err)
	}
	target := *prior
	target.MACFilterList = nil
	target.ScheduleWithDuration = prior.ScheduleWithDuration[:1]

	patch, err := Diff(prior, &target, DiffOptions{Clear: []string{"mac_filter_list"}})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if _, err := c.UpdateWLANPatch(ctx, "default", "wlan1", patch); err != nil {
		t.Fatalf("UpdateWLANPatch: %v", err)
	}

	puts, _ := f.requests()
	want := `{"mac_filter_list":[],"schedule_with_duration":[{"duration_minutes":600,"name":"weekday",` +
		`"start_days_of_week":["mon"],"start_hour":8,"start_minute":0}]}`
	if !slices.Equal(puts, []string{want}) {
		t.Errorf("PUT bodies = %q, want %q", puts, want)
	}
	if v, ok := f.stored(patchTestWLANPath, "mac_filter_list").([]any); !ok || len(v) != 0 {
		t.Errorf("mac_filter_list = %v, want []", v)
	}
	if v, _ := f.stored(patchTestWLANPath, "schedule_with_duration").([]any); len(v) != 1 {
		t.Errorf("schedule_with_duration = %v, want one entry", v)
	}
}

func TestUpdateWLANPatch_ReadsBackWithGet(t *testing.T) {
	tests := []struct {
		name     string
		patch    map[string]any
		wantPuts int
	}{
		{name: "empty patch sends no PUT", patch: map[string]any{}, wantPuts: 0},
		{name: "PUT without change returns empty data", patch: map[string]any{"name": "guest"}, wantPuts: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, c := newPatchTestController(t, map[string]string{patchTestWLANPath: patchTestWLAN})

			got, err := c.UpdateWLANPatch(context.Background(), "default", "wlan1", tt.patch)
			if err != nil {
				t.Fatalf("UpdateWLANPatch: %v", err)
			}
			puts, gets := f.requests()
			if len(puts) != tt.wantPuts || gets != 1 {
				t.Errorf("PUTs = %d, GETs = %d, want %d and 1", len(puts), gets, tt.wantPuts)
			}
			if got.Name != "guest" {
				t.Errorf("returned name = %q, want guest", got.Name)
			}
		})
	}
}

// The stored LAN lacks dhcpd_ntp_1, which the prior encoding holds. Sending
// dhcpd_ntp_enabled alone then fails validation of the merged document, so the
// NTP group sends the unchanged address with the flag.
func TestUpdateNetworkPatch_NTPGroup(t *testing.T) {
	ntp := []string{"dhcpd_ntp_enabled", "dhcpd_ntp_1", "dhcpd_ntp_2"}

	tests := []struct {
		name     string
		opts     DiffOptions
		wantBody string
		wantErr  string
	}{
		{
			name:     "without group",
			wantBody: `{"dhcpd_ntp_enabled":true}`,
			wantErr:  "api.err.NtpAddressInvalid",
		},
		{
			name:     "with group",
			opts:     DiffOptions{Groups: [][]string{ntp}},
			wantBody: `{"dhcpd_ntp_1":"198.51.100.1","dhcpd_ntp_enabled":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			f, c := newPatchTestController(t, map[string]string{patchTestLANPath: patchTestLAN})
			f.validate = func(doc map[string]any) string {
				if addr, _ := doc["dhcpd_ntp_1"].(string); doc["dhcpd_ntp_enabled"] == true && addr == "" {
					return "api.err.NtpAddressInvalid"
				}
				return ""
			}

			prior, err := c.GetNetwork(ctx, "default", "lan1")
			if err != nil {
				t.Fatalf("GetNetwork: %v", err)
			}
			prior.DHCPDNtp1 = Ptr("198.51.100.1")
			target := *prior
			target.DHCPDNtpEnabled = true

			patch, err := Diff(prior, &target, tt.opts)
			if err != nil {
				t.Fatalf("Diff: %v", err)
			}
			got, err := c.UpdateNetworkPatch(ctx, "default", "lan1", patch)

			puts, _ := f.requests()
			if !slices.Equal(puts, []string{tt.wantBody}) {
				t.Errorf("PUT bodies = %q, want %q", puts, tt.wantBody)
			}
			if tt.wantErr != "" {
				apiErr := &APIError{}
				if !errors.As(err, &apiErr) || apiErr.Message != tt.wantErr {
					t.Fatalf("error = %v, want APIError %s", err, tt.wantErr)
				}
				if v := f.stored(patchTestLANPath, "dhcpd_ntp_enabled"); v != false {
					t.Errorf("dhcpd_ntp_enabled = %v after a rejected PUT, want false", v)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateNetworkPatch: %v", err)
			}
			if !got.DHCPDNtpEnabled || got.DHCPDNtp1 == nil || *got.DHCPDNtp1 != "198.51.100.1" {
				t.Errorf("returned NTP = %v %v, want enabled with 198.51.100.1", got.DHCPDNtpEnabled, got.DHCPDNtp1)
			}
		})
	}
}
