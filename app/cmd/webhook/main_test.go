package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// captureLog redirects log.Default's output to a bytes.Buffer for the
// duration of fn, then restores it. Used to assert on boot-time log
// lines emitted by logStartup.
func captureLog(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(prev) })
	fn()
	return buf.String()
}

// --- resolvePort ----------------------------------------------------------

func TestResolvePortDefaultsTo8080(t *testing.T) {
	if got := resolvePort(""); got != defaultPort {
		t.Fatalf("resolvePort(\"\") = %q, want %q", got, defaultPort)
	}
}

func TestResolvePortUsesEnvWhenSet(t *testing.T) {
	if got := resolvePort("9090"); got != "9090" {
		t.Fatalf("resolvePort(\"9090\") = %q, want %q", got, "9090")
	}
}

// --- newServer ------------------------------------------------------------

func TestNewServerWiring(t *testing.T) {
	srv := newServer("8080", []byte("secret"))
	if srv == nil {
		t.Fatal("newServer returned nil")
	}
	if srv.Addr != ":8080" {
		t.Errorf("Addr = %q, want %q", srv.Addr, ":8080")
	}
	if srv.ReadHeaderTimeout != readHeaderTimeout {
		t.Errorf("ReadHeaderTimeout = %v, want %v", srv.ReadHeaderTimeout, readHeaderTimeout)
	}
	if srv.Handler == nil {
		t.Fatal("Handler is nil — newServer should attach a mux")
	}
}

func TestNewServerUsesProvidedPort(t *testing.T) {
	srv := newServer("9999", []byte("secret"))
	if srv.Addr != ":9999" {
		t.Errorf("Addr = %q, want %q", srv.Addr, ":9999")
	}
}

func TestNewServerRoutesHealth(t *testing.T) {
	// The handler returned by newServer should route GET /health to
	// handleHealth and respond 200 with the status:ok body.
	srv := newServer("8080", []byte(testSecret))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type=%q, want application/json", got)
	}
	if body := rec.Body.String(); !strings.Contains(body, `"status":"ok"`) {
		t.Errorf("body=%q, want contains status:ok", body)
	}
}

func TestNewServerRoutesMarketplaceWithSignedBody(t *testing.T) {
	// The handler returned by newServer should route POST
	// /webhook/github-marketplace through handleMarketplace, which means a
	// correctly-signed `purchased` payload should round-trip 200.
	secret := []byte(testSecret)
	srv := newServer("8080", secret)
	body := []byte(`{"action":"purchased","marketplace_purchase":{"account":{"login":"acme"},"plan":{"name":"p"}}}`)
	req := httptest.NewRequest(http.MethodPost, "/webhook/github-marketplace", bytes.NewReader(body))
	req.Header.Set("X-Hub-Signature-256", sign(secret, body))
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /webhook status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if got, _ := io.ReadAll(rec.Body); !strings.Contains(string(got), `"ok":true`) {
		t.Errorf("body=%q, want contains ok:true", got)
	}
}

func TestNewServerRejectsWrongMethodOnWebhook(t *testing.T) {
	// Go 1.22+ method-prefixed mux returns 405 on method mismatch. This
	// proves the route is registered as POST-only (not as a catch-all).
	srv := newServer("8080", []byte(testSecret))
	req := httptest.NewRequest(http.MethodGet, "/webhook/github-marketplace", nil)
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET /webhook status = %d, want 405", rec.Code)
	}
}

// --- serve ----------------------------------------------------------------

func TestServeReturnsNilOnGracefulClose(t *testing.T) {
	// Pre-close the server so ListenAndServe returns ErrServerClosed.
	// serve() translates that to nil — graceful shutdown is not an error.
	srv := newServer("0", []byte(testSecret))
	if err := srv.Close(); err != nil {
		t.Fatalf("pre-close: %v", err)
	}
	if err := serve(srv); err != nil {
		t.Fatalf("serve() = %v, want nil for graceful close", err)
	}
}

// --- logStartup -----------------------------------------------------------

func TestLogStartupWarnsOnMissingSecret(t *testing.T) {
	out := captureLog(t, func() {
		logStartup(nil, "8080")
	})
	if !strings.Contains(out, "GITHUB_WEBHOOK_SECRET is not set") {
		t.Errorf("log = %q, want missing-secret warning", out)
	}
	if !strings.Contains(out, "listening on :8080") {
		t.Errorf("log = %q, want listening line", out)
	}
}

func TestLogStartupSilentOnSecretPresent(t *testing.T) {
	out := captureLog(t, func() {
		logStartup([]byte("present"), "9090")
	})
	if strings.Contains(out, "GITHUB_WEBHOOK_SECRET is not set") {
		t.Errorf("log = %q, should not warn when secret is set", out)
	}
	if !strings.Contains(out, "listening on :9090") {
		t.Errorf("log = %q, want listening line", out)
	}
}

func TestServeReturnsErrOnBindFailure(t *testing.T) {
	// Bind a listener to grab a port, then ask serve() to bind to that
	// same port. ListenAndServe should fail with "address already in use"
	// (or similar). serve() must surface that error to the caller.
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("pre-listen: %v", err)
	}
	defer lis.Close()

	srv := &http.Server{
		Addr:              lis.Addr().String(),
		Handler:           http.NewServeMux(),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	done := make(chan error, 1)
	go func() { done <- serve(srv) }()

	select {
	case err := <-done:
		if err == nil {
			t.Fatalf("serve() = nil, want bind-failure error")
		}
	case <-time.After(2 * time.Second):
		_ = srv.Close()
		<-done
		t.Fatal("serve() did not return within 2s on bind failure")
	}
}
