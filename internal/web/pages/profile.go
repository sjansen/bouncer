package pages

import (
	"fmt"
	"io"
	"net/http"

	"github.com/crewjam/saml/samlsp"
	"github.com/sjansen/bouncer/internal/authz"
	"github.com/sjansen/bouncer/internal/web/config"
)

var _ Response = &ProfilePage{}

// ProfilePage shows information about a user.
type ProfilePage struct {
	Page
}

// WriteContent writes an HTTP response body.
func (p *ProfilePage) WriteContent(w io.Writer) {
	p.Title = "Profile"
	if err := tmpls.ExecuteTemplate(w, "profile.html", p); err != nil {
		fmt.Println(err)
	}
}

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

	page := &ProfilePage{}
	page.User = user
	WriteResponse(w, page)
}
