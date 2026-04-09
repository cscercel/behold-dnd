package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/cscercel/beyond-dnd/internal/db"
	"github.com/google/uuid"
)


// Encounter handlers
func (a *API) handleListEncounters(w http.ResponseWriter, r *http.Request) {
	encounters, err := a.queries.ListEncounters(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list encounters")
		return
	}

	respondJSON(w, http.StatusOK, encounters)
}

func (a *API) handleGetActiveEncounter(w http.ResponseWriter, r *http.Request) {
	encounter, err := a.queries.GetActiveEncounter(r.Context())
	if err != nil {
		respondError(w, http.StatusNotFound, "no active encounter")
		return
	}

	respondJSON(w, http.StatusOK, encounter)
}

func (a *API) handleCreateEncounter(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name	string	`json:"name"`
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


// Participant handlers
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

// handleAddParticipant accepts either a character_id (full sheet) or 
// manual stats (for quick ad-hoc enemies)
func (a *API) handleAddParticipant(w http.ResponseWriter, r *http.Request) {
	encounterID, err := uuid.Parse(chi.URLParam(r, "encounterID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid encounter id")
		return
	}

	var body struct {
		CharacterID	string	`json:"character_id"`
		Initiative	int32	`json:"initiative"`

		// Only used when adding directly from a character sheet
		Name		string	`json:"name"`
		CurrentHP	int32	`json:"current_hp"`
		MaxHP		int32	`json:"max_hp"`
		TempHP		int32	`json:"temp_hp"`
		ArmorClass	int32	`json:"armor_class"`
		Speed		int32	`json:"speed"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// If a character_id was provided, prefill stats from character sheet
	if body.CharacterID != "" {
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

	// Otherwise use manually provided stats
	if body.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required when not using a character id")
		return
	}

	participant, err := a.queries.AddParticipant(r.Context(), db.AddParticipantParams{
		EncounterID: encounterID,
		Name: body.Name,
		Initiative: body.Initiative,
		CurrentHp: body.CurrentHP,
		MaxHp: body.MaxHP,
		TempHp: body.TempHP,
		ArmorClass: body.ArmorClass,
		Speed: body.Speed,
	})
	
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to add participant")
		return
	}

	respondJSON(w, http.StatusCreated, participant)
} 

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

func (a *API) handleParticipantDamage(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Amount	int	`json:"amount"`
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

func (a *API) handleParticipantHeal(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Amount	int	`json:"amount"`
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

func (a *API) handleUpdateTempHP(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Amount	int32	`json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be a positive number")
		return
	}

	participant, err := a.queries.UpdateParticipantTempHP(r.Context(), db.UpdateParticipantTempHPParams{
		ID: participantID,
		TempHp: body.Amount,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update temp HP")
		return
	}

	respondJSON(w, http.StatusOK, participant)
}

func (a *API) handleUpdateParticipantInitiative(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Initiative	int32	`json:"initiative"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Initiative <= 0 {
		respondError(w, http.StatusBadRequest, "initiative must be a positive number")
		return
	}

	participant, err := a.queries.UpdateParticipantInitiative(r.Context(), db.UpdateParticipantInitiativeParams{
		ID: participantID,
		Initiative: body.Initiative,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update initiative")
		return
	}

	respondJSON(w, http.StatusOK, participant)
}

func (a *API) handleUpdateParticipantConditions(w http.ResponseWriter, r *http.Request) {
	participantID, err := uuid.Parse(chi.URLParam(r, "participantID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var body struct {
		Conditions	[]string	`json:"conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	participant, err := a.queries.UpdateParticipantConditions(r.Context(), db.UpdateParticipantConditionsParams{
		ID: participantID,
		Conditions: body.Conditions,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update conditions")
		return
	}

	respondJSON(w, http.StatusOK, participant)
}

func (a *API) handleUpdateParticipantToggleConcentration(w http.ResponseWriter, r *http.Request) {
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
