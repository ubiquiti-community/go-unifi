package unifi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

// TestDataRequest_RateLimitedThenSucceeds proves that with API-key auth (no
// login step at all) a controller 429 on an ordinary data request is retried
// by the transport — honoring Retry-After — instead of failing the call. A
// Terraform plan issues hundreds of serial reads through the cloud connector
// proxy, and a single rate-limit hit must not abort the run.
func TestDataRequest_RateLimitedThenSucceeds(t *testing.T) {
	var hits int32
	const rateLimitedAttempts = 2

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.URL.Path == "/proxy/network/api/s/default/rest/networkconf" {
			if atomic.AddInt32(&hits, 1) <= rateLimitedAttempts {
				w.Header().Set("Retry-After", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c, err := New(context.Background(), &Config{BaseURL: srv.URL, APIKey: "test-key"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if _, err := c.ListNetwork(context.Background(), "default"); err != nil {
		t.Fatalf("ListNetwork should survive %d rate-limited attempts, got: %v", rateLimitedAttempts, err)
	}
	if n := atomic.LoadInt32(&hits); n != rateLimitedAttempts+1 {
		t.Fatalf("expected %d attempts (2 rate-limited + 1 success), got %d", rateLimitedAttempts+1, n)
	}
}

// TestDataRequest_RateLimitExhaustionSurfacesError proves a persistent 429 on a
// data request still surfaces as a RateLimitError once the transport's retry
// budget is spent, rather than looping forever.
func TestDataRequest_RateLimitExhaustionSurfacesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	retryMax := 1
	c, err := New(context.Background(), &Config{BaseURL: srv.URL, APIKey: "test-key", RetryMax: &retryMax})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = c.ListNetwork(context.Background(), "default")
	var rle *RateLimitError
	if !errors.As(err, &rle) {
		t.Fatalf("expected RateLimitError after retry exhaustion, got: %v", err)
	}
}
