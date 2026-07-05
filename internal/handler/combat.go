package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CombatHandler struct {
	service *service.CombatService
}

func NewCombatHandler(service *service.CombatService) *CombatHandler {
	return &CombatHandler{service: service}
}

func (h *CombatHandler) RegisterRoutes(r chi.Router, authMiddleware, dmOnlyMiddleware func(http.Handler) http.Handler) {
	r.Route("/combat", func(r chi.Router) {
		r.Use(authMiddleware, dmOnlyMiddleware)

		r.Get("/", h.handleListEncounters)
		r.Post("/", h.handleCreateEncounter)

		r.Route("/{encounterID}", func(r chi.Router) {
			r.Get("/", h.handleGetEncounter)
			r.Delete("/", h.handleDeleteEncounter)
			r.Post("/start", h.handleStartEncounter)
			r.Post("/end", h.handleEndEncounter)
			r.Post("/next-round", h.handleNextRound)

			r.Get("/participants", h.handleListParticipants)
			r.Post("/participants", h.handleAddParticipant)

			r.Route("/participants/{participantID}", func(r chi.Router) {
				r.Delete("/", h.handleRemoveParticipant)
				r.Post("/damage", h.handleParticipantDamage)
				r.Post("/heal", h.handleParticipantHeal)
				r.Post("/temp-hp", h.handleParticipantTempHP)
				r.Put("/initiative", h.handleParticipantInitiative)
				r.Put("/conditions", h.handleParticipantConditions)
				r.Post("/toggle-concentration", h.handleParticipantToggleConcentration)
				r.Post("/deactivate", h.handleDeactivateParticipant)
			})
		})
	})
}

// @Summary      List all combat encounters
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   db.CombatEncounter
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat [get]
func (h *CombatHandler) handleListEncounters(w http.ResponseWriter, r *http.Request) {
	encounters, err := h.service.ListEncounters(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list encounters", err)
		return
	}
	respondWithJSON(w, http.StatusOK, encounters)
}

// @Summary      Create a combat encounter
// @Tags         combat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body      object{name=string}  true  "Encounter name"
// @Success      201  {object}  db.CombatEncounter
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat [post]
func (h *CombatHandler) handleCreateEncounter(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	encounter, err := h.service.CreateEncounter(r.Context(), body.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create encounter", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, encounter)
}

// @Summary      Get a combat encounter
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID  path      string  true  "Encounter ID"
// @Success      200  {object}  db.CombatEncounter
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /combat/{encounterID} [get]
func (h *CombatHandler) handleGetEncounter(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	encounter, err := h.service.GetEncounter(r.Context(), encounterID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "encounter not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounter)
}

// @Summary      Start a combat encounter
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID  path      string  true  "Encounter ID"
// @Success      200  {object}  db.CombatEncounter
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/start [post]
func (h *CombatHandler) handleStartEncounter(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	encounter, err := h.service.StartEncounter(r.Context(), encounterID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed to start encounter", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounter)
}

// @Summary      End a combat encounter
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID  path      string  true  "Encounter ID"
// @Success      200  {object}  db.CombatEncounter
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/end [post]
func (h *CombatHandler) handleEndEncounter(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	encounter, err := h.service.EndEncounter(r.Context(), encounterID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed to end encounter", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounter)
}

// @Summary      Advance to the next round
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID  path      string  true  "Encounter ID"
// @Success      200  {object}  db.CombatEncounter
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/next-round [post]
func (h *CombatHandler) handleNextRound(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	encounter, err := h.service.NextRound(r.Context(), encounterID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed to advance round", err)
		return
	}

	respondWithJSON(w, http.StatusOK, encounter)
}

// @Summary      Delete a combat encounter
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID  path      string  true  "Encounter ID"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID} [delete]
func (h *CombatHandler) handleDeleteEncounter(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	if err := h.service.DeleteEncounter(r.Context(), encounterID); err != nil {
		respondWithError(w, http.StatusNotFound, "failed to delete encounter", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}

// @Summary      List participants in an encounter
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID  path      string  true  "Encounter ID"
// @Success      200  {array}   db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants [get]
func (h *CombatHandler) handleListParticipants(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	participants, err := h.service.ListParticipants(r.Context(), encounterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list participants", err)
		return
	}

	respondWithJSON(w, http.StatusOK, participants)
}

// @Summary      Add a participant to an encounter
// @Tags         combat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID  path      string  true  "Encounter ID"
// @Param        body body object{character_id=string,initiative=int} true "Participant data. Copies stats from the character sheet."
// @Success      201  {object}  db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants [post]
func (h *CombatHandler) handleAddParticipant(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid encounter id", err)
		return
	}

	var body struct {
		CharacterID string `json:"character_id"`
		Initiative  int32  `json:"initiative"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if body.CharacterID == "" {
		respondWithError(w, http.StatusBadRequest, "character_id is required", nil)
		return
	}

	characterID, err := uuid.Parse(body.CharacterID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	participant, err := h.service.AddCharacterToEncounter(r.Context(), encounterID, characterID, body.Initiative)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to add participant", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, participant)
}

// @Summary      Remove a participant from an encounter
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID    path      string  true  "Encounter ID"
// @Param        participantID  path      string  true  "Participant ID"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants/{participantID} [delete]
func (h *CombatHandler) handleRemoveParticipant(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid participant id", err)
		return
	}

	if err := h.service.RemoveParticipant(r.Context(), participantID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to remove participant", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}

// @Summary      Deal damage to a participant
// @Tags         combat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID    path      string             true  "Encounter ID"
// @Param        participantID  path      string             true  "Participant ID"
// @Param        body           body      object{amount=int} true  "Damage amount"
// @Success      200  {object}  db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants/{participantID}/damage [post]
func (h *CombatHandler) handleParticipantDamage(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid participant id", err)
		return
	}

	var body struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount must be a positive number", err)
		return
	}

	participant, err := h.service.ApplyDamageToParticipant(r.Context(), participantID, body.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to apply damage", err)
		return
	}

	respondWithJSON(w, http.StatusOK, participant)
}

// @Summary      Heal a participant
// @Tags         combat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID    path      string             true  "Encounter ID"
// @Param        participantID  path      string             true  "Participant ID"
// @Param        body           body      object{amount=int} true  "Heal amount"
// @Success      200  {object}  db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants/{participantID}/heal [post]
func (h *CombatHandler) handleParticipantHeal(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid participant id", err)
		return
	}

	var body struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount must be a positive number", err)
		return
	}

	participant, err := h.service.HealParticipant(r.Context(), participantID, body.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to heal participant", err)
		return
	}

	respondWithJSON(w, http.StatusOK, participant)
}

// @Summary      Give Temp HP to a participant
// @Tags         combat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID    path      string             true  "Encounter ID"
// @Param        participantID  path      string             true  "Participant ID"
// @Param        body           body      object{amount=int} true  "TempHP amount"
// @Success      200  {object}  db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants/{participantID}/temp-hp [post]
func (h *CombatHandler) handleParticipantTempHP(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid participant id", err)
		return
	}

	var body struct {
		Amount int32 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount must be a positive number", err)
		return
	}

	participant, err := h.service.AddParticipantTempHP(r.Context(), participantID, body.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update temp HP", err)
		return
	}

	respondWithJSON(w, http.StatusOK, participant)
}

// @Summary      Update a participant's initiative
// @Tags         combat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID    path      string                  true  "Encounter ID"
// @Param        participantID  path      string                  true  "Participant ID"
// @Param        body           body      object{initiative=int}  true  "Initiative value"
// @Success      200  {object}  db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants/{participantID}/initiative [put]
func (h *CombatHandler) handleParticipantInitiative(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid participant id", err)
		return
	}

	var body struct {
		Initiative int32 `json:"initiative"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Initiative <= 0 {
		respondWithError(w, http.StatusBadRequest, "initiative must be a positive number", err)
		return
	}

	participant, err := h.service.UpdateParticipantInitiative(r.Context(), participantID, body.Initiative)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update initiative", err)
		return
	}

	respondWithJSON(w, http.StatusOK, participant)
}

// @Summary      Update conditions on a participant
// @Tags         combat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID    path      string                       true  "Encounter ID"
// @Param        participantID  path      string                       true  "Participant ID"
// @Param        body           body      object{conditions=[]string}  true  "Conditions list"
// @Success      200  {object}  db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants/{participantID}/conditions [put]
func (h *CombatHandler) handleParticipantConditions(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid participant id", err)
		return
	}

	var body struct {
		Conditions []string `json:"conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	participant, err := h.service.UpdateParticipantConditions(r.Context(), participantID, body.Conditions)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update conditions", err)
		return
	}

	respondWithJSON(w, http.StatusOK, participant)
}

// @Summary      Toggle concentration for a participant
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID    path      string  true  "Encounter ID"
// @Param        participantID  path      string  true  "Participant ID"
// @Success      200  {object}  db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants/{participantID}/toggle-concentration [post]
func (h *CombatHandler) handleParticipantToggleConcentration(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid participant id", err)
		return
	}

	participant, err := h.service.ToggleParticipantConcentration(r.Context(), participantID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to toggle concentration", err)
		return
	}

	respondWithJSON(w, http.StatusOK, participant)
}

// @Summary      Deactivate a participant (knocked out or fled)
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID    path      string  true  "Encounter ID"
// @Param        participantID  path      string  true  "Participant ID"
// @Success      200  {object}  db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants/{participantID}/deactivate [post]
func (h *CombatHandler) handleDeactivateParticipant(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid participant id", err)
		return
	}

	participant, err := h.service.DeactivateParticipant(r.Context(), participantID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to deactivate participant", err)
		return
	}

	respondWithJSON(w, http.StatusOK, participant)
}
