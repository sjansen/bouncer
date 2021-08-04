package server

import (
	"time"

	"github.com/go-chi/chi"
	cmw "github.com/go-chi/chi/middleware"

	"github.com/sjansen/bouncer/internal/web/handlers"
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

	requireLogin := s.saml.RequireAccount
	r.Method("GET", "/", requireLogin(
		handlers.NewRoot(),
	))
	r.Method("GET", "/b/whoami/", requireLogin(
		handlers.NewWhoAmI(),
	))
	r.Mount("/saml/", s.saml)
}
