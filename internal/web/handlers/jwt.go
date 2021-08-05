package handlers

import (
	"net/http"

	"github.com/sjansen/bouncer/internal/keyring"
	"github.com/sjansen/bouncer/internal/web/config"
)

// JWT returns a fresh JSON Web Token.
type JWT struct {
	keyring *keyring.KeyRing
	subject string
}

// NewJWT creates a new handler.
func NewJWT(cfg *config.Config) *JWT {
	return &JWT{
		keyring: cfg.KeyRing,
		subject: cfg.AppURL.String(),
	}
}

// ServeHTTP handles HTTP requests.
func (h *JWT) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwt, err := h.keyring.NewJWT(h.subject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write(jwt)
	}
}
