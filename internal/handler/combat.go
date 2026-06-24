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

type CombatHandler struct {
	service	*service.CombatService
	authService *service.AuthService
}

func NewCombatHandler(service *service.CombatService, authService *service.AuthService) *CombatHandler {
	return &CombatHandler{service: service, authService: authService}
}

func (h *CombatHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/encounters", func(r chi.Router) {
	})
}

func (h *CombatHandler) handleCreateEncounter(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	encounter, err := h.service.CreateEncouter(r.Context(), body.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create encounter", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, encounter)
} 

func (h *CombatHandler) handleGetEncounter(w http.ResponseWriter, r *http.Request) {
	combatID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	encounter, err := h.service.GetEncounter(r.Context(), combatID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to load encounter", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounter)
}

func (h *CombatHandler) handleListEncounters(w http.ResponseWriter, r *http.Request) {
	encounters, err := h.service.ListEncouters(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to load encounters", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounters)
}

func (h *CombatHandler) handleGetActiveEncounters(w http.ResponseWriter, r *http.Request) {
	encounters, err := h.service.GetActiveEncouters(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to load active encounters", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounters)
}

func (h *CombatHandler) handleStartEncounter(w http.ResponseWriter, r *http.Request) {
	combatID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	encounter, err := h.service.StartEncounter(r.Context(), combatID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to start encounter", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounter)
}

func (h *CombatHandler) handleEndEncounter(w http.ResponseWriter, r *http.Request) {
	combatID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	encounter, err := h.service.EndEncounter(r.Context(), combatID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to end encounter", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounter)
}

func (h *CombatHandler) handleNextRound(w http.ResponseWriter, r *http.Request) {
	combatID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	encounter, err := h.service.NextRound(r.Context(), combatID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to move encounter to next round", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounter)
}

func (h *CombatHandler) handleDeleteEncounter(w http.ResponseWriter, r *http.Request) {
	combatID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	if err := h.service.DeleteEncounter(r.Context(), combatID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete encounter", err)
		return
	}
	
	respondWithJSON(w, http.StatusNoContent, "")
}

func (h *CombatHandler) handleGetParticipant(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid participant id", err)
		return
	}

	participant, err := h.service.GetParticipant(r.Context(), participantID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to load participant", err)
		return
	}

	respondWithJSON(w, http.StatusOK, participant)
}
