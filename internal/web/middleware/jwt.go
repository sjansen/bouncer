package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sjansen/bouncer/internal/keyring"
)

type JWT struct {
	KeyRing *keyring.KeyRing
	Secure  bool
	Subject string
}

func (m *JWT) SetJWT(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt, err := m.KeyRing.NewJWT(m.Subject)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		cookie := &http.Cookie{
			Name:     "auth_token",
			Value:    string(jwt),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   1 * 60 * 60,
			// Browsers ignore Strict cookies set on redirect.
			SameSite: http.SameSiteLaxMode,
			Secure:   m.Secure,
		}
		http.SetCookie(w, cookie)
		handler.ServeHTTP(w, r)
	})
}
