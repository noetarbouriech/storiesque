package story

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

type Story struct {
	Id            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	First_page_id int32  `json:"first_page_id"`
}

func (s *Service) Routes(r chi.Router) {
	r.Get("/story", s.getStories)
	r.Post("/story", s.createStory)
	r.Delete("/story/{id}", s.deleteStory)
}

func (s *Service) getStories(w http.ResponseWriter, r *http.Request) {
	stories, err := s.queries.ListStories(context.Background())
	rStories := []Story{}
	for _, story := range stories {
		rStory := Story{
			Id:            story.ID,
			Title:         story.Title,
			Description:   story.Description.String,
			First_page_id: story.FirstPageID.Int32,
		}
		rStories = append(rStories, rStory)
	}
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	render.JSON(w, r, rStories)
}

func (s *Service) createStory(w http.ResponseWriter, r *http.Request) {
	var story Story
	errJson := json.NewDecoder(r.Body).Decode(&story)
	if errJson != nil {
		render.JSON(w, r, map[string]string{"message": "issue with json decoding"})
		return
	}

	if len(story.Title) == 0 || len(story.Title) > 32 {
		render.JSON(w, r, map[string]string{"message": "issue with title length"})
		return
	}
	if len(story.Description) > 512 {
		render.JSON(w, r, map[string]string{"message": "issue with description"})
		return
	}

	_, errDB := s.queries.CreateStory(context.Background(), db.CreateStoryParams{
		Title:       story.Title,
		Description: sql.NullString{String: story.Description, Valid: true},
	})
	if errDB != nil {
		log.Fatal(errDB.Error())
		return
	}

	render.JSON(w, r, map[string]string{"message": "successfully created"})
}

func (s *Service) deleteStory(w http.ResponseWriter, r *http.Request) {
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		render.JSON(w, r, map[string]string{"message": "Impossible to parse int"})
		return
	}
	errDB := s.queries.DeleteStory(context.Background(), int64(id))
	if errDB != nil {
		render.JSON(w, r, map[string]string{"message": "No story found"})
		return
	}
	render.JSON(w, r, map[string]string{"message": "successfully deleted"})
}
