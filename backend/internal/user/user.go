package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
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
	Username string `json:"username" validate:"required,gte=0,lte=24"`
	Password string `json:"password" validate:"required,gte=0,lte=128"`
	Email    string `json:"email"    validate:"required,gte=0,lte=128,email"`
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
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
	r.Delete("/user/{id}", s.deleteUser)
}

func (s *Service) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.queries.ListUsers(context.Background())
	if err != nil {
		utils.Response(w, r, 404, "user not found")
		return
	}
	rUsers := []User{}
	for _, user := range users {
		rUser := User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
		}
		rUsers = append(rUsers, rUser)
	}
	render.JSON(w, r, rUsers)
}

func (s *Service) getUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := s.queries.GetUserWithUsername(context.Background(), username)
	if err != nil {
		utils.Response(w, r, 404, "user not found")
		return
	}
	userJson := User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}
	render.JSON(w, r, userJson)
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

	// get id in param
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse user id")
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

func (s *Service) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		utils.Response(w, r, 400, "impossible to parse user id")
		return
	}
	errDB := s.queries.DeleteUser(context.Background(), int64(id))
	if errDB != nil {
		utils.Response(w, r, 404, "user not found")
		return
	}
	utils.Response(w, r, 200, "user successfully deleted")
}
