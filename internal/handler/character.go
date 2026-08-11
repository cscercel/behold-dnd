package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CharacterHandler struct {
	service          *service.CharacterService
	inventoryHandler *InventoryHandler
	spellHandler     *SpellHandler
	featureHandler   *FeatureHandler
}

func NewCharacterHandler(
	service *service.CharacterService,
	inventoryHandler *InventoryHandler,
	spellHandler *SpellHandler,
	featureHandler *FeatureHandler,
) *CharacterHandler {
	return &CharacterHandler{
		service:          service,
		inventoryHandler: inventoryHandler,
		spellHandler:     spellHandler,
		featureHandler:   featureHandler,
	}
}

func (h *CharacterHandler) RegisterRoutes(r chi.Router, authMiddleware, dmOnlyMiddleware func(http.Handler) http.Handler) {
	r.Route("/characters", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/", h.handleListCharacters)
		r.Post("/", h.handleCreateCharacter)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.handleGetCharacter)
			r.Delete("/", h.handleDeleteCharacter)
			r.Patch("/info", h.handleUpdateCharacterInfo)
			r.Patch("/ability-scores", h.handleUpdateCharacterAbilityScores)
			r.Patch("/skills", h.handleUpdateCharacterSkills)
			r.Patch("/level", h.handleUpdateCharacterLevel)
			r.Patch("/training", h.handleUpdateCharacterTraining)
			r.Patch("/currency", h.handleUpdateCharacterCurrency)

			// Game mechanics
			r.Post("/damage", h.handleDamage)
			r.Post("/heal", h.handleHeal)
			r.Post("/temp-hp", h.handleTempHP)
			r.Post("/death-save", h.handleDeathSave)
			r.Post("/long-rest", h.handleLongRest)
			r.Post("/short-rest", h.handleShortRest)
			r.Put("/conditions", h.handleUpdateConditions)

			r.Route("/inventory", h.inventoryHandler.RegisterRoutes)
			r.Route("/spells", h.spellHandler.RegisterSpellRoutes)
			r.Route("/spell-slots", h.spellHandler.RegisterSlotRoutes)
			r.Route("/features", h.featureHandler.RegisterRoutes)
		})
	})

	// DM only routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware, dmOnlyMiddleware)
		r.Get("/players", h.handleListPlayerCharacters)
		r.Get("/npcs", h.handleListNPCs)
	})
}

// @Summary      List all characters
// @Tags         characters
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   db.Character
// @Failure      401  {object}  object{error=string}
// @Router       /characters [get]
func (h *CharacterHandler) handleListCharacters(w http.ResponseWriter, r *http.Request) {
	role, _ := RoleFromContext(r.Context())
	userID, _ := UserIDFromContext(r.Context())

	if role == "dm" {
		characters, err := h.service.ListCharacters(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to list characters", err)
			return
		}
		respondWithJSON(w, http.StatusOK, characters)
		return
	}

	characters, err := h.service.ListUserCharacters(r.Context(), userID)
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
func (h *CharacterHandler) handleListPlayerCharacters(w http.ResponseWriter, r *http.Request) {
	characters, err := h.service.ListPlayerCharacters(r.Context())
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
func (h *CharacterHandler) handleListNPCs(w http.ResponseWriter, r *http.Request) {
	characters, err := h.service.ListNPCs(r.Context())
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
func (h *CharacterHandler) handleGetCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	character, err := h.requireCharacterAccess(r, id)
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
func (h *CharacterHandler) handleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var params db.CreateCharacterParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "missing user id", nil)
		return
	}
	params.OwnerID = userID

	// TEXT[] columns are NOT NULL with a '{}' default, but an explicit NULL
	if params.TrainingArmor == nil {
		params.TrainingArmor = []string{}
	}
	if params.TrainingWeapons == nil {
		params.TrainingWeapons = []string{}
	}
	if params.TrainingTools == nil {
		params.TrainingTools = []string{}
	}
	if params.TrainingLanguages == nil {
		params.TrainingLanguages = []string{}
	}
	if params.Conditions == nil {
		params.Conditions = []string{}
	}
	if params.Resistances == nil {
		params.Resistances = []string{}
	}
	if params.Vulnerabilities == nil {
		params.Vulnerabilities = []string{}
	}
	if params.Immunities == nil {
		params.Immunities = []string{}
	}

	character, err := h.service.CreateCharacter(r.Context(), params)
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
func (h *CharacterHandler) handleUpdateCharacterInfo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterInfoParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.ID = id

	character, err := h.service.UpdateCharacterInfo(r.Context(), params)
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
func (h *CharacterHandler) handleUpdateCharacterAbilityScores(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterAbilityScoresParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.ID = id

	character, err := h.service.UpdateCharacterAbilityScores(r.Context(), params)
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
func (h *CharacterHandler) handleUpdateCharacterSkills(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterSkillsParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.ID = id

	character, err := h.service.UpdateCharacterSkills(r.Context(), params)
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
func (h *CharacterHandler) handleUpdateCharacterLevel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterLevelParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.ID = id

	character, err := h.service.UpdateCharacterLevel(r.Context(), params)
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
func (h *CharacterHandler) handleUpdateCharacterTraining(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterTrainingParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.ID = id

	character, err := h.service.UpdateCharacterTraining(r.Context(), params)
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
// @Router       /characters/{id}/currency [patch]
func (h *CharacterHandler) handleUpdateCharacterCurrency(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var params db.UpdateCharacterCurrencyParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	params.ID = id

	character, err := h.service.UpdateCharacterCurrency(r.Context(), params)
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
func (h *CharacterHandler) handleDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	if err := h.service.DeleteCharacter(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete character", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
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
func (h *CharacterHandler) handleDamage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
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

	character, err := h.service.ApplyDamage(r.Context(), id, body.Amount)
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
func (h *CharacterHandler) handleHeal(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
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

	character, err := h.service.Heal(r.Context(), id, body.Amount)
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
func (h *CharacterHandler) handleTempHP(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
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

	character, err := h.service.AddTempHP(r.Context(), id, body.Amount)
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
func (h *CharacterHandler) handleDeathSave(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	var body struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	character, err := h.service.RecordDeathSave(r.Context(), id, body.Success)
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
func (h *CharacterHandler) handleLongRest(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
		respondWithError(w, http.StatusForbidden, "you do not own this character", err)
		return
	}

	character, err := h.service.LongRest(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to apply long rest", err)
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
// @Param        body body      object{hit_dice_remaining=int,current_hp=int} true  "Short rest details"
// @Success      200  {object}  db.Character
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /characters/{id}/short-rest [post]
func (h *CharacterHandler) handleShortRest(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
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

	character, err := h.service.ShortRest(r.Context(), id, body.HitDiceRemaining, body.CurrentHp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to apply short rest", err)
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
func (h *CharacterHandler) handleUpdateConditions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid character id", err)
		return
	}

	if _, err := h.requireCharacterAccess(r, id); err != nil {
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

	character, err := h.service.UpdateConditions(r.Context(), id, body.Conditions)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update conditions", err)
		return
	}

	respondWithJSON(w, http.StatusOK, character)
}

// requireCharacterAccess checks the requester either is the DM or owns the character.
func (h *CharacterHandler) requireCharacterAccess(r *http.Request, characterID uuid.UUID) (db.GetCharacterRow, error) {
	character, err := h.service.GetCharacter(r.Context(), characterID)
	if err != nil {
		return db.GetCharacterRow{}, fmt.Errorf("not found")
	}

	role, _ := RoleFromContext(r.Context())
	if role == "dm" {
		return character, nil
	}

	userID, _ := UserIDFromContext(r.Context())
	if character.OwnerID != userID {
		return db.GetCharacterRow{}, fmt.Errorf("forbidden")
	}

	return character, nil
}
