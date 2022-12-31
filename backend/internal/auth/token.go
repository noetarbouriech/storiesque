package auth

import (
	"net/http"
	"time"

	"github.com/noetarbouriech/storiesque/backend/internal/db"
)

func (s *Service) CreateAccessToken(userDB db.User) (*http.Cookie, error) {

	expireTime := time.Now().Add(1 * time.Second)
	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{
		"name":  userDB.Username,   // username
		"id":    userDB.ID,         // user id
		"admin": userDB.IsAdmin,    // user is_admin
		"iat":   time.Now(),        // issued time
		"exp":   expireTime.Unix(), // expire time
	})
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:       "jwt",
		Value:      tokenString,
		Path:       "/",
		Domain:     s.apiDomain,
		Expires:    expireTime,
		RawExpires: "",
		MaxAge:     10000,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "",
		Unparsed:   []string{},
	}, nil
}

func (s *Service) CreateRefreshToken() (*http.Cookie, error) {

	expireTime := time.Now().Add(1 * time.Minute)
	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{
		"iat": time.Now(),        // issued time
		"exp": expireTime.Unix(), // expire time
	})
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:       "refresh",
		Value:      tokenString,
		Path:       "/",
		Domain:     s.apiDomain,
		Expires:    expireTime,
		RawExpires: "",
		MaxAge:     10000,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "",
		Unparsed:   []string{},
	}, nil
}
