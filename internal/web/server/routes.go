package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	cmw "github.com/go-chi/chi/middleware"

	"github.com/sjansen/bouncer/internal/web/handlers"
	"github.com/sjansen/bouncer/internal/web/images"
	"github.com/sjansen/bouncer/internal/web/middleware"
	"github.com/sjansen/bouncer/internal/web/pages"
)

func (s *Server) addRoutes() {
	jwt := &middleware.JWT{
		KeyRing: s.config.KeyRing,
		Secure:  !s.config.Insecure,
		Subject: s.config.AppURL.Host,
	}
	requireLogin := chi.Chain(
		s.saml.RequireAccount,
		jwt.SetJWT,
	).Handler

	r := chi.NewRouter()
	s.router = r

	r.Use(
		cmw.RequestID,
		cmw.RealIP,
		cmw.Logger,
		cmw.Recoverer,
		cmw.Timeout(5*time.Second),
		cmw.Heartbeat("/ping"),
		s.sess.LoadAndSave,
		s.relaystate.LoadAndSave,
	)

	r.Route("/b", func(r chi.Router) {
		r.Get("/jwks/", handlers.NewJWKS(s.config).ServeHTTP)
		r.Mount("/saml/", s.saml)
		r.Route("/", func(r chi.Router) {
			r.Use(requireLogin)
			r.Get("/", pages.NewRoot(s.config).ServeHTTP)
			r.Get("/redirect/", handlers.NewRedirect(s.config).ServeHTTP)
			r.Get("/whoami/", pages.NewWhoAmI(s.config).ServeHTTP)
			r.Method("GET", "/*",
				http.StripPrefix("/b/", images.NewHandler()),
			)
		})
	})
}
