package unifi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetClientByMAC_CaseInsensitive(t *testing.T) {
	const (
		site         = "default"
		storedMAC    = "02:9d:44:7c:a1:b6"
		requestedMAC = "02:9D:44:7C:A1:B6"
	)

	clientPath := "/proxy/network/api/s/" + site + "/rest/user"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == clientPath {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[{"_id":"client-id","mac":"` + storedMAC + `"}]}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	client, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		t.Fatalf("client init: %v", err)
	}

	got, err := client.GetClientByMAC(context.Background(), site, requestedMAC)
	if err != nil {
		t.Fatalf("GetClientByMAC(%q) errored: %v", requestedMAC, err)
	}
	if got.MAC != storedMAC {
		t.Errorf("MAC = %q, want %q", got.MAC, storedMAC)
	}
}
