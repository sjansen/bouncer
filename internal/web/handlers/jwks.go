package handlers

import (
	"net/http"

	"github.com/sjansen/bouncer/internal/keyring"
)

// JWKS exposes the current JSON Web Key Set.
type JWKS struct {
	keyring *keyring.KeyRing
}

// NewJWKS creates a new handler.
func NewJWKS(keyring *keyring.KeyRing) *JWKS {
	return &JWKS{keyring: keyring}
}

// ServeHTTP handles reqeusts for the root page.
func (h *JWKS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(h.keyring.JWKSetAsJSON())
}
