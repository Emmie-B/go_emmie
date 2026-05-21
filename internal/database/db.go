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

func (d *DB) Disconnect() error {
	if err := d.Prisma.Disconnect(); err != nil {
		log.Printf("Error disconnecting from database: %v", err)
		return err
	}

	log.Println("Disconnected from database successfully")
	return nil
}

// HealthCheck verifies database connectivity.
func (d *DB) HealthCheck(ctx context.Context) error {
	// Simple query to verify connection
	_, err := d.User.FindMany().Take(1).Exec(ctx)
	if err != nil {
		// If the generated Prisma client returns ErrNotFound, treat as healthy
		if db.IsErrNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}
