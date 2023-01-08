package story

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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

type StoryCard struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
	HasImg      bool   `json:"has_img"`
	Featured    bool   `json:"is_featured"`
}

type StoryCreation struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"         validate:"required,gte=0,lte=48"`
	Description string `json:"description"   validate:"lte=512"`
	AuthorId    int64  `json:"author_id"`
}

type StoryDetails struct {
	Id            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	First_page_id int64  `json:"first_page_id"`
	AuthorName    string `json:"author_name"`
	HasImg        bool   `json:"has_img"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (s *Service) PublicRoutes(r chi.Router) {
	r.Get("/story", s.getStories)                  // get list of stories
	r.Get("/story/{id}", s.getStory)               // get a single story
	r.Get("/story/featured", s.getFeaturedStories) // get stories that are featured

	r.Get("/page/{id}", s.getPage) // get a single page
}

func (s *Service) UserRoutes(r chi.Router) {
	r.Post("/story", s.createStory)                // create a story
	r.Patch("/story/{id}", s.updateStory)          // update a story
	r.Patch("/story/{id}/featured", s.setFeatured) // set a story as featured on homepage
	r.Delete("/story/{id}", s.deleteStory)         // delete a story

	r.Post("/page/{id}", s.createPage)   // add a choice to a given page
	r.Patch("/page/{id}", s.updatePage)  // update a page
	r.Delete("/page/{id}", s.deletePage) // delete a page (cascade delete choices)
}

func (s *Service) getStories(w http.ResponseWriter, r *http.Request) {

	// get stories with filter by name
	title := r.URL.Query().Get("title")

	// get page number in query
	page := r.URL.Query().Get("page")
	if len(page) == 0 {
		page = "1"
	}

	stories, err := s.queries.SearchStories(context.Background(), db.SearchStoriesParams{
		Column1: sql.NullString{String: title, Valid: true},
		Column2: page,
	})
	if err != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}

	// map stories to json struct
	rStories := []StoryCard{}
	for _, story := range stories {
		rStory := StoryCard{
			Id:          story.ID,
			Title:       story.Title,
			Description: story.Description.String,
			AuthorName:  story.AuthorName,
			HasImg:      story.HasImg,
			Featured:    story.Featured,
		}
		rStories = append(rStories, rStory)
	}

	render.JSON(w, r, rStories)
}

func (s *Service) getStory(w http.ResponseWriter, r *http.Request) {

	// get id in url
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse story id")
		return
	}

	// get story from db
	story, err := s.queries.GetStory(context.Background(), int64(id))
	if err != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}
	storyJson := StoryDetails{
		Id:            story.ID,
		Title:         story.Title,
		Description:   story.Description.String,
		First_page_id: story.FirstPageID.Int64,
		AuthorName:    story.AuthorName,
		HasImg:        story.HasImg,
	}

	render.JSON(w, r, storyJson)
}

func (s *Service) createStory(w http.ResponseWriter, r *http.Request) {
	var story StoryCreation

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

	// get infos from jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		utils.Response(w, r, http.StatusUnauthorized, "not logged in")
		return
	}

	// create story in db
	storyDB, err := s.queries.CreateStory(context.Background(), db.CreateStoryParams{
		Title:       story.Title,
		Description: sql.NullString{String: story.Description, Valid: true},
		Author:      int64(claims["id"].(float64)),
	})
	if err != nil {
		utils.Response(w, r, 500, err.Error())
		return
	}

	// render story created in db
	render.JSON(w, r, StoryCreation{
		Id:          storyDB.ID,
		Title:       storyDB.Title,
		Description: storyDB.Description.String,
		AuthorId:    storyDB.Author,
	})
}

func (s *Service) updateStory(w http.ResponseWriter, r *http.Request) {

	// get id in url
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "story id bad format")
		return
	}

	// get story from db
	author, err := s.queries.GetStoryAuthor(context.Background(), int64(id))
	if err != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}

	// check if user is authorized
	if !utils.IsOwner(r, int(author)) {
		utils.Response(w, r, 401, "user is not the owner of the given resource")
		return
	}

	// map body to json
	var story StoryCreation
	err = json.NewDecoder(r.Body).Decode(&story)
	if err != nil {
		utils.Response(w, r, 500, "error while decoding json")
		return
	}

	// validate page form
	err = validate.Struct(story)
	if err != nil {
		utils.Response(w, r, 400, "invalid input")
		return
	}

	// update db
	err = s.queries.UpdateStory(context.Background(), db.UpdateStoryParams{
		ID: int64(id),

		TitleDoUpdate: story.Title != "",
		Title:         story.Title,

		DescriptionDoUpdate: story.Description != "",
		Description:         story.Description,
	})
	if err != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}

	utils.Response(w, r, 200, "story successfully updated")
}

func (s *Service) deleteStory(w http.ResponseWriter, r *http.Request) {

	// get id in url
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "story id bad format")
		return
	}

	// get story from db
	author, err := s.queries.GetStoryAuthor(context.Background(), int64(id))
	if err != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}

	// check if user is authorized
	if !utils.IsOwner(r, int(author)) {
		utils.Response(w, r, 401, "user is not the owner of the given resource")
		return
	}

	// delete story in db
	errDB := s.queries.DeleteStory(context.Background(), int64(id))
	if errDB != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}

	utils.Response(w, r, 200, "story successfully deleted")
}

func (s *Service) setFeatured(w http.ResponseWriter, r *http.Request) {

	// get id in param
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse story id")
		return
	}

	// check if user is authorized
	if !utils.IsAdmin(r) {
		utils.Response(w, r, 401, "user is not admin")
		return
	}

	// set story as featured in db
	errDB := s.queries.SetStoryAsFeatured(context.Background(), int64(id))
	if errDB != nil {
		fmt.Println(errDB.Error())
		utils.Response(w, r, 404, "story not found")
		return
	}

	utils.Response(w, r, 200, "story successfully updated")
}

func (s *Service) getFeaturedStories(w http.ResponseWriter, r *http.Request) {

	stories, err := s.queries.GetFeaturedStories(context.Background())
	if err != nil {
		utils.Response(w, r, 404, "story not found")
		return
	}

	// map stories to json struct
	rStories := []StoryCard{}
	for _, story := range stories {
		rStory := StoryCard{
			Id:          story.ID,
			Title:       story.Title,
			Description: story.Description.String,
			AuthorName:  story.AuthorName,
			HasImg:      story.HasImg,
			Featured:    true,
		}
		rStories = append(rStories, rStory)
	}

	render.JSON(w, r, rStories)
}
