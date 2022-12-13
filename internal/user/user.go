package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/noetarbouriech/storiesque/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

type UserCreation struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
}

func (s *Service) Routes(r chi.Router) {
	r.Get("/user", s.getUsers)
	r.Get("/user/{id}", s.getUser)
	r.Post("/user", s.createUser)
	r.Put("/user/{id}", s.updateUser)
	r.Delete("/user/{id}", s.deleteUser)

	r.Post("/login", s.login)
	r.Post("/signup", s.signUp)
}

func (s *Service) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.queries.ListUsers(context.Background())
	rUsers := []User{}
	for _, user := range users {
		rUser := User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		rUsers = append(rUsers, rUser)
	}
	if err != nil {
		render.JSON(w, r, map[string]string{"message": "no user found"})
		return
	}
	render.JSON(w, r, rUsers)
}

func (s *Service) getUser(w http.ResponseWriter, r *http.Request) {
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		render.JSON(w, r, map[string]string{"message": "impossible to parse id"})
		return
	}
	user, err := s.queries.GetUser(context.Background(), int64(id))
	if err != nil {
		render.JSON(w, r, map[string]string{"message": "no user found"})
		return
	}
	userJson := User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	render.JSON(w, r, userJson)
}

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	var user UserCreation
	errJson := json.NewDecoder(r.Body).Decode(&user)
	if errJson != nil {
		render.JSON(w, r, map[string]string{"message": "issue with json decoding"})
		return
	}

	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if errPassword != nil {
		render.JSON(w, r, map[string]string{"message": "issue with password hashing"})
		return
	}
	_, errDB := s.queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     user.Username,
		PasswordHash: string(hashedPassword),
		Email:        user.Email,
	})
	if errDB != nil {
		log.Fatal(errDB.Error())
		return
	}

	render.JSON(w, r, map[string]string{"message": "successfully created"})
}

func (s *Service) updateUser(w http.ResponseWriter, r *http.Request) {
	// map body to json
	var user UserCreation
	errJson := json.NewDecoder(r.Body).Decode(&user)
	if errJson != nil {
		render.JSON(w, r, map[string]string{"message": "issue with json decoding"})
		return
	}

	// hash password
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if errPassword != nil {
		render.JSON(w, r, map[string]string{"message": "issue with password hashing"})
		return
	}

	// get id in param
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		render.JSON(w, r, map[string]string{"message": "impossible to parse id"})
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
		render.JSON(w, r, map[string]string{"message": "no user found"})
		return
	}

	render.JSON(w, r, map[string]string{"message": "successfully updated"})
}

func (s *Service) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, errInt := strconv.Atoi(chi.URLParam(r, "id"))
	if errInt != nil {
		render.JSON(w, r, map[string]string{"message": "impossible to parse id"})
		return
	}
	errDB := s.queries.DeleteUser(context.Background(), int64(id))
	if errDB != nil {
		render.JSON(w, r, map[string]string{"message": "no user found"})
		return
	}
	render.JSON(w, r, map[string]string{"message": "successfully deleted"})
}
