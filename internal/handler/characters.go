package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/db"
	_ "github.com/cscercel/behold-dnd/internal/db" // ONLY required for Swagger to pick up db interfaces
	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/cscercel/behold-dnd/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CharacterHandler struct {
	service	*service.CharacterService
	authService *service.AuthService
}

func NewCharacterHandler(service *service.CharacterService, authService *service.AuthService) *CharacterHandler {
	return &CharacterHandler{service: service, authService: authService}
}

func (h *CharacterHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/characters", func(r chi.Router) {
		// Player routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(h.authService))

			r.Get("/auth/me", h.handleMe)

			r.Get("/", h.handleListCharacters)
			r.Post("/", h.handleCreateCharacter)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.handleGetCharacter)
				r.Patch("/", h.handleUpdateCharacter)
				r.Delete("/", h.handleDeleteCharacter)


				// Game mechanics
				r.Post("/damage", h.handleDamage)
				r.Post("/heal", h.handleHeal)
				r.Post("/temp-hp", h.handleTempHP)
				r.Post("/death-save", h.handleDeathSave)
				r.Post("/long-rest", h.handleLongRest)
				r.Post("/short-rest", h.handleShortRest)
				r.Put("/conditions", h.handleUpdateConditions)
			})
		})

		// DM routes
		r.Group(func(r chi.Router) {
		})
	})
}

func (h *CharacterHandler) handleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var body db.CreateCharacterParams

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.CreateCharacter(r.Context(), body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create character", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, character)
}

func (h *CharacterHandler) handleGetCharacter(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	character, err := h.service.GetCharacter(r.Context(), characterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleListCharacters(w http.ResponseWriter, r *http.Request) {
	characters, err := h.service.ListCharacters(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get characters", err)
		return
	}

	respondWithJSON(w, http.StatusOK, characters)
}

func (h *CharacterHandler) handleListUserCharacters(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user id", err)
		return
	}
	characters, err := h.service.ListUserCharacters(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get user characters", err)
		return
	}

	respondWithJSON(w, http.StatusOK, characters)
}

func (h *CharacterHandler) handleListPlayerCharacters(w http.ResponseWriter, r *http.Request) {
	characters, err := h.service.ListPlayerCharacters(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get player characters", err)
		return
	}

	respondWithJSON(w, http.StatusOK, characters)
}

func (h *CharacterHandler) handleListNPCs(w http.ResponseWriter, r *http.Request) {
	npcs, err := h.service.ListNPCs(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get npcs", err)
		return
	}

	respondWithJSON(w, http.StatusOK, npcs)
}

func (h *CharacterHandler) handleUpdateCharacterInfo(w http.ResponseWriter, r *http.Request) {
	var body db.UpdateCharacterInfoParams
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateCharacterInfo(r.Context(), body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update player", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleUpdateCharacterAbilityScores(w http.ResponseWriter, r *http.Request) {
	var body db.UpdateCharacterAbilityScoresParams
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateCharacterAbilityScores(r.Context(), body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update player", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleUpdateCharacterSkills(w http.ResponseWriter, r *http.Request) {
	var body db.UpdateCharacterSkillsParams
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateCharacterSkills(r.Context(), body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update player", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleUpdateCharacterLevel(w http.ResponseWriter, r *http.Request) {
	var body db.UpdateCharacterLevelParams
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateCharacterLevel(r.Context(), body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update player", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleUpdateCharacterTraining(w http.ResponseWriter, r *http.Request) {
	var body db.UpdateCharacterTrainingParams
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateCharacterTraining(r.Context(), body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update player", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleUpdateCharacterCurrency(w http.ResponseWriter, r *http.Request) {
	var body db.UpdateCharacterCurrencyParams
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateCharacterCurrency(r.Context(), body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update player", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleHealCharacter(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var body struct {
		Amount	int32	`json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.HealCharacter(r.Context(), characterID, body.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not heal character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleDamageCharacter(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var body struct {
		Amount	int32	`json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.DamageCharacter(r.Context(), characterID, body.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not apply damage to character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleUpdateCharacterTempHP(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var body struct {
		Amount	int32	`json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateCharacterTempHP(r.Context(), characterID, body.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update character temp hp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleUpdateCharacterDeathSave(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var body struct {
		Success	bool `json:"success"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateCharacterDeathSave(r.Context(), characterID, body.Success)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update character death saves", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleUpdateConditions(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var body struct {
		Conditions	[]string `json:"conditions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateConditions(r.Context(), characterID, body.Conditions)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update character conditions", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleLongRest(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	character, err := h.service.LongRest(r.Context(), characterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not long rest character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleShortRest(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var body struct {
		HitDiceRemaining int32	`json:"hit_dice_remaining"`
		CurrentHp	int32	`json:"current_hp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.ShortRest(r.Context(), characterID, body.HitDiceRemaining, body.CurrentHp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not short rest character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

func (h *CharacterHandler) handleDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if err := h.service.DeleteCharacter(r.Context(), characterID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete character", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
