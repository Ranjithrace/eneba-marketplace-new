package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"eneba-marketplace/backend/internal/api"
	"eneba-marketplace/backend/internal/config"
	"eneba-marketplace/backend/internal/database"
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

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           api.NewRouter(db),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      25 * time.Second,
	}
	log.Printf("API listening on port %s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}
