package handlers_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func InitTestDB() *pgxpool.Pool {
	dsn := os.Getenv("EB_DSN")
	if dsn == "" {
		log.Fatal("EB_DSN environment variable must be set")
	}
	dsn = dsn + "&search_path=test"

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v\n", err)
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		log.Fatal("Failed to ping test DB: %v\n", err)
	}

	return dbpool
}

func CloseTestDB(dbpool *pgxpool.Pool) {
	if dbpool != nil {
		dbpool.Close()
	}
}

func BeforeEachTest(t *testing.T, dbpool *pgxpool.Pool) {
	ctx := context.Background()

	// Truncate tables
	_, err := dbpool.Exec(ctx, "TRUNCATE TABLE gamer_activity;")
	require.NoError(t, err)

	_, err = dbpool.Exec(ctx, "TRUNCATE TABLE gamer_profile CASCADE;")
	require.NoError(t, err)

	// Seed test users
	_, err = dbpool.Exec(ctx, `
        INSERT INTO gamer_profile (first_name, last_name, student_number, membership_tier)
        VALUES
        ('John','Doe','11223344',1),
        ('Jane','Doe','87654321',2);
    `)
	require.NoError(t, err)
}

func AfterEachTest(t *testing.T, dbpool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := dbpool.Exec(ctx, "TRUNCATE TABLE gamer_activity;")
	require.NoError(t, err)
}
