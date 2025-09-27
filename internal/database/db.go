package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var DB *pgxpool.Pool

func Init() {
	dsn := os.Getenv("EB_DSN")
	if dsn == "" {
		log.Fatal("EB_DSN environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	if err = DB.Ping(ctx); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Connected to database (pgxpool)")
}

// Ping checks if the database connection is still alive
func Ping() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}
	return DB.Ping(context.Background())
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		DB.Close()
	}
	return nil
}
