package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cscercel/behold-dnd/internal/db"
	_ "github.com/cscercel/behold-dnd/internal/db" // ONLY required for Swagger to pick up db interfaces
	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CharacterHandler struct {
	service	*service.CharacterService
}

func NewCharacterHandler(service *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{service: service}
}

func (h *CharacterHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/characters", func(r chi.Router) {
		// Public routes

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
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

func (h *CharacterHandler) handleUpdateCharacter(w http.ResponseWriter, r *http.Request) {
	var body db.UpdateCharacterParams
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.UpdateCharacter(r.Context(), body)
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
}

func (h *CharacterHandler) handleLongRest(w http.ResponseWriter, r *http.Request) {
}

func (h *CharacterHandler) handleShortRest(w http.ResponseWriter, r *http.Request) {
}

func (h *CharacterHandler) handleDeleteCharacter(w http.ResponseWriter, r *http.Request) {
}
