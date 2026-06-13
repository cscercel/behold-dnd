package main


import (
	"context"
	"fmt"
	"log"
	"net/http"


	"github.com/cscercel/behold-dnd/internal/api"
	"github.com/cscercel/behold-dnd/internal/config"
	"github.com/cscercel/behold-dnd/internal/database"

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
	a := api.New(pool, cfg)
	router := a.Routes()

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Behold DnD server listening on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
