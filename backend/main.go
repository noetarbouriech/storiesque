package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/noetarbouriech/storiesque/backend/internal/api"
	"github.com/noetarbouriech/storiesque/backend/internal/auth"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/story"
	"github.com/noetarbouriech/storiesque/backend/internal/user"
)

func main() {
	// init database
	pg, err := db.NewPostgres("db", "5432", "postgres", "test")
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
		r.Use(authService.Verifier())
		r.Use(authService.Authenticator)

		r.Group(userService.UserRoutes)
		r.Group(storyService.UserRoutes)
	})

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", router)
}
