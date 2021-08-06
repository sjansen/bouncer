package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/sjansen/bouncer/internal/web/config"
)

// Redirect forwards authenticated requests to a target.
type Redirect struct{}

// NewRedirect creates a new handler.
func NewRedirect(cfg *config.Config) *Redirect {
	return &Redirect{}
}

// ServeHTTP handles HTTP requests.
func (h *Redirect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if !strings.HasPrefix(target, "/") {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("400 Bad Request"))
		return
	}

	url, err := url.Parse(target)
	if err != nil || url.Host != "" {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("403 Forbidden"))
		return
	}

	http.Redirect(w, r, target, http.StatusFound)
}
