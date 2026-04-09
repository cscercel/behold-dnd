package api

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/beyond-dnd/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)


func (a *API) handleListSpells(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	spells, err := a.queries.ListSpells(r.Context(), characterID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list spells")
		return
	}

	respondJSON(w, http.StatusOK, spells)
}

func (a *API) handleCreateSpell(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var params db.CreateSpellParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Assign spell to character
	params.CharacterID = characterID

	spell, err := a.queries.CreateSpell(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create spell")
		return
	}

	respondJSON(w, http.StatusOK, spell)
}

func (a *API) handleUpdateSpell(w http.ResponseWriter, r *http.Request) {
	spellID, err := uuid.Parse(chi.URLParam(r, "spellID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid spell id")
		return
	}
	
	var params db.UpdateSpellParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params.ID = spellID

	spell, err := a.queries.UpdateSpell(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update inventory")
		return
	}

	respondJSON(w, http.StatusOK, spell)
}

func (a *API) handleDeleteSpell(w http.ResponseWriter, r *http.Request) {
	spellID, err := uuid.Parse(chi.URLParam(r, "spellID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid spell item")
		return
	}

	if err := a.queries.DeleteSpell(r.Context(), spellID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete spell")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *API) handleToggleSpellPrepared(w http.ResponseWriter, r *http.Request) {
	spellID, err := uuid.Parse(chi.URLParam(r, "spellID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid spell id")
		return
	}

	spell, err := a.queries.ToggleSpellPrepared(r.Context(), spellID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to toggle spell preparation")
		return
	}

	respondJSON(w, http.StatusOK, spell)
}

func (a *API) handleListSpellSlots(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	slots, err := a.queries.ListSpellSlots(r.Context(), characterID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list spell slots")
		return
	}
	
	respondJSON(w, http.StatusOK, slots)
}

func (a *API) handleUpsertSpellSlot(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var params db.UpsertSpellSlotParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params.CharacterID = characterID

	slot, err := a.queries.UpsertSpellSlot(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to upsert spell slots")
		return
	}

	respondJSON(w, http.StatusOK, slot)
}

func (a *API) handleUseSpellSlot(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var params db.UseSpellSlotParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params.CharacterID = characterID

	slot, err := a.queries.UseSpellSlot(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusBadRequest, "failed to use spell slot")
		return
	}

	respondJSON(w, http.StatusOK, slot)
}
