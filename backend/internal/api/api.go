package api

import (
	"github.com/cscercel/beyond-dnd/internal/db"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)


type API struct {
	queries	*db.Queries
	pool	*pgxpool.Pool
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

	r.Route("/characters", func(r chi.Router) {
		r.Get("/", a.handleListCharacters)
		r.Post("/", a.handleCreateCharacter)
		r.Get("/{id}", a.handleGetCharacter)
	})

	return r
}
