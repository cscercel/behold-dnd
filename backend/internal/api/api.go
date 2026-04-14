package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cscercel/behold-dnd/internal/config"
	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/cscercel/behold-dnd/internal/service"
	appMiddleware "github.com/cscercel/behold-dnd/internal/middleware"

	_ "github.com/cscercel/behold-dnd/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)


type API struct {
	queries				*db.Queries
	pool				*pgxpool.Pool
	authService			*service.AuthService
	characterService	*service.CharacterService
	inventoryService	*service.InventoryService
	spellService		*service.SpellService
	combatService		*service.CombatService
}

func New(pool *pgxpool.Pool, cfg *config.Config) *API {
	queries := db.New(pool)
	return &API{
		queries: queries,
		pool: pool,
		authService: service.NewAuthService(queries, cfg.JWTSecret, cfg.JWTExpiryHours, cfg.RegistrationCode),
		characterService: service.NewCharacterService(queries),
		inventoryService: service.NewInventoryService(queries),
		spellService: service.NewSpellService(queries),
		combatService: service.NewCombatService(queries),
	}
}

func (a *API) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, 
	}))

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"),
	))

	// Public routes
	r.Get("/health", a.handleHealth)
	r.Post("/auth/register", a.handleRegister)
	r.Post("/auth/login", a.handleLogin)


	// Protected routes
	r.Group(func(r chi.Router) {

		r.Use(appMiddleware.Authenticate(a.authService))

		r.Get("/auth/me", a.handleMe)

		// Character Routes
		r.Route("/characters", func(r chi.Router) {
			r.Post("/", a.handleCreateCharacter)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", a.handleGetCharacter)
				r.Patch("/", a.handleUpdateCharacter)
				r.Delete("/", a.handleDeleteCharacter)


				// Game mechanics
				r.Post("/damage", a.handleDamage)
				r.Post("/heal", a.handleHeal)
				r.Post("/temp-hp", a.handleTempHP)
				r.Post("/death-save", a.handleDeathSave)
				r.Post("/long-rest", a.handleLongRest)
				r.Post("/short-rest", a.handleShortRest)
				r.Put("/conditions", a.handleUpdateConditions)

				// Inventory
				r.Route("/inventory", func(r chi.Router) {
					r.Get("/", a.handleListInventory)
					r.Post("/", a.handleCreateInventoryItem)
					r.Patch("/{itemID}", a.handleUpdateInventoryItem)
					r.Delete("/{itemID}", a.handleDeleteInventoryItem)
					r.Post("/{itemID}/attune", a.handleAttuneItem)
					r.Post("/{itemID}/unattune", a.handleUnattuneItem)
				})

				// Spells
				r.Route("/spells", func(r chi.Router) {
					r.Get("/", a.handleListSpells)
					r.Post("/", a.handleCreateSpell)
					r.Patch("/{spellID}", a.handleUpdateSpell)
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

		// DM only routes
		r.Group(func(r chi.Router) {

			r.Use(appMiddleware.RequireDM)
			
			// List all characters
			r.Route("/list-characters", func(r chi.Router) {
				r.Get("/", a.handleListCharacters)
			})

			// Combat Routes
			r.Route("/combat", func(r chi.Router) {
				r.Get("/", a.handleListEncounters)
				r.Post("/", a.handleCreateEncounter)
				r.Get("/active", a.handleGetActiveEncounter)

				r.Route("/{encounterID}", func(r chi.Router) {
					r.Get("/", a.handleGetEncounter)
					r.Delete("/", a.handleDeleteEncounter)
					r.Post("/start", a.handleStartEncounter)
					r.Post("/end", a.handleEndEncounter)
					r.Post("/next-round", a.handleNextRound)

					// Participants
					r.Get("/participants", a.handleListParticipants)
					r.Post("/participants", a.handleAddParticipant)

					r.Route("/participants/{participantID}", func(r chi.Router) {
						r.Delete("/", a.handleRemoveParticipant)
						r.Post("/damage", a.handleParticipantDamage)
						r.Post("/heal", a.handleParticipantHeal)
						r.Post("/temp-hp", a.handleParticipantTempHP)
						r.Put("/initiative", a.handleParticipantInitiative)
						r.Put("/conditions", a.handleParticipantConditions)
						r.Post("/toggle-concentration", a.handleParticipantToggleConcentration)
						r.Post("/deactivate", a.handleDeactivateParticipant)
					})
				})
			})
		})
	})

	return r
}
