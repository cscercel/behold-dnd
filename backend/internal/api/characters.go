package api

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/beyond-dnd/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)


func (a *API) handleListCharacters(w http.ResponseWriter, r *http.Request) {
	characters, err := a.queries.ListCharacters(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list characters")
		return
	}

	respondJSON(w, http.StatusOK, characters)
}

func (a *API) handleGetCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))

	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	character, err := a.queries.GetCharacter(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "character not found")
		return
	}

	respondJSON(w, http.StatusOK, character)
}

func (a *API) handleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var params db.CreateCharacterParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	character, err := a.queries.CreateCharacter(r.Context(), params)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to create character")
			return
		}

	respondJSON(w, http.StatusCreated, character)
}

func (a *API) handleUpdateCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}
	
	var params db.UpdateCharacterParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params.ID = id

	character, err := a.queries.UpdateCharacter(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update character")
		return
	}

	respondJSON(w, http.StatusOK, character)
}

func (a *API) handleDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	if err := a.queries.DeleteCharacter(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete character")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *API) handleDamage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var body struct {
		Amount	int	`json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be a positive number")
		return
	}

	character, err := a.characterService.ApplyDamage(r.Context(), id, body.Amount)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to apply damage")
		return
	}

	respondJSON(w, http.StatusOK, character)
}

func (a *API) handleHeal(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var body struct {
		Amount	int	`json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be a positive number")
		return
	}

	character, err := a.characterService.Heal(r.Context(), id, body.Amount)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to heal character")
		return
	}

	respondJSON(w, http.StatusOK, character)
}

func (a *API) handleTempHP(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var body struct {
		Amount	int	`json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "amount must be a positive number")
		return
	}

	character, err := a.characterService.AddTempHP(r.Context(), id, body.Amount)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to add temp HP")
		return
	}

	respondJSON(w, http.StatusOK, character)
}

func (a *API) handleDeathSave(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var body struct {
		Success	bool `json:"success"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "amount must be a positive number")
		return
	}

	character, err := a.characterService.RecordDeathSave(r.Context(), id, body.Success)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to record death save")
		return
	}

	respondJSON(w, http.StatusOK, character)
}

func (a *API) handleLongRest(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	character, err := a.queries.LongRest(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to apply long rest")
		return
	}

	respondJSON(w, http.StatusOK, character)
}

func (a *API) handleShortRest(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var body struct {
		HitDiceRemaining	int	`json:"hit_dice_remaining"`
		CurrentHp			int	`json:"current_hp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	character, err := a.queries.ShortRest(r.Context(), db.ShortRestParams{
		ID:	id,
		HitDiceRemaining: int32(body.HitDiceRemaining),
		CurrentHp: int32(body.CurrentHp),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to apply long rest")
		return
	}

	respondJSON(w, http.StatusOK, character)
}

func (a *API) handleUpdateConditions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var body struct {
		Conditions []string	`json:"conditions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	character, err := a.queries.UpdateConditions(r.Context(), db.UpdateConditionsParams{
		ID:	id,
		Conditions: body.Conditions,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update conditions")
		return
	}

	respondJSON(w, http.StatusOK, character)
}
