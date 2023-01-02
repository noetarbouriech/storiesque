package user

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
	"github.com/noetarbouriech/storiesque/backend/internal/story"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

type UserCreation struct {
	Username string `json:"username" validate:"required,gte=2,lte=24,alphanum"`
	Password string `json:"password" validate:"required,gte=4,lte=128"`
	Email    string `json:"email"    validate:"required,gte=0,lte=128,email"`
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
	HasImg   bool   `json:"has_img"`
}

type UserDetails struct {
	Id       int64             `json:"id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	IsAdmin  bool              `json:"is_admin"`
	Stories  []story.StoryCard `json:"stories"`
	HasImg   bool              `json:"has_img"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (s *Service) PublicRoutes(r chi.Router) {
	r.Get("/user", s.getUsers)
	r.Get("/user/{username}", s.getUser)
}

func (s *Service) UserRoutes(r chi.Router) {
	r.Post("/user", s.createUser)
	r.Put("/user/{id}", s.updateUser)
	r.Put("/user/{id}/admin", s.setAdmin)
	r.Delete("/user/{id}", s.deleteUser)

	r.Get("/shelf", s.getShelf)
	r.Get("/shelf/{id}", s.isOnShelf)
	r.Post("/shelf/{id}", s.addToShelf)
	r.Delete("/shelf/{id}", s.removeFromShelf)
}

func (s *Service) getUsers(w http.ResponseWriter, r *http.Request) {

	// get stories with filter by name
	username := r.URL.Query().Get("username")

	// get page number in query
	page := r.URL.Query().Get("page")
	if len(page) == 0 {
		page = "1"
	}

	// get list of users from db
	users, err := s.queries.SearchUsers(context.Background(), db.SearchUsersParams{
		Column1: sql.NullString{String: username, Valid: true},
		Column2: page,
	})
	if err != nil {
		utils.Response(w, r, 404, "user not found")
		return
	}

	// map list to json
	rUsers := []User{}
	for _, user := range users {
		rUser := User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
			HasImg:   user.HasImg,
		}
		rUsers = append(rUsers, rUser)
	}

	render.JSON(w, r, rUsers)
}

func (s *Service) getUser(w http.ResponseWriter, r *http.Request) {

	// get username in uri
	user, err := s.queries.GetUserDetails(
		context.Background(),
		chi.URLParam(r, "username"),
	)
	if err != nil {
		utils.Response(w, r, 404, "user not found")
		return
	}

	// list stories written by given user
	stories := []story.StoryCard{}
	if user[0].StoryID.Valid {
		for _, line := range user {
			story := story.StoryCard{
				Id:          line.StoryID.Int64,
				Title:       line.Title.String,
				Description: line.Description.String,
				AuthorName:  line.Username,
				HasImg:      line.StoryHasImg.Bool,
			}
			stories = append(stories, story)
		}
	}

	// render json of user
	render.JSON(w, r, UserDetails{
		Id:       user[0].ID,
		Username: user[0].Username,
		Email:    user[0].Email,
		IsAdmin:  user[0].IsAdmin,
		HasImg:   user[0].HasImg,
		Stories:  stories,
	})
}

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	var user UserCreation

	// decode json in struct
	errJson := json.NewDecoder(r.Body).Decode(&user)
	if errJson != nil {
		utils.Response(w, r, 500, "error while decoding json")
		return
	}

	// validate user form
	err := validate.Struct(user)
	if err != nil {
		utils.Response(w, r, 400, "invalid input")
		return
	}

	// hash password
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if errPassword != nil {
		utils.Response(w, r, 500, "error while hashing password")
		return
	}

	// create user in db
	_, err = s.queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     user.Username,
		PasswordHash: string(hashedPassword),
		Email:        user.Email,
	})
	if err != nil {
		utils.Response(w, r, 500, err.Error())
		log.Fatal(err.Error())
		return
	}

	utils.Response(w, r, 200, "user successfully created")
}

func (s *Service) updateUser(w http.ResponseWriter, r *http.Request) {

	// get id in param
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse user id")
		return
	}

	// check if user is authorized
	if !utils.IsOwner(r, id) {
		utils.Response(w, r, 401, "user is not the owner of the given resource")
		return
	}

	// map body to json
	var user UserCreation
	errJson := json.NewDecoder(r.Body).Decode(&user)
	if errJson != nil {
		utils.Response(w, r, 500, "error while decoding json")
		return
	}

	// hash password
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if errPassword != nil {
		utils.Response(w, r, 500, "error while hashing password")
		return
	}

	// update db
	errDB := s.queries.UpdateUser(context.Background(), db.UpdateUserParams{
		ID: int64(id),

		UsernameDoUpdate: user.Username != "",
		Username:         user.Username,

		PasswordHashDoUpdate: user.Password != "",
		PasswordHash:         string(hashedPassword),

		EmailDoUpdate: user.Email != "",
		Email:         user.Email,
	})
	if errDB != nil {
		utils.Response(w, r, 404, "user not found")
		return
	}

	utils.Response(w, r, 200, "user successfully updated")
}

func (s *Service) setAdmin(w http.ResponseWriter, r *http.Request) {

	// get id in param
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse user id")
		return
	}

	// check if user is authorized
	if !utils.IsAdmin(r) {
		utils.Response(w, r, 401, "user is not admin")
		return
	}

	// set user as admin in db
	errDB := s.queries.SetAdmin(context.Background(), int64(id))
	if errDB != nil {
		utils.Response(w, r, 404, "user not found")
		return
	}

	utils.Response(w, r, 200, "user successfully updated")
}

func (s *Service) deleteUser(w http.ResponseWriter, r *http.Request) {

	// get id in param
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse user id")
		return
	}

	// check if user is authorized
	if !utils.IsOwner(r, id) {
		utils.Response(w, r, 401, "user is not the owner of the given resource")
		return
	}

	// delete user in db
	errDB := s.queries.DeleteUser(context.Background(), int64(id))
	if errDB != nil {
		utils.Response(w, r, 404, "user not found")
		return
	}

	utils.Response(w, r, 200, "user successfully deleted")
}
