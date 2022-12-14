package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/noetarbouriech/storiesque/backend/internal/api"
	"github.com/noetarbouriech/storiesque/backend/internal/auth"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/img"
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

	// init s3
	minio, err := img.NewMinio(
		os.Getenv("S3_HOST"),
		os.Getenv("S3_USER"),
		os.Getenv("S3_PASSWORD"),
	)

	// init router
	router := api.CreateRouter()

	// get brcrypt password cost from env
	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		log.Fatal("Error reading bcrypt cost")
	}

	// init services
	authService := auth.NewService(queries, os.Getenv("JWT_SECRET"), os.Getenv("API_DOMAIN"), cost)
	storyService := story.NewService(queries)
	userService := user.NewService(queries)
	imgService := img.NewService(queries, minio)

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
		r.Group(imgService.UserRoutes)
	})

	fmt.Println("Starting Storiesque api on " + os.Getenv("API_DOMAIN") + os.Getenv("API_PORT"))
	http.ListenAndServe(":"+os.Getenv("API_PORT"), router)
}
