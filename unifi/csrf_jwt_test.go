package unifi

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// jwtWith builds a fake JWT ("hdr.<payload>.sig") whose payload carries the
// given csrfToken claim.
func jwtWith(csrf string) string {
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"csrfToken":"` + csrf + `"}`))
	return "hdr." + payload + ".sig"
}

// jwtWithExp is jwtWith plus a standard exp claim (epoch seconds).
func jwtWithExp(csrf string, exp int64) string {
	payload := base64.RawURLEncoding.EncodeToString(
		[]byte(fmt.Sprintf(`{"csrfToken":%q,"exp":%d}`, csrf, exp)),
	)
	return "hdr." + payload + ".sig"
}

// TestLogin_CSRFFromJWTCookieWhenHeaderAbsent proves login succeeds against
// UniFi OS builds that return the CSRF token only inside the TOKEN cookie's
// JWT and omit the X-Csrf-Token response header. The client must recover the
// token from the cookie (isLoggedIn requires a non-empty CSRF token on
// new-style controllers, so without recovery login fails outright) and send
// it on subsequent requests.
func TestLogin_CSRFFromJWTCookieWhenHeaderAbsent(t *testing.T) {
	shortLoginBackoff(t)

	var lastSeenCSRF string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			// Deliberately no X-Csrf-Token header: the token travels only in
			// the JWT cookie, as some UniFi OS builds do.
			http.SetCookie(w, &http.Cookie{Name: "TOKEN", Value: jwtWith("cookie-csrf")})
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path == "/proxy/network/api/s/default/self" {
			lastSeenCSRF = r.Header.Get("X-Csrf-Token")
			_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		t.Fatalf("login should succeed via JWT-cookie CSRF recovery, got: %v", err)
	}

	if err := c.do(context.Background(), http.MethodGet, "api/s/default/self", nil, nil); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if lastSeenCSRF != "cookie-csrf" {
		t.Errorf("expected X-Csrf-Token %q on follow-up request, got %q", "cookie-csrf", lastSeenCSRF)
	}
}

// TestLogin_CSRFHeaderTakesPrecedenceOverCookie verifies the existing header
// behavior is unchanged: when a response carries both an X-Csrf-Token header
// and a JWT cookie, the header wins.
func TestLogin_CSRFHeaderTakesPrecedenceOverCookie(t *testing.T) {
	var lastSeenCSRF string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			http.SetCookie(w, &http.Cookie{Name: "TOKEN", Value: jwtWith("cookie-csrf")})
			w.Header().Set("X-Csrf-Token", "header-csrf")
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path == "/proxy/network/api/s/default/self" {
			lastSeenCSRF = r.Header.Get("X-Csrf-Token")
			_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}

	if err := c.do(context.Background(), http.MethodGet, "api/s/default/self", nil, nil); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if lastSeenCSRF != "header-csrf" {
		t.Errorf("expected header token to win, got X-Csrf-Token %q", lastSeenCSRF)
	}
}

// TestCSRFFromJWTCookie exercises the claim extraction against well-formed and
// malformed cookie values.
func TestCSRFFromJWTCookie(t *testing.T) {
	noClaim := "hdr." + base64.RawURLEncoding.EncodeToString([]byte(`{"userId":"u"}`)) + ".sig"

	cases := []struct {
		name    string
		cookies []*http.Cookie
		want    string
		wantExp time.Time
	}{
		{"TOKEN cookie with claim", []*http.Cookie{{Name: "TOKEN", Value: jwtWith("abc")}}, "abc", time.Time{}},
		{"UOS_TOKEN cookie with claim", []*http.Cookie{{Name: "UOS_TOKEN", Value: jwtWith("xyz")}}, "xyz", time.Time{}},
		{"claim with exp", []*http.Cookie{{Name: "TOKEN", Value: jwtWithExp("abc", 1786000000)}}, "abc", time.Unix(1786000000, 0)},
		{"unrelated cookie name ignored", []*http.Cookie{{Name: "unifises", Value: jwtWith("abc")}}, "", time.Time{}},
		{"not a JWT", []*http.Cookie{{Name: "TOKEN", Value: "opaque-session-id"}}, "", time.Time{}},
		{"bad base64 payload", []*http.Cookie{{Name: "TOKEN", Value: "hdr.!!!.sig"}}, "", time.Time{}},
		{"no csrfToken claim", []*http.Cookie{{Name: "TOKEN", Value: noClaim}}, "", time.Time{}},
		{"no cookies", nil, "", time.Time{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, gotExp := csrfFromJWTCookie(tc.cookies)
			if got != tc.want {
				t.Errorf("csrfFromJWTCookie() token = %q, want %q", got, tc.want)
			}
			if !gotExp.Equal(tc.wantExp) {
				t.Errorf("csrfFromJWTCookie() exp = %v, want %v", gotExp, tc.wantExp)
			}
		})
	}
}

// TestRelogin_PurgesStaleUOSTokenCookie verifies that a re-login does not send
// a stale UOS_TOKEN session cookie: controllers that name their session cookie
// UOS_TOKEN error out when a login request carries the expired token.
func TestRelogin_PurgesStaleUOSTokenCookie(t *testing.T) {
	shortLoginBackoff(t)

	var loginCookies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			names := make([]string, 0, len(r.Cookies()))
			for _, ck := range r.Cookies() {
				names = append(names, ck.Name)
			}
			loginCookies = append(loginCookies, strings.Join(names, ","))
			// Path=/ as real UniFi OS controllers set it; the purge relies on it.
			http.SetCookie(w, &http.Cookie{Name: "UOS_TOKEN", Value: jwtWith("cookie-csrf"), Path: "/"})
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}

	// Force a re-login with the previous session's UOS_TOKEN still in the jar.
	if err := c.login(context.Background()); err != nil {
		t.Fatalf("unexpected re-login error: %v", err)
	}

	if len(loginCookies) != 2 {
		t.Fatalf("expected 2 login requests, saw %d", len(loginCookies))
	}
	if strings.Contains(loginCookies[1], "UOS_TOKEN") {
		t.Errorf("re-login request still carried the stale UOS_TOKEN cookie (cookies: %q)", loginCookies[1])
	}
}

// TestLogin_JWTCookieExpTracksTokenExpiry verifies that when the CSRF token is
// recovered from the JWT cookie and the response has no X-Token-Expire-Time
// header, the JWT's exp claim drives session-expiry tracking.
func TestLogin_JWTCookieExpTracksTokenExpiry(t *testing.T) {
	shortLoginBackoff(t)

	exp := time.Now().Add(2 * time.Hour).Unix()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			http.SetCookie(w, &http.Cookie{Name: "TOKEN", Value: jwtWithExp("cookie-csrf", exp)})
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}

	if !c.tokenExpiry.Equal(time.Unix(exp, 0)) {
		t.Errorf("tokenExpiry = %v, want %v (from the JWT exp claim)", c.tokenExpiry, time.Unix(exp, 0))
	}
}

// TestConcurrentRequestsCSRFTokenRace drives parallel requests against a
// controller that re-issues the TOKEN cookie on every response, so the CSRF
// capture path runs concurrently under the loginMu read lock. Run with -race.
func TestConcurrentRequestsCSRFTokenRace(t *testing.T) {
	shortLoginBackoff(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handleNewStyleSetup(w, r) {
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == loginPathNew {
			http.SetCookie(w, &http.Cookie{Name: "TOKEN", Value: jwtWith("csrf-0")})
			w.WriteHeader(http.StatusOK)
			return
		}
		// Re-issue the session cookie on every data response, as UniFi OS does.
		http.SetCookie(w, &http.Cookie{Name: "TOKEN", Value: jwtWith("csrf-" + r.URL.Query().Get("i"))})
		_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
	}))
	defer srv.Close()

	c, err := New(context.Background(), &Config{
		BaseURL:  srv.URL,
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}

	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = c.do(context.Background(), http.MethodGet, "api/s/default/self", nil, nil,
				map[string]string{"i": fmt.Sprint(i)})
		}()
	}
	wg.Wait()
}
