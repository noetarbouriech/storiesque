package story

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Story struct {
	id    int64  `json:"id"`
	title string `json:"title"`
	// ...
}

func Routes(r chi.Router) {
	r.Get("/stories", getStories)
}

func getStories(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("stories"))
}
