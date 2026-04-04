package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cscercel/beyond-dnd/internal/db"
	"github.com/cscercel/beyond-dnd/internal/service"
)


type API struct {
	queries				*db.Queries
	pool				*pgxpool.Pool
	characterService	*service.CharacterService
	inventoryService	*service.InventoryService
	spellService		*service.SpellService
}

func New(pool *pgxpool.Pool) *API {
	queries := db.New(pool)
	return &API{
		queries: queries,
		pool: pool,
		characterService: service.NewCharacterService(queries),
		inventoryService: service.NewInventoryService(queries),
		spellService: service.NewSpellService(queries),
	}
}

func (a *API) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", a.handleHealth)

	// Character Routes
	r.Route("/characters", func(r chi.Router) {
		r.Get("/", a.handleListCharacters)
		r.Post("/", a.handleCreateCharacter)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", a.handleGetCharacter)
			r.Put("/", a.handleUpdateCharacter)
			r.Delete("/", a.handleDeleteCharacter)


			// Game mechanics
			r.Post("/damage", a.handleDamage)
			r.Post("/heal", a.handleHeal)
			r.Post("temp-hp", a.handleTempHP)
			r.Post("/death-save", a.handleDeathSave)
			r.Post("long-rest", a.handleLongRest)
			r.Post("short-rest", a.handleShortRest)
			r.Put("/conditions", a.handleUpdateConditions)

			// Inventory
			r.Route("/inventory", func(r chi.Router) {
				r.Get("/", a.handleListInventory)
				r.Post("/", a.handleCreateInventoryItem)
				r.Put("/{itemID}", a.handleUpdateInventoryItem)
				r.Delete("/{itemID}", a.handleDeleteInventoryItem)
				r.Post("/{itemID}/attune", a.handleAttuneItem)
				r.Post("/{itemID}/unattune", a.handleUnattuneItem)
			})

			// Spells
			r.Route("/spells", func(r chi.Router) {
				r.Get("/", a.handleListSpells)
				r.Post("/", a.handleCreateSpell)
				r.Put("/{spellID}", a.handleUpdateSpell)
				r.Delete("/{spellID}", a.handleDeleteSpell)
				r.Post("/{spellID}/toggle-prepared", a.handleToggleSpellPrepared)
			})

			// Spell slots
			r.Route("/spell-slots", func(r chi.Router) {
				r.Get("/", a.handleListSpellSlots)
				r.Put("/", a.handleUpsertSpellSlot)
				r.Post("/use", a.handleUseSpellSlot)
			})
		})
	})

	return r
}
