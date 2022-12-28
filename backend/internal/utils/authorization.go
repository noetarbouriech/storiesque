package utils

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

// check if a logged user is the owner of the resource he wants to use
func IsOwner(r *http.Request, id int) bool {

	// get infos from jwt
	// cannot be error since the jwt is verified for userRouters
	_, claims, _ := jwtauth.FromContext(r.Context())

	// return true if admin or owner
	return claims["admin"].(bool) || int64(claims["id"].(float64)) == int64(id)
}
