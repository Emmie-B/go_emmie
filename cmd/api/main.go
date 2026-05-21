package main

import (
	"fmt"
	"go_emmie/internal/config"
	"go_emmie/internal/database"
	"go_emmie/internal/routes"
	"log"
	"os"
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

	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to run server on port %d: %v", cfg.Server.Port, err)
	}

	log.Printf("Server running on port %d", cfg.Server.Port)

	// router.Run(":" + fmt.Sprint(cfg.Server.Port))
}
