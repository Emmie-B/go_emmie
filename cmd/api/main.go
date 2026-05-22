package main

import (
	"fmt"
	"go_emmie/internal/config"
	"go_emmie/internal/database"
	"go_emmie/internal/routes"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 💥 THE FIX: Explicitly inject the string into the environment for Prisma's engine
	os.Setenv("DATABASE_URL", cfg.Database.URL)

	dbClient := database.Connect()

	defer func() {
		if err := dbClient.Disconnect(); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	r := routes.New(cfg, dbClient)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server running on port %d", cfg.Server.Port)
	log.Printf("🚀  server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}
