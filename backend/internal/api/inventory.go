package api

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/beyond-dnd/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)


func (a *API) handleListInventory(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	items, err := a.queries.ListInventoryItems(r.Context(), characterID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list inventory")
		return
	}

	respondJSON(w, http.StatusOK, items)
}

func (a *API) handleCreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	var params db.CreateInventoryItemParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	// Assign item to character
	params.CharacterID = characterID

	item, err := a.queries.CreateInventoryItem(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create inventory")
		return
	}

	respondJSON(w, http.StatusCreated, item)
}

func (a *API) handleUpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var params db.UpdateInventoryItemParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	params.ID = itemID

	item, err := a.queries.UpdateInventoryItem(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update inventory")
		return
	}

	respondJSON(w, http.StatusOK, item)
}

func (a *API) handleDeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := a.queries.DeleteInventoryItem(r.Context(), itemID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete item")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *API) handleAttuneItem(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := a.inventoryService.AttuneItem(r.Context(), characterID, itemID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, item)
}

func (a *API) handleUnattuneItem(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := a.inventoryService.UnattuneItem(r.Context(), characterID, itemID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, item)
}
