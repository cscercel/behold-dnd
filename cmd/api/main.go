package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/cscercel/behold-dnd/internal/config"
	"github.com/cscercel/behold-dnd/internal/database"
	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/cscercel/behold-dnd/internal/handler"
	"github.com/cscercel/behold-dnd/internal/service"

	_ "github.com/cscercel/behold-dnd/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Behold D&D API
// @version         1.0.0
// @description     API for managing your D&D campaign characters, inventory, spells and combat.

// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by your JWT token
func main() {

	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	// Connect to Database
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Instantiate API
	queries := db.New(pool)

	// Services
	authService := service.NewAuthService(queries, cfg.JWTSecret, cfg.JWTExpiryHours, cfg.RegistrationCode)
	spellService := service.NewSpellService(queries)
	characterService := service.NewCharacterService(queries, spellService)
	inventoryService := service.NewInventoryService(queries)
	combatService := service.NewCombatService(queries)

	// Handlers
	userHandler := handler.NewUserHandler(authService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	spellHandler := handler.NewSpellHandler(spellService)
	characterHandler := handler.NewCharacterHandler(characterService, inventoryHandler, spellHandler)
	combatHandler := handler.NewCombatHandler(combatService)

	// Router
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"),
	))

	// Public routes

	// Small `Mandatory` test route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	userHandler.RegisterPublicRoutes(r)

	// Local middleware
	authMiddleware := handler.AuthenticateMiddleware(authService)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		userHandler.RegisterProtectedRoutes(r)
	})

	characterHandler.RegisterRoutes(r, authMiddleware, handler.RequireDM)
	combatHandler.RegisterRoutes(r, authMiddleware, handler.RequireDM)

	// Start server
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	log.Printf("server running on port: %v", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
