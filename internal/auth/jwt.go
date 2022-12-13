package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

// same as jwtauth.Authenticator but in json
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"message": err.Error()})
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
