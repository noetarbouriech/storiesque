package story

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
)

type Page struct {
	Id      int64   `json:"id"`
	Title   string  `json:"title"             validate:"gte=0,lte=32"`
	Body    string  `json:"body"              validate:"gte=0,lte=4096"`
	Choices []int64 `json:"choices,omitempty"`
}

func (s *Service) getPage(w http.ResponseWriter, r *http.Request) {
	// get id in param
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse page id")
		return
	}

	// get page in db
	page, err := s.queries.GetPage(context.Background(), int64(id))
	if err != nil {
		utils.Response(w, r, 404, "page not found")
		return
	}

	// get choices for given page in db
	choices, err := s.queries.ListChoices(context.Background(), int64(id))
	if err != nil {
		utils.Response(w, r, 404, "choices not found")
		return
	}

	pageJson := Page{
		Id:      page.ID,
		Title:   page.Title,
		Body:    page.Body,
		Choices: choices,
	}
	render.JSON(w, r, pageJson)
}

func (s *Service) createPage(w http.ResponseWriter, r *http.Request) {
	// get id in param
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse page id")
		return
	}

	// check if page in param exists
	exists, _ := s.queries.GetPage(context.Background(), int64(id))
	if exists.ID == 0 {
		utils.Response(w, r, 404, fmt.Sprintf("page with id %d doesn't exist", id))
		return
	}

	// create page in db
	page, err := s.queries.CreatePage(context.Background(), db.CreatePageParams{
		Title: "New Page",
		Body:  "Hello world !",
	})
	if err != nil {
		utils.Response(w, r, 500, err.Error())
		log.Fatal(err.Error())
		return
	}

	// create choice in db
	_, err = s.queries.CreateChoices(context.Background(), db.CreateChoicesParams{
		PageID: int64(id),
		PathID: page.ID,
	})
	if err != nil {
		utils.Response(w, r, 500, err.Error())
		log.Fatal(err.Error())
		return
	}

	utils.Response(w, r, 201, "page successfully created")
}

func (s *Service) updatePage(w http.ResponseWriter, r *http.Request) {
	// map body to json
	var page Page
	err := json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		utils.Response(w, r, 500, "error while decoding json")
		return
	}

	// get id in param
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.Response(w, r, 400, "impossible to parse page id")
		return
	}

	// validate page form
	err = validate.Struct(page)
	if err != nil {
		utils.Response(w, r, 400, "invalid input")
		return
	}

	// check if page in param exists
	exists, _ := s.queries.GetPage(context.Background(), int64(id))
	if exists.ID == 0 {
		utils.Response(w, r, 404, fmt.Sprintf("page with id %d doesn't exist", id))
		return
	}

	// update db
	errDB := s.queries.UpdatePage(context.Background(), db.UpdatePageParams{
		ID: int64(id),

		TitleDoUpdate: page.Title != "",
		Title:         page.Title,

		BodyDoUpdate: page.Body != "",
		Body:         page.Body,
	})
	if errDB != nil {
		utils.Response(w, r, 404, "page not found")
		return
	}

	utils.Response(w, r, 200, "page successfully updated")
}

func (s *Service) deletePage(w http.ResponseWriter, r *http.Request) {
	// get id
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.Response(w, r, 400, "page id bad format")
		return
	}

	// check if page in param exists
	exists, _ := s.queries.GetPage(context.Background(), int64(id))
	fmt.Println(exists)
	if exists.ID == 0 {
		utils.Response(w, r, 404, fmt.Sprintf("page with id %d doesn't exist", id))
		return
	}

	// delete page in db
	err = s.queries.DeletePage(context.Background(), int64(id))
	if err != nil {
		utils.Response(w, r, 404, "page not found")
		return
	}

	utils.Response(w, r, 200, "page successfully deleted")
}
