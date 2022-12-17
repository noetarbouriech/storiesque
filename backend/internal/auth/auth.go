package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/user"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries   *db.Queries
	tokenAuth *jwtauth.JWTAuth
	apiDomain string
}

func NewService(queries *db.Queries, jwt_secret string, api_domain string) *Service {
	return &Service{
		queries:   queries,
		tokenAuth: jwtauth.New("HS256", []byte(jwt_secret), nil),
		apiDomain: api_domain,
	}
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		utils.Response(w, r, 500, "error while decoding json")
		return
	}

	// get user
	user, err := s.queries.GetUserWithEmail(context.Background(), credentials.Email)
	if err != nil {
		utils.Response(w, r, 404, "account not found")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		utils.Response(w, r, 400, "wrong password")
		return
	}

	// create token
	expireTime := time.Now().Add(15 * time.Minute)
	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{
		"name": user.Username,     // username
		"id":   user.ID,           // user id
		"iat":  time.Now(),        // issued time
		"exp":  expireTime.Unix(), // expire time
	})
	if err != nil {
		utils.Response(w, r, 500, "error with token creation")
		return
	}

	// put token in client cookies
	http.SetCookie(w, &http.Cookie{
		Name:       "jwt",
		Value:      tokenString,
		Path:       "",
		Domain:     s.apiDomain,
		Expires:    expireTime,
		RawExpires: "",
		MaxAge:     10000,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "",
		Unparsed:   []string{},
	})

	utils.Response(w, r, 200, "successfully logged in")
}

func (s *Service) signUp(w http.ResponseWriter, r *http.Request) {
	var user user.UserCreation

	// decode json
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.Response(w, r, 500, "error while decoding json")
		return
	}

	// validate user form
	err = validate.Struct(user)
	if err != nil {
		utils.Response(w, r, 400, "invalid input")
		return
	}

	// check if user email already exists
	_, err = s.queries.GetUserWithEmail(context.Background(), user.Email)
	if err == nil {
		utils.Response(w, r, 400, "a user with this email address already exists")
		return
	}

	// check if username already exists
	_, err = s.queries.GetUserWithUsername(context.Background(), user.Username)
	if err == nil {
		utils.Response(w, r, 400, "a user with this username already exists")
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
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

	utils.Response(w, r, 201, "account successfully created")
}
