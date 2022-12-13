package utils

import (
	"net/http"

	"github.com/go-chi/render"
)

// unified json http responses
func Response(w http.ResponseWriter, r *http.Request, code int, msg string) {
	render.Status(r, code)
	render.JSON(w, r, map[string]string{"message": msg})
}
