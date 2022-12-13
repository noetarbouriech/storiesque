package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/noetarbouriech/storiesque/internal/db"
	"github.com/noetarbouriech/storiesque/internal/user"
	"golang.org/x/crypto/bcrypt"
)

var tokenAuth *jwtauth.JWTAuth

type Service struct {
	queries    *db.Queries
	jwt_secret string
}

func NewService(queries *db.Queries, jwt_secret string) *Service {
	return &Service{
		queries:    queries,
		jwt_secret: jwt_secret,
	}
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) init() {
	tokenAuth = jwtauth.New("HS256", []byte(s.jwt_secret), nil)
}

func (s *Service) PublicRoutes(r chi.Router) {
	r.Post("/login", s.login)
	r.Post("/signup", s.signUp)
}

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	// decode json
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"message": "error with json"})
		return
	}

	// get user
	user, err := s.queries.GetUserWithEmail(context.Background(), credentials.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"message": "account not found"})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"message": "wrong password"})
		return
	}

	// create token
	expireTime := time.Now().Add(15 * time.Minute)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"name": user.Username,     // username
		"id":   user.ID,           // user id
		"iat":  time.Now(),        // issued time
		"exp":  expireTime.Unix(), // expire time
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"message": "error with token creation"})
		return
	}

	// put token in client cookies
	http.SetCookie(w, &http.Cookie{
		Name:       "jwt",
		Value:      tokenString,
		Path:       "",
		Domain:     "localhost",
		Expires:    expireTime,
		RawExpires: "",
		MaxAge:     10000,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "",
		Unparsed:   []string{},
	})

	render.JSON(w, r, map[string]string{"message": "successfully logged in"})
}

func (s *Service) signUp(w http.ResponseWriter, r *http.Request) {
	var user user.UserCreation

	// decode json
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"message": "error with json"})
		return
	}

	// check fields
	if user.Email == "" || user.Username == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"message": "missing fields"})
		return
	}

	// check if user email already exists
	_, err = s.queries.GetUserWithEmail(context.Background(), user.Email)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"message": "user with this email already exists"})
		return
	}

	// check if username already exists
	_, err = s.queries.GetUserWithUsername(context.Background(), user.Username)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"message": "user with this username already exists"})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"message": "issue with password hashing"})
		return
	}

	// create user in db
	_, err = s.queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     user.Username,
		PasswordHash: string(hashedPassword),
		Email:        user.Email,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"message": "error with user creation in db"})
		log.Fatal(err.Error())
		return
	}

	render.JSON(w, r, map[string]string{"message": "successfully created"})
}
