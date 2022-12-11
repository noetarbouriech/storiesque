package story

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/noetarbouriech/storiesque/internal/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) Routes(r chi.Router) {
	r.Get("/story", s.getStories)
	r.Post("/story", s.createStory)
}

func (s *Service) getStories(w http.ResponseWriter, r *http.Request) {
	stories, err := s.queries.ListStories(context.Background())
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	render.JSON(w, r, stories)
}

func (s *Service) createStory(w http.ResponseWriter, r *http.Request) {
	var story db.Story
	json.NewDecoder(r.Body).Decode(&story)

	if len(story.Title) == 0 || len(story.Title) > 32 {
		render.JSON(w, r, map[string]string{"message": "issue with title length"})
		return
	}
	if len(story.Description) > 512 {
		render.JSON(w, r, map[string]string{"message": "description too long"})
		return
	}

	_, err := s.queries.CreateStory(context.Background(), db.CreateStoryParams{
		Title:       story.Title,
		Description: story.Description,
	})
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	render.JSON(w, r, map[string]string{"message": "successfully created"})
}
