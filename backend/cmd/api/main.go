package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"eneba-marketplace/backend/internal/api"
	"eneba-marketplace/backend/internal/config"
	"eneba-marketplace/backend/internal/database"
	"eneba-marketplace/backend/internal/web"
)

func main() {
	cfg, err := config.FromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler := api.NewRouter(db)

	// Optional: serve the compiled React application from Go.
	// If unset, the existing API-only development mode is preserved.
	if distDir := os.Getenv("WEB_DIST_DIR"); distDir != "" {
		index := filepath.Join(distDir, "index.html")

		if _, err := os.Stat(index); err != nil {
			log.Fatalf("React build not found at %s: %v", index, err)
		}

		handler = web.NewHandler(handler, distDir)
		log.Printf("Serving React application from %s", distDir)
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      25 * time.Second,
	}

	log.Printf("API listening on port %s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}
