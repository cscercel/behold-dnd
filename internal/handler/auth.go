package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/middleware"
)


// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body object{username=string,email=string,password=string,role=string,registration_code=string} true "Registration details"
// @Success      201  {object}  object{id=string,username=string,email=string,role=string}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /auth/register [post]
func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username 			string `json:"username"`
		Email				string	`json:"email"`
		Password			string	`json:"password"`
		Role				string	`json:"role"`
		RegistrationCode	string	`json:"registration_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Username == "" || body.Email == "" || body.Password == "" || body.Role == "" {
		respondError(w, http.StatusBadRequest, "username, email, password and role are missing")
		return
	}

	user, err := a.authService.Register(r.Context(), body.Username, body.Email, body.Password, body.Role, body.RegistrationCode)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"id":	user.ID,
		"username":	user.Username,
		"email": user.Email,
		"role": user.Role,
		"created_at": user.CreatedAt,
	})
}

// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body object{email=string,password=string} true "Login credentials"
// @Success      200  {object}  object{token=string}
// @Failure      401  {object}  object{error=string}
// @Router       /auth/login [post]
func (a *API) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email		string	`json:"email"`
		Password	string	`json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, user, err := a.authService.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user": map[string]any{
			"id": user.ID,
			"username": user.Username,
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
func (a *API) handleMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, err := a.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"id": user.ID,
		"username": user.Username,
		"email": user.Email,
		"role": user.Role,
	})
}
