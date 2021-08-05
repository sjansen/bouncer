package handlers

import (
	"net/http"

	"github.com/sjansen/bouncer/internal/keyring"
	"github.com/sjansen/bouncer/internal/web/config"
)

// JWKS exposes the current JSON Web Key Set.
type JWKS struct {
	keyring *keyring.KeyRing
}

// NewJWKS creates a new handler.
func NewJWKS(cfg *config.Config) *JWKS {
	return &JWKS{keyring: cfg.KeyRing}
}

// ServeHTTP handles HTTP requests.
func (h *JWKS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(h.keyring.JWKSetAsJSON())
}
