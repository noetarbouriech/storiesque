package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/noetarbouriech/storiesque/internal/db"
	"github.com/noetarbouriech/storiesque/internal/story"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	dbInstance, err := sql.Open("postgres", "port=5431 user=postgres password=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(dbInstance)
	if err != nil {
		log.Fatal(err.Error())
	}

	insertedStory, err := queries.CreateStory(ctx, "Test Story")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(insertedStory)

	stories, err := queries.ListStories(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(stories)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.CleanPath)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Set-Cookie"},
		AllowCredentials:   true,
		MaxAge:             300,
		OptionsPassthrough: false,
		Debug:              false,
	}))

	r.Group(story.Routes)

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", r)
}
