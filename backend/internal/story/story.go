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
	"github.com/go-playground/validator/v10"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

type Story struct {
	Id            int64  `json:"id"`
	Title         string `json:"title"         validate:"required,gte=0,lte=32"`
	Description   string `json:"description"   validate:"lte=512"`
	First_page_id int64  `json:"first_page_id"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (s *Service) PublicRoutes(r chi.Router) {
	r.Get("/story", s.getStories)    // get list of stories
	r.Get("/story/{id}", s.getStory) // get a single story

	r.Get("/page/{id}", s.getPage) // get a single page
}

func (s *Service) UserRoutes(r chi.Router) {
	r.Post("/story", s.createStory)        // create a story
	r.Delete("/story/{id}", s.deleteStory) // delete a story

	r.Post("/page/{id}", s.createPage)   // add a choice to a given page
	r.Put("/page/{id}", s.updatePage)    // update a page
	r.Delete("/page/{id}", s.deletePage) // delete a page (cascade delete choices)
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

	// translate json to struct
	errJson := json.NewDecoder(r.Body).Decode(&story)
	if errJson != nil {
		utils.Response(w, r, 500, "error while decoding json")
		return
	}

	// validate form
	err := validate.Struct(story)
	if err != nil {
		utils.Response(w, r, 400, "invalid input")
		return
	}

	// create story in db
	_, err = s.queries.CreateStory(context.Background(), db.CreateStoryParams{
		Title:       story.Title,
		Description: sql.NullString{String: story.Description, Valid: true},
	})
	if err != nil {
		utils.Response(w, r, 500, err.Error())
		log.Fatal(err.Error())
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
