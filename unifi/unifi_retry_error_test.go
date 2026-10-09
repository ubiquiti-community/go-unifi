package unifi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// A 5xx that outlives the retry budget must report the controller's status and
// message, not an opaque "giving up after N attempt(s)".
func TestDo_ExhaustedRetriesReportStatusAndBody(t *testing.T) {
	tests := []struct {
		name string
		body string
		want []string
	}{
		{
			name: "v2 JSON error",
			body: `{"code":"api.err.Boom","message":"rule is invalid"}`,
			want: []string{"api.err.Boom: rule is invalid", "500"},
		},
		{
			name: "non-JSON gateway page",
			body: `<html>upstream failure</html>`,
			want: []string{"upstream failure", "500"},
		},
		{
			name: "empty body",
			body: ``,
			want: []string{"empty response body", "500"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if handleNewStyleSetup(w, r) {
						return
					}
					if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
						w.Header().Set("X-Csrf-Token", "tok")
						w.WriteHeader(http.StatusOK)
						return
					}
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(tt.body))
				}),
			)
			t.Cleanup(srv.Close)

			retryMax := 1
			c, err := New(context.Background(), &Config{
				BaseURL:  srv.URL,
				Username: "admin",
				Password: "admin",
				RetryMax: &retryMax,
			})
			if err != nil {
				t.Fatalf("client init: %v", err)
			}
			c.c.RetryWaitMin = time.Millisecond
			c.c.RetryWaitMax = time.Millisecond

			_, err = c.CreateNat(context.Background(), "default", &Nat{Description: "x"})
			if err == nil {
				t.Fatal("expected an error")
			}
			if strings.Contains(err.Error(), "giving up") {
				t.Errorf("error is opaque: %v", err)
			}
			for _, want := range tt.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q is missing %q", err, want)
				}
			}
		})
	}
}
