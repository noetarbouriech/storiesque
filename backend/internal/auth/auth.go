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
	"github.com/go-playground/validator/v10"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/user"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
	passwordvalidator "github.com/wagslane/go-password-validator"
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
	r.Get("/logout", s.logout)
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
	userDB, err := s.queries.GetUserWithEmail(context.Background(), credentials.Email)
	if err != nil {
		utils.Response(w, r, 404, "account not found")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash), []byte(credentials.Password)); err != nil {
		utils.Response(w, r, 400, "wrong password")
		return
	}

	// create access token
	accessToken, err := s.CreateAccessToken(userDB)
	if err != nil {
		utils.Response(w, r, 500, "error with token creation")
		return
	}
	http.SetCookie(w, accessToken)

	// create refresh token
	refreshToken, err := s.CreateRefreshToken()
	if err != nil {
		utils.Response(w, r, 500, "error with token creation")
		return
	}
	http.SetCookie(w, refreshToken)

	render.JSON(w, r, user.User{
		Id:       userDB.ID,
		Username: userDB.Username,
		Email:    userDB.Email,
		IsAdmin:  userDB.IsAdmin,
	})
}

func (s *Service) logout(w http.ResponseWriter, r *http.Request) {

	// replace access token with an expired one
	http.SetCookie(w, &http.Cookie{
		Name:       "jwt",
		Value:      "",
		Path:       "",
		Domain:     s.apiDomain,
		Expires:    time.Unix(0, 0),
		RawExpires: "",
		MaxAge:     0,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "",
		Unparsed:   []string{},
	})

	// replace refresh token with an expired one
	http.SetCookie(w, &http.Cookie{
		Name:       "refresh",
		Value:      "",
		Path:       "",
		Domain:     s.apiDomain,
		Expires:    time.Unix(0, 0),
		RawExpires: "",
		MaxAge:     0,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
		Raw:        "",
		Unparsed:   []string{},
	})

	utils.Response(w, r, 200, "successfully logged out")
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

	// check if password is secure enough
	err = passwordvalidator.Validate(user.Password, 60)
	if err != nil {
		utils.Response(w, r, 400, err.Error())
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
