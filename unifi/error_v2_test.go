package unifi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestV2APIErrorMessageSurfaced verifies that an error from the v2 API — which
// uses a {"code","message","errorCode"} body, unlike the v1 {"meta":{"rc","msg"}}
// shape — surfaces the controller's actual message instead of a bare HTTP 400.
// Regression guard for the Device Supervisor PoE validation error
// (api.err.PurePoeRequiresUplinkException) showing up as just "(400 Bad Request)".
func TestV2APIErrorMessageSurfaced(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{
			"code": "api.err.PurePoeRequiresUplinkException",
			"message": "Invalid pure PoE supervisor configuration: PORT_NOT_POE_CAPABLE",
			"errorCode": 400
		}`))
	}))
	defer srv.Close()

	c, err := New(context.Background(), &Config{BaseURL: srv.URL, APIKey: "test-key"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var out []PowerSupervisor
	err = c.do(context.Background(), http.MethodGet, "v2/api/site/default/power-supervisors", nil, &out)
	if err == nil {
		t.Fatal("expected an error from the v2 400 response, got nil")
	}
	if !strings.Contains(err.Error(), "PORT_NOT_POE_CAPABLE") {
		t.Errorf("v2 error message not surfaced; got: %v", err)
	}
	if !strings.Contains(err.Error(), "api.err.PurePoeRequiresUplinkException") {
		t.Errorf("v2 error code not surfaced; got: %v", err)
	}
}
