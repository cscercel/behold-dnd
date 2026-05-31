package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/cscercel/behold-dnd/internal/db" // ONLY required for Swagger to pick up db interfaces
	"github.com/cscercel/behold-dnd/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CharacterHandler struct {
	service	*service.CharacterService
}

func NewCharacterHandler(service *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{service: service}
}

func (h *CharacterHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/characters", func(r chi.Router) {
		// Public routes

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
		})
	})
}

func (h *CharacterHandler)
