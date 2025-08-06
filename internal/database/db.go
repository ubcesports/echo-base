package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SCHEMA"),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}

	log.Println("Connected to database")
}

// Ping checks if the database connection is still alive
func Ping() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}
	return DB.Ping()
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
