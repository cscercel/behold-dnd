package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/cscercel/behold-dnd/internal/db" // ONLY required for Swagger to pick up db interfaces
	"github.com/cscercel/behold-dnd/internal/middleware"
	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	service	*service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/auth", func(r chi.Router) {
		// Public Routes
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)

		// Private Routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(h.service))

			r.Get("/me", h.handleMe)
		})
	})
}

// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body object{username=string,email=string,password=string,role=string,registration_code=string} true "Registration details"
// @Success      201  {object}  object{id=string,username=string,email=string,role=string}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /auth/register [post]
func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email				string	`json:"email"`
		Password			string	`json:"password"`
		Role				string	`json:"role"`
		RegistrationCode	string	`json:"registration_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if body.Email == "" || body.Password == "" || body.Role == "" {
		respondWithError(w, http.StatusBadRequest, "email, password and role are missing", fmt.Errorf(""))
		return
	}

	user, err := h.service.Register(r.Context(), body.Email, body.Password, body.Role, body.RegistrationCode)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to register", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body object{email=string,password=string} true "Login credentials"
// @Success      200  {object}  object{token=string}
// @Failure      401  {object}  object{error=string}
// @Router       /auth/login [post]
func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email		string	`json:"email"`
		Password	string	`json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	token, user, err := h.service.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to login", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user": map[string]any{
			"id": user.ID,
			"email": user.Email,
			"role": user.Role,
		},
	})
}

// @Summary      Get current user
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{id=string,username=string,email=string,role=string}
// @Failure      401  {object}  object{error=string}
// @Router       /auth/me [get]
func (h *AuthHandler) handleMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "not authenticated", fmt.Errorf(""))
		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
