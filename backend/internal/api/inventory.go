package api

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      List inventory items for a character
// @Tags         inventory
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Character ID"
// @Success      200  {array}   db.InventoryItem
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/inventory [get]
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

// @Summary      Add an item to a character's inventory
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                      true  "Character ID"
// @Param        body body      db.CreateInventoryItemParams true  "Item data"
// @Success      201  {object}  db.InventoryItem
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/inventory [post]
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

// @Summary      Update an inventory item
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                      true  "Character ID"
// @Param        itemID path      string                      true  "Item ID"
// @Param        body   body      db.UpdateInventoryItemParams true  "Item data"
// @Success      200  {object}  db.InventoryItem
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/inventory/{itemID} [patch]
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

// @Summary      Delete an inventory item
// @Tags         inventory
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string  true  "Character ID"
// @Param        itemID path      string  true  "Item ID"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/inventory/{itemID} [delete]
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

// @Summary      Attune to a magic item
// @Tags         inventory
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string  true  "Character ID"
// @Param        itemID path      string  true  "Item ID"
// @Success      200  {object}  db.InventoryItem
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /characters/{id}/inventory/{itemID}/attune [post]
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


// @Summary      Remove attunement from a magic item
// @Tags         inventory
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string  true  "Character ID"
// @Param        itemID path      string  true  "Item ID"
// @Success      200  {object}  db.InventoryItem
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /characters/{id}/inventory/{itemID}/unattune [post]
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
