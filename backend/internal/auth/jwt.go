package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
)

func (s *Service) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(s.tokenAuth)
}

// same as jwtauth.Authenticator but in json
func (s *Service) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			utils.Response(w, r, http.StatusUnauthorized, "not logged in")
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			utils.Response(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
