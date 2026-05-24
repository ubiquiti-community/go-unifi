package unifi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// handleNewStyleSetup answers the unauthenticated probe (GET /) and the
// best-effort version lookup so New() can reach the login step against a
// new-style (UniFi OS) controller. Returns true if it handled the request.
func handleNewStyleSetup(w http.ResponseWriter, r *http.Request) bool {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/":
		w.WriteHeader(http.StatusOK) // 200 => new-style API
		return true
	case r.URL.Path == "/proxy/network/status":
		_, _ = w.Write([]byte(`{"meta":{"server_version":"8.0.0"}}`))
		return true
	}
	return false
}

// shortLoginBackoff shrinks the inter-attempt backoff so rate-limit tests run
// fast, restoring the production default afterwards. Honored Retry-After waits
// are unaffected (they come from the response, not loginRetryBackoff).
func shortLoginBackoff(t *testing.T) {
	t.Helper()
	prev := loginRetryBackoff
	loginRetryBackoff = 5 * time.Millisecond
	t.Cleanup(func() { loginRetryBackoff = prev })
}

// TestLogin_RateLimitedThenSucceeds proves a normal workflow survives the
// controller's login rate-limit: the controller returns several HTTP 429s before
// accepting the login, and New() must retry until it succeeds rather than failing
// with a confusing "unable to login" error.
func TestLogin_RateLimitedThenSucceeds(t *testing.T) {
	shortLoginBackoff(t)

	var loginHits int32
	const rateLimitedAttempts = 5

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			n := atomic.AddInt32(&loginHits, 1)
			if n <= rateLimitedAttempts {
				w.Header().Set("Retry-After", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		t.Fatalf("expected login to survive %d rate-limited attempts, got error: %v", rateLimitedAttempts, err)
	}
	if got := atomic.LoadInt32(&loginHits); got != rateLimitedAttempts+1 {
		t.Errorf("expected %d login attempts, got %d", rateLimitedAttempts+1, got)
	}
}

// TestLogin_HonorsRetryAfter verifies the controller's Retry-After hint is waited
// out (rather than retried immediately).
func TestLogin_HonorsRetryAfter(t *testing.T) {
	shortLoginBackoff(t)

	var loginHits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			if atomic.AddInt32(&loginHits, 1) == 1 {
				w.Header().Set("Retry-After", "1") // 1 second
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	start := time.Now()
	_, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed < 900*time.Millisecond {
		t.Errorf("expected to honor Retry-After ~1s, but login took only %s", elapsed)
	}
}

// TestLogin_ExhaustionReturnsRateLimitError verifies that a persistent rate-limit
// surfaces a clear, typed RateLimitError instead of an opaque "giving up" message.
func TestLogin_ExhaustionReturnsRateLimitError(t *testing.T) {
	shortLoginBackoff(t)
	prevMax := loginRetryMax
	loginRetryMax = 3
	t.Cleanup(func() { loginRetryMax = prevMax })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	if err == nil {
		t.Fatal("expected an error when the controller never clears the rate-limit")
	}
	var rle *RateLimitError
	if !errors.As(err, &rle) {
		t.Errorf("expected error to unwrap to *RateLimitError, got: %v", err)
	}
}

// TestLogin_EmptyBodyResponseRetries verifies that a 2xx response that does not
// establish a session (e.g. an empty body under throttling) is retried rather
// than accepted as a broken "logged-in" state.
func TestLogin_EmptyBodyResponseRetries(t *testing.T) {
	shortLoginBackoff(t)

	var loginHits int32
	const emptyAttempts = 2
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			n := atomic.AddInt32(&loginHits, 1)
			if n <= emptyAttempts {
				w.WriteHeader(http.StatusOK) // 200 but no CSRF token => no session
				return
			}
			w.Header().Set("X-Csrf-Token", "tok")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		t.Fatalf("expected retry past empty-body responses, got error: %v", err)
	}
	if got := atomic.LoadInt32(&loginHits); got != emptyAttempts+1 {
		t.Errorf("expected %d login attempts, got %d", emptyAttempts+1, got)
	}
}

// TestLogin_APIKeySkipsLogin verifies API-key auth issues no POST /api/auth/login
// at all (and sends the key as a header), which is the definitive way to avoid the
// login rate-limit across separate Terraform invocations.
func TestLogin_APIKeySkipsLogin(t *testing.T) {
	var loginHits int32
	var sawAPIKey atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != "" {
			sawAPIKey.Store(true)
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			atomic.AddInt32(&loginHits, 1)
		}
		if r.URL.Path == "/proxy/network/status" {
			_, _ = w.Write([]byte(`{"meta":{"server_version":"8.0.0"}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := New(context.Background(), &Config{
		BaseURL: srv.URL,
		APIKey:  "secret-key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := atomic.LoadInt32(&loginHits); got != 0 {
		t.Errorf("expected 0 login attempts with API key, got %d", got)
	}
	if !sawAPIKey.Load() {
		t.Error("expected requests to carry the X-Api-Key header")
	}
}
