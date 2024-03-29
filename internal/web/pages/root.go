package pages

import (
	"fmt"
	"io"
	"net/http"

	"github.com/crewjam/saml/samlsp"
	"github.com/sjansen/bouncer/internal/authz"
	"github.com/sjansen/bouncer/internal/build"
	"github.com/sjansen/bouncer/internal/web/config"
)

var _ Response = &RootPage{}

// RootPage is the default application starting poing.
type RootPage struct {
	Page

	GitSHA    string
	Timestamp string
}

// WriteContent writes an HTTP response body.
func (p *RootPage) WriteContent(w io.Writer) {
	if err := tmpls.ExecuteTemplate(w, "root.html", p); err != nil {
		fmt.Println(err)
	}
}

// Root is the default app starting page.
type Root struct{}

// NewRoot creates a new root page handler.
func NewRoot(cfg *config.Config) *Root {
	return &Root{}
}

// ServeHTTP handles requests for the root page.
func (p *Root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user *authz.User
	s := samlsp.SessionFromContext(r.Context())
	if u, ok := s.(*authz.User); ok {
		user = u
	}

	page := &RootPage{
		GitSHA:    build.GitSHA,
		Timestamp: build.Timestamp,
	}
	page.User = user
	WriteResponse(w, page)
}
