package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/story"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
)

func (s *Service) getShelf(w http.ResponseWriter, r *http.Request) {

	// get infos from jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		utils.Response(w, r, http.StatusUnauthorized, "not logged in")
		return
	}

	// get page number in query
	page := r.URL.Query().Get("page")
	if len(page) == 0 {
		page = "1"
	}

	// get shelf of current user
	stories, err := s.queries.GetShelf(context.Background(), db.GetShelfParams{
		OwnerID: int64(claims["id"].(float64)),
		Column2: page,
	})
	if err != nil {
		utils.Response(w, r, 404, "shelf is empty")
		return
	}

	// map list to json
	rStories := []story.StoryCard{}
	for _, storyDB := range stories {
		rStory := story.StoryCard{
			Id:          storyDB.ID,
			Title:       storyDB.Title,
			Description: storyDB.Description.String,
			AuthorName:  storyDB.AuthorName,
		}
		rStories = append(rStories, rStory)
	}

	render.JSON(w, r, rStories)
}

func (s *Service) isOnShelf(w http.ResponseWriter, r *http.Request) {

	// get infos from jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		utils.Response(w, r, http.StatusUnauthorized, "not logged in")
		return
	}

	// get id in url
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse story id")
		return
	}

	// get shelf of current user
	_, err = s.queries.GetOnShelf(context.Background(), db.GetOnShelfParams{
		OwnerID: int64(claims["id"].(float64)),
		StoryID: int64(id),
	})
	if err != nil {
		utils.Response(w, r, 404, "story not on shelf")
		return
	}

	utils.Response(w, r, 200, "story is in shelf")
}

func (s *Service) addToShelf(w http.ResponseWriter, r *http.Request) {

	// get infos from jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		utils.Response(w, r, http.StatusUnauthorized, "not logged in")
		return
	}

	// get id in url
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse story id")
		return
	}

	// add story in shelf in db
	err = s.queries.AddToShelf(context.Background(), db.AddToShelfParams{
		OwnerID: int64(claims["id"].(float64)),
		StoryID: int64(id),
	})
	if err != nil {
		utils.Response(w, r, 400, "story already in shelf")
		return
	}

	utils.Response(w, r, 200, "story added to shelf successfully")
}

func (s *Service) removeFromShelf(w http.ResponseWriter, r *http.Request) {

	// get infos from jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		utils.Response(w, r, http.StatusUnauthorized, "not logged in")
		return
	}

	// get id in url
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse story id")
		return
	}

	// add story in shelf in db
	err = s.queries.RemoveFromShelf(context.Background(), db.RemoveFromShelfParams{
		OwnerID: int64(claims["id"].(float64)),
		StoryID: int64(id),
	})
	if err != nil {
		utils.Response(w, r, 500, err.Error())
		return
	}

	utils.Response(w, r, 200, "story removed from shelf successfully")
}
