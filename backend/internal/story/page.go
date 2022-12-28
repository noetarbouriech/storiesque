package story

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
)

type PageChoices struct {
	Id     int64  `json:"page_id"`
	Action string `json:"action"`
}

type Page struct {
	Id      int64         `json:"id"`
	Author  int64         `json:"author_id"`
	Action  string        `json:"action"             validate:"gte=0,lte=64"`
	Body    string        `json:"body"              validate:"gte=0,lte=4096"`
	Choices []PageChoices `json:"choices"`
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

	// map choices
	choicesStruct := []PageChoices{}
	for _, choice := range choices {
		choicesStruct = append(choicesStruct, PageChoices{
			Id:     choice.PathID,
			Action: choice.Action,
		})
	}

	// render page in json
	render.JSON(w, r, Page{
		Id:      page.ID,
		Author:  page.Author,
		Action:  page.Action,
		Body:    page.Body,
		Choices: choicesStruct,
	})
}

func (s *Service) createPage(w http.ResponseWriter, r *http.Request) {

	// get id in param
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse page id")
		return
	}

	// check if page in param exists
	parentPage, _ := s.queries.GetPage(context.Background(), int64(id))
	if parentPage.ID == 0 {
		utils.Response(w, r, 404, fmt.Sprintf("page with id %d doesn't exist", id))
		return
	}

	// check if user is authorized
	if !utils.IsOwner(r, int(parentPage.Author)) {
		utils.Response(w, r, 401, "user is not the owner of the given resource")
		return
	}

	// get infos from jwt
	// cannot be error since the jwt is verified for userRouters
	_, claims, _ := jwtauth.FromContext(r.Context())

	// create page in db
	page, err := s.queries.CreatePage(context.Background(), db.CreatePageParams{
		Action: "New choice",
		Author: int64(claims["id"].(float64)),
		Body:   "Hello world !",
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

	// print json of the choice added
	render.JSON(w, r, PageChoices{
		Action: page.Action,
		Id:     page.ID,
	})
}

func (s *Service) updatePage(w http.ResponseWriter, r *http.Request) {

	// get id in param
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.Response(w, r, 400, "impossible to parse page id")
		return
	}

	// check if page in param exists
	parentPage, _ := s.queries.GetPage(context.Background(), int64(id))
	if parentPage.ID == 0 {
		utils.Response(w, r, 404, fmt.Sprintf("page with id %d doesn't exist", id))
		return
	}

	// check if user is authorized
	if !utils.IsOwner(r, int(parentPage.Author)) {
		utils.Response(w, r, 401, "user is not the owner of the given resource")
		return
	}

	// map body to json
	var page Page
	err = json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		utils.Response(w, r, 500, "error while decoding json")
		return
	}

	// validate page form
	err = validate.Struct(page)
	if err != nil {
		utils.Response(w, r, 400, "invalid input")
		return
	}

	// update db
	err = s.queries.UpdatePage(context.Background(), db.UpdatePageParams{
		ID: int64(id),

		ActionDoUpdate: page.Action != "",
		Action:         page.Action,

		BodyDoUpdate: page.Body != "",
		Body:         page.Body,
	})
	if err != nil {
		utils.Response(w, r, 404, "page not found")
		return
	}

	utils.Response(w, r, 200, "page successfully updated")
}

func (s *Service) deletePage(w http.ResponseWriter, r *http.Request) {

	// get id in param
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.Response(w, r, 400, "page id bad format")
		return
	}

	// check if page in param exists
	parentPage, _ := s.queries.GetPage(context.Background(), int64(id))
	if parentPage.ID == 0 {
		utils.Response(w, r, 404, fmt.Sprintf("page with id %d doesn't exist", id))
		return
	}

	// check if user is authorized
	if !utils.IsOwner(r, int(parentPage.Author)) {
		utils.Response(w, r, 401, "user is not the owner of the given resource")
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
