package unifi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// wanSlaTestServer is a new-style (UniFi OS) controller holding WAN SLA profiles the way Network 11 does: POST assigns
// _id and monitor ids, per-id GET/PUT/DELETE work, DELETE answers with an empty body, and a missing id is a 404 with
// api.err.WanSlaNotFound.
func wanSlaTestServer(t *testing.T, site string) (*httptest.Server, map[string]WANSla) {
	t.Helper()
	store := map[string]WANSla{}
	base := "/proxy/network/v2/api/site/" + site + "/wan-slas"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		writeJSON := func(v any) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(v)
		}
		id := ""
		if len(r.URL.Path) > len(base)+1 {
			id = r.URL.Path[len(base)+1:]
		}
		var in WANSla
		if r.Body != nil {
			b, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(b, &in)
		}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == base:
			out := []WANSla{}
			for _, v := range store {
				out = append(out, v)
			}
			writeJSON(out)
		case r.Method == http.MethodPost && r.URL.Path == base:
			in.ID = "00000000-0000-4000-8000-000000000001"
			for i := range in.Monitors {
				n := int64(i + 1)
				in.Monitors[i].MonitorID = &n
			}
			store[in.ID] = in
			writeJSON(in)
		case id != "":
			if _, ok := store[id]; !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"code":"api.err.WanSlaNotFound","errorCode":404,"message":"WAN SLA not found"}`))
				return
			}
			switch r.Method {
			case http.MethodGet:
				writeJSON(store[id])
			case http.MethodPut:
				store[id] = in
				writeJSON(in)
			case http.MethodDelete:
				delete(store, id)
				w.WriteHeader(http.StatusOK) // empty body, as the controller does
			}
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, store
}

func ptr64(v int64) *int64 { return &v }

func TestWANSla_CRUD(t *testing.T) {
	const site = "default"
	srv, store := wanSlaTestServer(t, site)
	c, err := New(context.Background(), &Config{BaseURL: srv.URL, Username: "admin", Password: "admin"})
	if err != nil {
		t.Fatalf("client init: %v", err)
	}
	ctx := context.Background()

	created, err := c.CreateWANSla(ctx, site, &WANSla{
		Name:              "wan-standard",
		MonitorHealthType: WANSlaMonitorHealthMajority,
		Monitors: []WANSlaMonitor{
			{Type: WANSlaMonitorTypeICMP, Target: "1.1.1.1", Interval: 2, TimePeriod: 30, Alert: WANSlaAlert{
				ThresholdPolicy: "loss-or-all-others", LossThreshold: ptr64(5), LatencyThreshold: ptr64(150), JitterThreshold: ptr64(30),
			}},
			{Type: WANSlaMonitorTypeDNS, Target: "8.8.8.8", Domain: "example.com", Interval: 5, TimePeriod: 30, Alert: WANSlaAlert{
				ThresholdPolicy: "any", LossThreshold: ptr64(5),
			}},
		},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" || created.Monitors[1].MonitorID == nil || *created.Monitors[1].MonitorID != 2 {
		t.Fatalf("create should return the id and monitor ids, got %+v", created)
	}

	got, err := c.GetWANSla(ctx, site, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Monitors[1].Domain != "example.com" || got.Monitors[1].Alert.LatencyThreshold != nil {
		t.Fatalf("get: domain kept and an unset threshold stays unset, got %+v", got.Monitors[1])
	}

	got.Name = "wan-standard-2"
	if _, err := c.UpdateWANSla(ctx, site, got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if store[created.ID].Name != "wan-standard-2" || *store[created.ID].Monitors[0].MonitorID != 1 {
		t.Fatalf("update should send the name and keep monitor ids, stored %+v", store[created.ID])
	}

	list, err := c.ListWANSla(ctx, site)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v %v", list, err)
	}

	if err := c.DeleteWANSla(ctx, site, created.ID); err != nil {
		t.Fatalf("delete (empty body): %v", err)
	}
	_, err = c.GetWANSla(ctx, site, created.ID)
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("get after delete should be *NotFoundError, got %v", err)
	}
}

// The JSON the controller sends must decode, including the monitor ids it assigns as numbers and a monitor without a
// domain.
func TestWANSla_DecodesControllerJSON(t *testing.T) {
	raw := `{"_id":"00000000-0000-4000-8000-000000000001","monitor_health_type":"MAJORITY","monitors":[{"alert":{"jitter_threshold":30,"latency_threshold":150,"loss_threshold":5,"threshold_policy":"loss-or-all-others"},"interval":2,"monitor_id":1,"target":"1.1.1.1","time_period":30,"type":"ICMP"}],"name":"api-probe"}`
	var s WANSla
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	m := s.Monitors[0]
	if *m.MonitorID != 1 || *m.Alert.LossThreshold != 5 || m.Alert.ThresholdPolicy != "loss-or-all-others" || m.Domain != "" {
		t.Fatalf("decoded %+v", m)
	}
}

// A WAN network must carry wan_sla when it is encoded: the WAN encoder lists its fields explicitly.
func TestWANSla_NetworkEncodesWANSla(t *testing.T) {
	id := "00000000-0000-4000-8000-000000000001"
	b, err := json.Marshal(&Network{Purpose: PurposeWAN, WANSla: &id})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if m["wan_sla"] != id {
		t.Fatalf("wan_sla = %v in %s", m["wan_sla"], b)
	}
	empty := ""
	b, _ = json.Marshal(&Network{Purpose: PurposeWAN, WANSla: &empty})
	_ = json.Unmarshal(b, &m)
	if v, ok := m["wan_sla"]; !ok || v != "" {
		t.Fatalf("an empty wan_sla (detach) must be sent, got %s", b)
	}
}
