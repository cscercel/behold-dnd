package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FeatureHandler struct {
	service *service.FeatureService
}

func NewFeatureHandler(service *service.FeatureService) *FeatureHandler {
	return &FeatureHandler{service: service}
}

// RegisterFeatureRoutes mounts feature routes. Call within a router scoped to /characters/{id}/features.
func (h *FeatureHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.handleListFeatures)
	r.Post("/", h.handleCreateFeature)
	r.Patch("/{featureID}", h.handleUpdateFeature)
	r.Delete("/{featureID}", h.handleDeleteFeature)
}

// @Summary      List features for a character
// @Tags         features
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Character ID"
// @Success      200  {array}   db.Feature
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/features [get]
func (h *FeatureHandler) handleListFeatures(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	features, err := h.service.ListFeatures(r.Context(), characterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list features", err)
		return
	}

	respondWithJSON(w, http.StatusOK, features)
}

// @Summary      Add a feature to a character
// @Tags         features
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string               true  "Character ID"
// @Param        body body      db.CreateFeatureParams true  "Feature data"
// @Success      201  {object}  db.Feature
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/features [post]
func (h *FeatureHandler) handleCreateFeature(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	var params db.CreateFeatureParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.CharacterID = characterID

	feature, err := h.service.CreateFeature(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create feature", err)
		return
	}

	respondWithJSON(w, http.StatusOK, feature)
}

// @Summary      Update a feature
// @Tags         features
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string               true  "Character ID"
// @Param        featureID path      string               true  "Feature ID"
// @Param        body    body      db.UpdateFeatureParams false  "Feature data"
// @Success      200  {object}  db.Feature
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/features/{featureID} [patch]
func (h *FeatureHandler) handleUpdateFeature(w http.ResponseWriter, r *http.Request) {
	featureID, err := uuid.Parse(chi.URLParam(r, "featureID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid feature id", err)
		return
	}

	var params db.UpdateFeatureParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.ID = featureID

	feature, err := h.service.UpdateFeature(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update feature", err)
		return
	}

	respondWithJSON(w, http.StatusOK, feature)
}

// @Summary      Delete a feature
// @Tags         features
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string  true  "Character ID"
// @Param        featureID path      string  true  "Feature ID"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/features/{featureID} [delete]
func (h *FeatureHandler) handleDeleteFeature(w http.ResponseWriter, r *http.Request) {
	featureID, err := uuid.Parse(chi.URLParam(r, "featureID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid feature id", err)
		return
	}

	if err := h.service.DeleteFeature(r.Context(), featureID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete feature", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
