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
	fmt.Printf("Loaded config: %+v\n", cfg)

	dbClient := database.Connect()

	if err != nil {
		log.Printf("Error connecting to database: %s", err)
	}
	
	defer func() {
		if err := dbClient.Disconnect(); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	defer dbClient.Disconnect()

	r := routes.New(cfg, dbClient)

	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))

	log.Printf("Server running on port %s", cfg.Server.Port)

	// router.Run(":" + fmt.Sprint(cfg.Server.Port))
}
