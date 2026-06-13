package api

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)


// @Summary      List spells for a character
// @Tags         spells
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Character ID"
// @Success      200  {array}   db.Spell
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/spells [get]
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

// @Summary      Add a spell to a character
// @Tags         spells
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string               true  "Character ID"
// @Param        body body      db.CreateSpellParams true  "Spell data"
// @Success      201  {object}  db.Spell
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/spells [post]
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

// @Summary      Update a spell
// @Tags         spells
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string               true  "Character ID"
// @Param        spellID path      string               true  "Spell ID"
// @Param        body    body      db.UpdateSpellParams false  "Spell data"
// @Success      200  {object}  db.Spell
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/spells/{spellID} [patch]
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

// @Summary      Delete a spell
// @Tags         spells
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string  true  "Character ID"
// @Param        spellID path      string  true  "Spell ID"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/spells/{spellID} [delete]
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

// @Summary      Toggle a spell's prepared status
// @Tags         spells
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string  true  "Character ID"
// @Param        spellID path      string  true  "Spell ID"
// @Success      200  {object}  db.Spell
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/spells/{spellID}/toggle-prepared [post]
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

// @Summary      List spell slots for a character
// @Tags         spells
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Character ID"
// @Success      200  {array}   db.SpellSlot
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/spell-slots [get]
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

// @Summary      Create or update a spell slot level
// @Tags         spells
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                    true  "Character ID"
// @Param        body body      db.UpsertSpellSlotParams  true  "Spell slot data"
// @Success      200  {object}  db.SpellSlot
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/spell-slots [put]
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


// @Summary      Use a spell slot
// @Tags         spells
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                    true  "Character ID"
// @Param        body body      object{spell_level=int}   true  "Spell level to use"
// @Success      200  {object}  db.SpellSlot
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /characters/{id}/spell-slots/use [post]
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
