package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/google/uuid"
)

// @Summary      List all features for a character
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
func (a *API) handleListFeatures(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	if _, err := a.requireCharacterAccess(r, characterID); err != nil {
		respondAccessError(w, err)
		return
	}

	features, err := a.queries.ListFeatures(r.Context(), characterID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list features")
		return
	}

	respondJSON(w, http.StatusOK, features)
}

// @Summary      Create a feature for a character
// @Tags         features
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                true  "Character ID"
// @Param        body body      db.CreateFeatureParams true  "Feature data"
// @Success      201  {object}  db.Feature
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/features [post]
func (a *API) handleCreateFeature(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	if _, err := a.requireCharacterAccess(r ,characterID); err != nil {
		respondAccessError(w, err)
		return
	}

	var params db.CreateFeatureParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params.CharacterID = characterID

	feature, err := a.queries.CreateFeature(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create feature")
		return
	}

	respondJSON(w, http.StatusCreated, feature)
}

// @Summary      Update a feature
// @Tags         features
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      string                true  "Character ID"
// @Param        featureID  path      string                true  "Feature ID"
// @Param        body       body      db.UpdateFeatureParams false "Feature data"
// @Success      200  {object}  db.Feature
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/features/{featureID} [patch]
func (a *API) handleUpdateFeature(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	if _, err := a.requireCharacterAccess(r, characterID); err != nil {
		respondAccessError(w, err)
		return
	}

	featureID, err := uuid.Parse(chi.URLParam(r, "featureID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid feature id")
		return
	}

	var params db.UpdateFeatureParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params.ID = featureID

	feature, err := a.queries.UpdateFeature(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update feature")
		return
	}

	respondJSON(w, http.StatusOK, feature)
}


// @Summary      Delete a feature
// @Tags         features
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      string  true  "Character ID"
// @Param        featureID  path      string  true  "Feature ID"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/features/{featureID} [delete]
func (a *API) handleDeleteFeature(w http.ResponseWriter, r *http.Request) {
	characterID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid character id")
		return
	}

	if _, err := a.requireCharacterAccess(r, characterID); err != nil {
		respondError(w, http.StatusBadRequest, "invalid feature id")
		return
	}

	featureID, err := uuid.Parse(chi.URLParam(r, "featureID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid feature id")
		return
	}

	if err := a.queries.DeleteFeature(r.Context(), featureID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete feature")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
