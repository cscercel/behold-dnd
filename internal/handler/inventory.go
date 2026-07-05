package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// RegisterRoutes mounts inventory routes. Call within a router already scoped to /characters/{id}/inventory.
func (h *InventoryHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.handleListInventory)
	r.Post("/", h.handleCreateInventoryItem)
	r.Patch("/{itemID}", h.handleUpdateInventoryItem)
	r.Delete("/{itemID}", h.handleDeleteInventoryItem)
	r.Post("/{itemID}/attune", h.handleAttuneItem)
	r.Post("/{itemID}/unattune", h.handleUnattuneItem)
}

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
func (h *InventoryHandler) handleListInventory(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	items, err := h.service.ListItems(r.Context(), characterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list inventory", err)
		return
	}

	respondWithJSON(w, http.StatusOK, items)
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
func (h *InventoryHandler) handleCreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var params db.CreateInventoryItemParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.CharacterID = characterID

	item, err := h.service.CreateItem(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create inventory item", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, item)
}

// @Summary      Update an inventory item
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string                      true  "Character ID"
// @Param        itemID path      string                      true  "Item ID"
// @Param        body   body      db.UpdateInventoryItemParams false  "Item data"
// @Success      200  {object}  db.InventoryItem
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/inventory/{itemID} [patch]
func (h *InventoryHandler) handleUpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid item id", err)
		return
	}

	var params db.UpdateInventoryItemParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.ID = itemID

	item, err := h.service.UpdateItem(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update inventory item", err)
		return
	}

	respondWithJSON(w, http.StatusOK, item)
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
func (h *InventoryHandler) handleDeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid item id", err)
		return
	}

	if err := h.service.DeleteItem(r.Context(), itemID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete item", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
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
func (h *InventoryHandler) handleAttuneItem(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid item id", err)
		return
	}

	item, err := h.service.AttuneItem(r.Context(), characterID, itemID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to attune item", err)
		return
	}

	respondWithJSON(w, http.StatusOK, item)
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
func (h *InventoryHandler) handleUnattuneItem(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid item id", err)
		return
	}

	item, err := h.service.UnattuneItem(r.Context(), characterID, itemID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to unattune item", err)
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}
