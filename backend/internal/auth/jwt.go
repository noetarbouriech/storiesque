package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
)

// use the tokenAuth defined in the service for the Verifier middleware
func (s *Service) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(s.tokenAuth)
}

// same as jwtauth.Authenticator but in json
func (s *Service) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			// try to refresh token
			if s.refresh(w, r) {
				next.ServeHTTP(w, r)
			} else {
				utils.Response(w, r, http.StatusUnauthorized, "not logged in")
			}
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			utils.Response(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		// token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// if refresh token is still valid, refresh the access token and the refresh token
func (s *Service) refresh(w http.ResponseWriter, r *http.Request) bool {
	// verify refresh token
	_, err := jwtauth.VerifyRequest(s.tokenAuth, r, refreshFromCookie)
	if err != nil {
		return false
	}

	// get infos from expired jwt token
	_, claims, _ := jwtauth.FromContext(r.Context())

	// create new access token
	accessToken, err := s.CreateAccessToken(db.User{
		ID:           int64(claims["id"].(float64)),
		Username:     claims["name"].(string),
		PasswordHash: "",
		IsAdmin:      claims["admin"].(bool),
		Email:        "",
	})
	if err != nil {
		utils.Response(w, r, 500, "error with token creation")
		return false
	}
	http.SetCookie(w, accessToken)

	// create new refresh token
	refreshToken, err := s.CreateRefreshToken()
	if err != nil {
		utils.Response(w, r, 500, "error with token creation")
		return false
	}
	http.SetCookie(w, refreshToken)

	return true
}

// TokenFromCookie tries to retreive the token string from a cookie named
// "refresh".
func refreshFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("refresh")
	if err != nil {
		return ""
	}
	return cookie.Value
}
