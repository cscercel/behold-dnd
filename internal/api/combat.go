package api

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      List all combat encounters
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   db.CombatEncounter
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat [get]
func (a *API) handleListEncounters(w http.ResponseWriter, r *http.Request) {
	encounters, err := a.queries.ListEncounters(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list encounters")
		return
	}

	respondJSON(w, http.StatusOK, encounters)
}

// @Summary      Get the currently active encounter
// @Tags         combat
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  db.CombatEncounter
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /combat/active [get]
func (a *API) handleGetActiveEncounter(w http.ResponseWriter, r *http.Request) {
	encounter, err := a.queries.GetActiveEncounter(r.Context())
	if err != nil {
		respondError(w, http.StatusNotFound, "no active encounter")
		return
	}

	respondJSON(w, http.StatusOK, encounter)
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
func (a *API) handleCreateEncounter(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	encounter, err := a.queries.CreateEncounter(r.Context(), body.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create encounter")
		return
	}

	respondJSON(w, http.StatusCreated, encounter)
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
func (a *API) handleGetEncounter(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid encounter id")
		return
	}

	encounter, err := a.queries.GetEncounter(r.Context(), encounterID)
	if err != nil {
		respondError(w, http.StatusNotFound, "encounter not found")
		return
	}

	respondJSON(w, http.StatusOK, encounter)
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
func (a *API) handleStartEncounter(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid encounter id")
		return
	}

	encounter, err := a.queries.StartEncounter(r.Context(), encounterID)
	if err != nil {
		respondError(w, http.StatusNotFound, "failed to start encounter")
		return
	}

	respondJSON(w, http.StatusOK, encounter)
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
func (a *API) handleEndEncounter(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid encounter id")
		return
	}

	encounter, err := a.queries.EndEncounter(r.Context(), encounterID)
	if err != nil {
		respondError(w, http.StatusNotFound, "failed to end encounter")
		return
	}

	respondJSON(w, http.StatusOK, encounter)
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
func (a *API) handleNextRound(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid encounter id")
		return
	}

	encounter, err := a.queries.NextRound(r.Context(), encounterID)
	if err != nil {
		respondError(w, http.StatusNotFound, "failed to advance round")
		return
	}

	respondJSON(w, http.StatusOK, encounter)
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
func (a *API) handleDeleteEncounter(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid encounter id")
		return
	}

	if err := a.queries.DeleteEncounter(r.Context(), encounterID); err != nil {
		respondError(w, http.StatusNotFound, "failed to advance round")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
func (a *API) handleListParticipants(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid encounter id")
		return
	}

	participants, err := a.queries.ListParticipants(r.Context(), encounterID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list participants")
		return
	}

	respondJSON(w, http.StatusOK, participants)
}

// @Summary      Add a participant to an encounter
// @Tags         combat
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        encounterID  path      string  true  "Encounter ID"
// @Param        body body object{character_id=string,initiative=int,name=string,current_hp=int,max_hp=int,armor_class=int,speed=int} true "Participant data. Provide character_id to copy stats from a character sheet, or fill in fields manually."
// @Success      201  {object}  db.CombatParticipant
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /combat/{encounterID}/participants [post]
func (a *API) handleAddParticipant(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid encounter id")
		return
	}

	var body struct {
		CharacterID string `json:"character_id"`
		Initiative  int32  `json:"initiative"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.CharacterID != "" {
		respondError(w, http.StatusBadRequest, "character_id is required")
		return
	}

	characterID, err := uuid.Parse(body.CharacterID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	participant, err := a.combatService.AddCharacterToEncounter(
		r.Context(), encounterID, characterID, body.Initiative,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to add participant")
		return
	}

	respondJSON(w, http.StatusCreated, participant)
	return
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
func (a *API) handleRemoveParticipant(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	if err := a.queries.RemoveParticipant(r.Context(), participantID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to remove participant")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
func (a *API) handleParticipantDamage(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be a positive number")
		return
	}

	participant, err := a.combatService.ApplyDamageToParticipant(r.Context(), participantID, body.Amount)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to apply damage")
		return
	}
	respondJSON(w, http.StatusOK, participant)
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
func (a *API) handleParticipantHeal(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be a positive number")
		return
	}

	participant, err := a.combatService.HealParticipant(r.Context(), participantID, body.Amount)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to heal participant")
		return
	}
	respondJSON(w, http.StatusOK, participant)
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
func (a *API) handleParticipantTempHP(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Amount int32 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be a positive number")
		return
	}

	participant, err := a.queries.UpdateParticipantTempHP(r.Context(), db.UpdateParticipantTempHPParams{
		ID:     participantID,
		TempHp: body.Amount,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update temp HP")
		return
	}

	respondJSON(w, http.StatusOK, participant)
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
func (a *API) handleParticipantInitiative(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Initiative int32 `json:"initiative"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Initiative <= 0 {
		respondError(w, http.StatusBadRequest, "initiative must be a positive number")
		return
	}

	participant, err := a.queries.UpdateParticipantInitiative(r.Context(), db.UpdateParticipantInitiativeParams{
		ID:         participantID,
		Initiative: body.Initiative,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update initiative")
		return
	}

	respondJSON(w, http.StatusOK, participant)
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
func (a *API) handleParticipantConditions(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Conditions []string `json:"conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	participant, err := a.queries.UpdateParticipantConditions(r.Context(), db.UpdateParticipantConditionsParams{
		ID:         participantID,
		Conditions: body.Conditions,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update conditions")
		return
	}

	respondJSON(w, http.StatusOK, participant)
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
func (a *API) handleParticipantToggleConcentration(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	participant, err := a.queries.ToggleParticipantConcentration(r.Context(), participantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to toggle concentration")
		return
	}

	respondJSON(w, http.StatusOK, participant)
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
func (a *API) handleDeactivateParticipant(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	participant, err := a.queries.DeactivateParticipant(r.Context(), participantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to deactivate participant")
		return
	}

	respondJSON(w, http.StatusOK, participant)
}
