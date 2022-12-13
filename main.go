package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/noetarbouriech/storiesque/internal/api"
	"github.com/noetarbouriech/storiesque/internal/auth"
	"github.com/noetarbouriech/storiesque/internal/db"
	"github.com/noetarbouriech/storiesque/internal/story"
	"github.com/noetarbouriech/storiesque/internal/user"
)

func main() {
	// init database
	pg, err := db.NewPostgres("localhost", "5431", "postgres", "test")
	if err != nil {
		log.Fatal(err.Error())
	}
	queries := db.New(pg.DB)

	// init router
	router := api.CreateRouter()

	// init services
	authService := auth.NewService(queries, "temp_secret")
	storyService := story.NewService(queries)
	userService := user.NewService(queries)

	// Public routes
	router.Group(func(r chi.Router) {
		r.Group(authService.PublicRoutes)
		r.Group(userService.PublicRoutes)
		r.Group(storyService.PublicRoutes)
	})

	// User routes
	router.Group(func(r chi.Router) {
		r.Use(auth.Verifier())
		r.Use(auth.Authenticator)

		r.Group(userService.UserRoutes)
		r.Group(storyService.UserRoutes)
	})

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", router)
}
