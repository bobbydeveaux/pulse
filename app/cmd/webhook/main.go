// Package main — Pulse GitHub Marketplace billing webhook.
//
// Single-binary stdlib-only HTTP server that receives marketplace_purchase
// events when a customer upgrades, downgrades, or cancels Pulse on the
// GitHub Marketplace. Designed for Cloud Run via StackRamp.
//
// Routes:
//   POST /webhook/github-marketplace   verify X-Hub-Signature-256, log event
//   GET  /health                       liveness probe
//
// Env vars:
//   GITHUB_WEBHOOK_SECRET   HMAC-SHA256 shared secret configured in the
//                           Marketplace listing. Required — server returns
//                           500 on every event until set.
//   PORT                    Listen port. Default 8080 (Cloud Run convention).
//
// Mirrors guardian/app/cmd/webhook (agentops-014); the GitHub Marketplace
// event format is identical between Guardian and Pulse listings.
package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

// defaultPort is the Cloud Run convention when $PORT is not set.
const defaultPort = "8080"

// readHeaderTimeout caps how long the server will wait to read request
// headers. Slow-loris defence; 10s matches Guardian's webhook.
const readHeaderTimeout = 10 * time.Second

// resolvePort returns env if non-empty, else defaultPort. Extracted so the
// port-defaulting branch is testable without touching os.Getenv.
func resolvePort(env string) string {
	if env == "" {
		return defaultPort
	}
	return env
}

// newServer builds the fully-wired *http.Server. It does not start it.
// Extracted so the mux + timeout wiring is testable as a pure construction.
func newServer(port string, secret []byte) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhook/github-marketplace", handleMarketplace(secret))
	mux.HandleFunc("GET /health", handleHealth)
	return &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}
}

// serve calls srv.ListenAndServe and translates http.ErrServerClosed (the
// normal shutdown signal) to nil. Any other error is returned unchanged so
// the caller can decide whether to log.Fatal or recover.
func serve(srv *http.Server) error {
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func main() {
	port := resolvePort(os.Getenv("PORT"))
	secret := []byte(os.Getenv("GITHUB_WEBHOOK_SECRET"))
	if len(secret) == 0 {
		log.Printf("marketplace: GITHUB_WEBHOOK_SECRET is not set — webhook will 500 on every event until configured")
	}

	srv := newServer(port, secret)
	log.Printf("marketplace: listening on :%s", port)
	if err := serve(srv); err != nil {
		log.Fatalf("marketplace: server error: %v", err)
	}
}
