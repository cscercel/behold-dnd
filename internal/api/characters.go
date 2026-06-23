package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/cscercel/behold-dnd/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      List all characters
// @Tags         characters
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   db.Character
// @Failure      401  {object}  object{error=string}
// @Router       /characters [get]
func (a *API) handleListCharacters(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.RoleFromContext(r.Context())
	userID, _ := middleware.UserIDFromContext(r.Context())

	if role == "dm" {
		// DM should see all characters
		characters, err := a.queries.ListCharacters(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to list characters", err)
			return
		}

		respondWithJSON(w, http.StatusOK, characters)
		return
	}

	// Players only see their own characters
	characters, err := a.queries.ListUserCharacters(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list characters", err)
		return
	}

	respondWithJSON(w, http.StatusOK, characters)
}

// @Summary      List all player characters (DM only)
// @Tags         characters
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   db.Character
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/players [get]
func (a *API) handleListPlayerCharacters(w http.ResponseWriter, r *http.Request) {
	characters, err := a.queries.ListPlayerCharacters(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list player characters", err)
		return
	}

	respondWithJSON(w, http.StatusOK, characters)
}

// @Summary      List all NPCs (DM only)
// @Tags         characters
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   db.Character
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/npcs [get]
func (a *API) handleListNPCs(w http.ResponseWriter, r *http.Request) {
	characters, err := a.queries.ListNPCs(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to list NPCs", err)
		return
	}

	respondWithJSON(w, http.StatusOK, characters)
}

// @Summary      Get a character
// @Tags         characters
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Character ID"
// @Success      200  {object}  db.GetCharacterRow
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /characters/{id} [get]
func (a *API) handleGetCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	character, err := a.requireCharacterAccess(r, id)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Create a character
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body      db.CreateCharacterParams true "Character data"
// @Success      201  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /characters [post]
func (a *API) handleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var params db.CreateCharacterParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := a.queries.CreateCharacter(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create character", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, character)
}

// @Summary      Update a character's info
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                   true  "Character ID"
// @Param        body body      db.UpdateCharacterInfoParams false  "Character data"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/info [patch]
func (a *API) handleUpdateCharacterInfo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterInfoParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	params.ID = id

	character, err := a.queries.UpdateCharacterInfo(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Update a character's ability scores
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                   true  "Character ID"
// @Param        body body      db.UpdateCharacterAbilityScoresParams false  "Character data"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/ability-scores [patch]
func (a *API) handleUpdateCharacterAbilityScores(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterAbilityScoresParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	params.ID = id

	character, err := a.queries.UpdateCharacterAbilityScores(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Update a character's skills
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                   true  "Character ID"
// @Param        body body      db.UpdateCharacterSkillsParams false  "Character data"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/skills [patch]
func (a *API) handleUpdateCharacterSkills(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterSkillsParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	params.ID = id

	character, err := a.queries.UpdateCharacterSkills(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Update a character's level
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                   true  "Character ID"
// @Param        body body      db.UpdateCharacterLevelParams false  "Character data"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/level [patch]
func (a *API) handleUpdateCharacterLevel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterLevelParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	params.ID = id

	character, err := a.queries.UpdateCharacterLevel(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Update a character's training
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                   true  "Character ID"
// @Param        body body      db.UpdateCharacterTrainingParams false  "Character data"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/training [patch]
func (a *API) handleUpdateCharacterTraining(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterTrainingParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	params.ID = id

	character, err := a.queries.UpdateCharacterTraining(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Update a character's currency
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                   true  "Character ID"
// @Param        body body      db.UpdateCharacterCurrencyParams false  "Character data"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/skills [patch]
func (a *API) handleUpdateCharacterCurrency(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterCurrencyParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	params.ID = id

	character, err := a.queries.UpdateCharacterCurrency(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Delete a character
// @Tags         characters
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Character ID"
// @Success      204
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id} [delete]
func (a *API) handleDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	if err := a.queries.DeleteCharacter(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete character", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Deal damage to a character
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string             true  "Character ID"
// @Param        body body      object{amount=int} true  "Damage amount"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/damage [post]
func (a *API) handleDamage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var body struct {
		Amount int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount must be a positive number", err)
		return
	}

	character, err := a.characterService.ApplyDamage(r.Context(), id, body.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to apply damage", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Heal a character
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string             true  "Character ID"
// @Param        body body      object{amount=int} true  "Heal amount"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/heal [post]
func (a *API) handleHeal(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var body struct {
		Amount int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount must be a positive number", err)
		return
	}

	character, err := a.characterService.Heal(r.Context(), id, body.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to heal character", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Add temporary HP to a character
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string             true  "Character ID"
// @Param        body body      object{amount=int} true  "Temp HP amount"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/temp-hp [post]
func (a *API) handleTempHP(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var body struct {
		Amount int `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount must be a positive number", err)
		return
	}

	character, err := a.characterService.AddTempHP(r.Context(), id, body.Amount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to add temp HP", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Record a death saving throw
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string               true  "Character ID"
// @Param        body body      object{success=bool} true  "Death save result"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/death-save [post]
func (a *API) handleDeathSave(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var body struct {
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "amount must be a positive number", err)
		return
	}

	character, err := a.characterService.RecordDeathSave(r.Context(), id, body.Success)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to record death save", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Long rest — restores HP, hit dice, resets death saves, conditions and spell slots
// @Tags         characters
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Character ID"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/long-rest [post]
func (a *API) handleLongRest(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	character, err := a.queries.LongRest(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to apply long rest", err)
		return
	}

	// Reset spell slots
	if err := a.spellService.LongRestSlots(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to reset spell slots", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Short rest — spend hit dice to regain HP
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                                        true  "Character ID"
// @Param        body body      object{hit_dice_used=int,hp_regained=int}     true  "Short rest details"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/short-rest [post]
func (a *API) handleShortRest(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var body struct {
		HitDiceRemaining int `json:"hit_dice_remaining"`
		CurrentHp        int `json:"current_hp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := a.queries.ShortRest(r.Context(), db.ShortRestParams{
		ID:               id,
		HitDiceRemaining: int32(body.HitDiceRemaining),
		CurrentHp:        int32(body.CurrentHp),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to apply long rest", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// @Summary      Update active conditions on a character
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                       true  "Character ID"
// @Param        body body      object{conditions=[]string}  true  "Conditions list"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/conditions [put]
func (a *API) handleUpdateConditions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := a.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var body struct {
		Conditions []string `json:"conditions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := a.queries.UpdateConditions(r.Context(), db.UpdateConditionsParams{
		ID:         id,
		Conditions: body.Conditions,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update conditions", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// Function to authenticate character ownership
func (a *API) requireCharacterAccess(r *http.Request, characterID uuid.UUID) (db.GetCharacterRow, error) {
	character, err := a.queries.GetCharacter(r.Context(), characterID)
	if err != nil {
		return db.GetCharacterRow{}, fmt.Errorf("not found")
	}

	role, _ := middleware.RoleFromContext(r.Context())
	if role == "dm" {
		return character, nil
	}

	userID, _ := middleware.UserIDFromContext(r.Context())
	if character.OwnerID != userID {
		return db.GetCharacterRow{}, fmt.Errorf("forbidden")
	}

	return character, nil
}
