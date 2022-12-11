package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/noetarbouriech/storiesque/internal/api"
	"github.com/noetarbouriech/storiesque/internal/db"
	"github.com/noetarbouriech/storiesque/internal/story"
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
	storyService := story.NewService(queries)
	router.Group(storyService.Routes)

	fmt.Println("Starting server on port 3000")
	http.ListenAndServe(":3000", router)
}
