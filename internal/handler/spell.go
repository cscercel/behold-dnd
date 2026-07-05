package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SpellHandler struct {
	service *service.SpellService
}

func NewSpellHandler(service *service.SpellService) *SpellHandler {
	return &SpellHandler{service: service}
}

// RegisterSpellRoutes mounts spell routes. Call within a router scoped to /characters/{id}/spells.
func (h *SpellHandler) RegisterSpellRoutes(r chi.Router) {
	r.Get("/", h.handleListSpells)
	r.Post("/", h.handleCreateSpell)
	r.Patch("/{spellID}", h.handleUpdateSpell)
	r.Delete("/{spellID}", h.handleDeleteSpell)
	r.Post("/{spellID}/toggle-prepared", h.handleToggleSpellPrepared)
}

// RegisterSlotRoutes mounts spell slot routes. Call within a router scoped to /characters/{id}/spell-slots.
func (h *SpellHandler) RegisterSlotRoutes(r chi.Router) {
	r.Get("/", h.handleListSpellSlots)
	r.Put("/", h.handleUpsertSpellSlot)
	r.Post("/use", h.handleUseSpellSlot)
}

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
func (h *SpellHandler) handleListSpells(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	spells, err := h.service.ListSpells(r.Context(), characterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list spells", err)
		return
	}

	respondWithJSON(w, http.StatusOK, spells)
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
func (h *SpellHandler) handleCreateSpell(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var params db.CreateSpellParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.CharacterID = characterID

	spell, err := h.service.CreateSpell(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create spell", err)
		return
	}

	respondWithJSON(w, http.StatusOK, spell)
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
func (h *SpellHandler) handleUpdateSpell(w http.ResponseWriter, r *http.Request) {
	spellID, err := uuid.Parse(chi.URLParam(r, "spellID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid spell id", err)
		return
	}

	var params db.UpdateSpellParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.ID = spellID

	spell, err := h.service.UpdateSpell(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update spell", err)
		return
	}

	respondWithJSON(w, http.StatusOK, spell)
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
func (h *SpellHandler) handleDeleteSpell(w http.ResponseWriter, r *http.Request) {
	spellID, err := uuid.Parse(chi.URLParam(r, "spellID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid spell id", err)
		return
	}

	if err := h.service.DeleteSpell(r.Context(), spellID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete spell", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
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
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/spells/{spellID}/toggle-prepared [post]
func (h *SpellHandler) handleToggleSpellPrepared(w http.ResponseWriter, r *http.Request) {
	spellID, err := uuid.Parse(chi.URLParam(r, "spellID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid spell id", err)
		return
	}

	spell, err := h.service.ToggleSpellPrepared(r.Context(), spellID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to toggle spell preparation", err)
		return
	}

	respondWithJSON(w, http.StatusOK, spell)
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
func (h *SpellHandler) handleListSpellSlots(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	slots, err := h.service.ListSpellSlots(r.Context(), characterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list spell slots", err)
		return
	}

	respondWithJSON(w, http.StatusOK, slots)
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
func (h *SpellHandler) handleUpsertSpellSlot(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var params db.UpsertSpellSlotParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.CharacterID = characterID

	slot, err := h.service.UpsertSpellSlot(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to upsert spell slot", err)
		return
	}

	respondWithJSON(w, http.StatusOK, slot)
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
func (h *SpellHandler) handleUseSpellSlot(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var params db.UseSpellSlotParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.CharacterID = characterID

	slot, err := h.service.UseSpellSlot(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to use spell slot", err)
		return
	}

	respondWithJSON(w, http.StatusOK, slot)
}
