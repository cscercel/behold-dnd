package main


import (
	"context"
	"fmt"
	"log"
	"net/http"


	"github.com/cscercel/beyond-dnd/internal/api"
	"github.com/cscercel/beyond-dnd/internal/config"
	"github.com/cscercel/beyond-dnd/internal/database"

)


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
	a := api.New(pool)
	router := a.Routes()

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Beyond DnD server listening on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
