package unifi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// stamgrTestServer stands up a new-style (UniFi OS) controller whose stamgr
// endpoint replies with the given status and body, capturing the last request
// body it received so tests can assert on the exact wire format.
func stamgrTestServer(t *testing.T, site string, status int, respJSON string) (*httptest.Server, *map[string]any) {
	t.Helper()
	lastBody := &map[string]any{}
	stamgrPath := "/proxy/network/api/s/" + site + "/cmd/stamgr"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == stamgrPath {
			body := map[string]any{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			*lastBody = body
			w.Header().Set("Content-Type", "application/json")
			if status != http.StatusOK {
				w.WriteHeader(status)
			}
			_, _ = w.Write([]byte(respJSON))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	return srv, lastBody
}

func newStamgrTestClient(t *testing.T, srv *httptest.Server) *ApiClient {
	t.Helper()
	c, err := New(context.Background(), &Config{BaseURL: srv.URL, Username: "admin", Password: "admin"})
	if err != nil {
		t.Fatalf("client init: %v", err)
	}
	return c
}

// TestAuthorizeClientByMAC_Success covers the classic controller reply where
// the authorized client object is echoed back in data.
func TestAuthorizeClientByMAC_Success(t *testing.T) {
	srv, lastBody := stamgrTestServer(t, "default", http.StatusOK,
		`{"meta":{"rc":"ok"},"data":[{"mac":"aa:bb:cc:dd:ee:ff"}]}`)
	c := newStamgrTestClient(t, srv)

	err := c.AuthorizeClientByMAC(context.Background(), "default", "aa:bb:cc:dd:ee:ff", "11:22:33:44:55:66", "480")
	if err != nil {
		t.Fatalf("AuthorizeClientByMAC errored: %v", err)
	}

	body := *lastBody
	if body["cmd"] != "authorize-guest" {
		t.Errorf("cmd = %v, want authorize-guest", body["cmd"])
	}
	if body["mac"] != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("mac = %v, want aa:bb:cc:dd:ee:ff", body["mac"])
	}
	if body["ap_mac"] != "11:22:33:44:55:66" {
		t.Errorf("ap_mac = %v, want 11:22:33:44:55:66", body["ap_mac"])
	}
	// The controller rejects string minutes: the value must be a JSON number.
	if m, ok := body["minutes"].(float64); !ok || m != 480 {
		t.Errorf("minutes = %v (%T), want JSON number 480", body["minutes"], body["minutes"])
	}
}

// TestAuthorizeClientByMAC_EmptyDataOK locks in the fix for UniFi OS builds
// that answer authorize-guest with rc:"ok" and an empty data array even though
// the client was authorized; the naive len(users)!=1 check misreads that as
// not-found.
func TestAuthorizeClientByMAC_EmptyDataOK(t *testing.T) {
	srv, _ := stamgrTestServer(t, "default", http.StatusOK,
		`{"meta":{"rc":"ok"},"data":[]}`)
	c := newStamgrTestClient(t, srv)

	err := c.AuthorizeClientByMAC(context.Background(), "default", "aa:bb:cc:dd:ee:ff", "", "")
	if err != nil {
		t.Fatalf("rc:ok with empty data should be success, got: %v", err)
	}
}

// TestAuthorizeClientByMAC_OmitsUnparseableMinutes verifies that an empty or
// non-numeric minutes value is omitted so the controller applies its default.
func TestAuthorizeClientByMAC_OmitsUnparseableMinutes(t *testing.T) {
	srv, lastBody := stamgrTestServer(t, "default", http.StatusOK,
		`{"meta":{"rc":"ok"},"data":[]}`)
	c := newStamgrTestClient(t, srv)

	if err := c.AuthorizeClientByMAC(context.Background(), "default", "aa:bb:cc:dd:ee:ff", "", "unlimited"); err != nil {
		t.Fatalf("AuthorizeClientByMAC errored: %v", err)
	}
	if _, present := (*lastBody)["minutes"]; present {
		t.Errorf("minutes should be omitted for unparseable input, body = %v", *lastBody)
	}
}

// TestAuthorizeClientByMAC_APIError surfaces a controller error response.
func TestAuthorizeClientByMAC_APIError(t *testing.T) {
	srv, _ := stamgrTestServer(t, "default", http.StatusBadRequest,
		`{"meta":{"rc":"error","msg":"api.err.UnknownStation"},"data":[]}`)
	c := newStamgrTestClient(t, srv)

	err := c.AuthorizeClientByMAC(context.Background(), "default", "aa:bb:cc:dd:ee:ff", "", "60")
	if err == nil {
		t.Fatal("expected error from rc:error response, got nil")
	}
}

// TestAuthorizeClientByMAC_MultipleClientsNotFound treats an ambiguous reply
// (more than one client echoed) as not-found, mirroring the other stamgr
// commands.
func TestAuthorizeClientByMAC_MultipleClientsNotFound(t *testing.T) {
	srv, _ := stamgrTestServer(t, "default", http.StatusOK,
		`{"meta":{"rc":"ok"},"data":[{"mac":"aa:aa:aa:aa:aa:aa"},{"mac":"bb:bb:bb:bb:bb:bb"}]}`)
	c := newStamgrTestClient(t, srv)

	err := c.AuthorizeClientByMAC(context.Background(), "default", "aa:bb:cc:dd:ee:ff", "", "")
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected NotFoundError for ambiguous reply, got: %v", err)
	}
}

// TestStamgrExistingCommandsUnchanged pins the pre-existing single-client
// contract of the other stamgr commands: an empty data array stays not-found
// for them (only authorize-guest is exempt).
func TestStamgrExistingCommandsUnchanged(t *testing.T) {
	srv, _ := stamgrTestServer(t, "default", http.StatusOK,
		`{"meta":{"rc":"ok"},"data":[]}`)
	c := newStamgrTestClient(t, srv)

	err := c.BlockClientByMAC(context.Background(), "default", "aa:bb:cc:dd:ee:ff")
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("BlockClientByMAC with empty data should stay NotFoundError, got: %v", err)
	}
}
