package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/noetarbouriech/storiesque/backend/internal/api"
	"github.com/noetarbouriech/storiesque/backend/internal/auth"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/story"
	"github.com/noetarbouriech/storiesque/backend/internal/user"
)

func main() {
	// environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// init database
	pg, err := db.NewPostgres(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	queries := db.New(pg.DB)

	// init router
	router := api.CreateRouter()

	// init services
	authService := auth.NewService(queries, os.Getenv("JWT_SECRET"), os.Getenv("API_DOMAIN"))
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

	fmt.Println("Starting Storiesque api on " + os.Getenv("API_PORT"))
	http.ListenAndServe(":"+os.Getenv("API_PORT"), router)
}
