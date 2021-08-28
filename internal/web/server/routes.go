package server

import (
	"time"

	"github.com/go-chi/chi"
	cmw "github.com/go-chi/chi/middleware"

	"github.com/sjansen/bouncer/internal/web/handlers"
	"github.com/sjansen/bouncer/internal/web/middleware"
)

func (s *Server) addRoutes() {
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
	r.Mount("/b/saml/", s.saml)

	jwt := &middleware.JWT{
		KeyRing: s.config.KeyRing,
		Secure:  !s.config.Insecure,
		Subject: s.config.AppURL.Host,
	}
	requireLogin := chi.Chain(
		s.saml.RequireAccount,
		jwt.SetJWT,
	).Handler

	r.Method("GET", "/b/", requireLogin(
		handlers.NewRoot(s.config),
	))
	r.Method("GET", "/b/jwks/", handlers.NewJWKS(s.config))
	r.Method("GET", "/b/redirect/", requireLogin(
		handlers.NewRedirect(s.config),
	))
	r.Method("GET", "/b/whoami/", requireLogin(
		handlers.NewWhoAmI(s.config),
	))
}
