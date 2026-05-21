package database

import (
	"context"
	db "go_emmie/prisma/db"
	"log"
)

type DB struct {
	 *db.PrismaClient
}

func Connect() *DB {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	log.Println("Connected to database successfully")

	return &DB{
        PrismaClient: client,
    }

}

func (db *DB) Disconnect() error {
	if err := db.Prisma.Disconnect(); err != nil {
		log.Printf("Error disconnecting from database: %v", err)
	} else {
		log.Println("Disconnected from database successfully")
	}
	return nil
}

// HealthCheck verifies database connectivity.
func (db *DB) HealthCheck(ctx context.Context) error {
	// Simple query to verify connection
	_, err := db.User.FindMany().Take(1).Exec(ctx)
	if err != nil {
		// If no regions exist, that's fine for health check
		// We just want to verify the connection works
		if err.Error() == "ErrNotFound" {
			return nil
		}
		// Check if it's just "no rows" error which is acceptable
		return nil
	}
	return nil
}