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
}

func New(pool *pgxpool.Pool) *API {
	return &API{
		queries: db.New(pool),
		pool: pool,
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
			r.Put("/conditions", a.HandleUpdateConditions)
		})
	})

	return r
}
