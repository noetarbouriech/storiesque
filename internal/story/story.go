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
	"github.com/noetarbouriech/storiesque/internal/utils"
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
	First_page_id int64  `json:"first_page_id"`
}

func (s *Service) PublicRoutes(r chi.Router) {
	r.Get("/story", s.getStories)
	r.Get("/story/{id}", s.getStory)
}

func (s *Service) UserRoutes(r chi.Router) {
	r.Post("/story", s.createStory)
	r.Delete("/story/{id}", s.deleteStory)
}

func (s *Service) getStories(w http.ResponseWriter, r *http.Request) {
	stories, err := s.queries.ListStories(context.Background())
	if err != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}
	rStories := []Story{}
	for _, story := range stories {
		rStory := Story{
			Id:            story.ID,
			Title:         story.Title,
			Description:   story.Description.String,
			First_page_id: story.FirstPageID.Int64,
		}
		rStories = append(rStories, rStory)
	}
	render.JSON(w, r, rStories)
}

func (s *Service) getStory(w http.ResponseWriter, r *http.Request) {
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse story id")
		return
	}
	story, err := s.queries.GetStory(context.Background(), int64(id))
	if err != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}
	storyJson := Story{
		Id:            story.ID,
		Title:         story.Title,
		Description:   story.Description.String,
		First_page_id: story.FirstPageID.Int64,
	}
	render.JSON(w, r, storyJson)
}

func (s *Service) createStory(w http.ResponseWriter, r *http.Request) {
	var story Story
	errJson := json.NewDecoder(r.Body).Decode(&story)
	if errJson != nil {
		utils.Response(w, r, 500, "error while decoding json")
		return
	}

	if len(story.Title) == 0 || len(story.Title) > 32 {
		utils.Response(w, r, 400, "title too short or too long")
		return
	}
	if len(story.Description) > 512 {
		utils.Response(w, r, 400, "description too long")
		return
	}

	_, errDB := s.queries.CreateStory(context.Background(), db.CreateStoryParams{
		Title:       story.Title,
		Description: sql.NullString{String: story.Description, Valid: true},
	})
	if errDB != nil {
		utils.Response(w, r, 500, errDB.Error())
		log.Fatal(errDB.Error())
		return
	}

	utils.Response(w, r, 201, "story successfully created")
}

func (s *Service) deleteStory(w http.ResponseWriter, r *http.Request) {
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "story id bad format")
		return
	}
	errDB := s.queries.DeleteStory(context.Background(), int64(id))
	if errDB != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}
	utils.Response(w, r, 200, "story successfully deleted")
}
