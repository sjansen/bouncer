package handlers

import (
	"net/http"

	"github.com/crewjam/saml/samlsp"

	"github.com/sjansen/bouncer/internal/authz"
	"github.com/sjansen/bouncer/internal/web/config"
	"github.com/sjansen/bouncer/internal/web/pages"
)

// WhoAmI shows information about the current user.
type WhoAmI struct{}

// NewRoot creates a new root page handler.
func NewWhoAmI(cfg *config.Config) *WhoAmI {
	return &WhoAmI{}
}

// WhoAmI shows information about the current user.
func (p *WhoAmI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user *authz.User
	s := samlsp.SessionFromContext(r.Context())
	if u, ok := s.(*authz.User); ok {
		user = u
	}

	page := &pages.ProfilePage{}
	page.User = user
	pages.WriteResponse(w, page)
}
