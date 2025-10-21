package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal"
	"github.com/ubcesports/echo-base/internal/database"
	"github.com/ubcesports/echo-base/internal/services"
)

const TIMEOUT = 10

func main() {
	if err := bootstrap(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func bootstrap(ctx context.Context, args []string) error {
	config.LoadEnv(".env")
	database.Init()

	// Initialize repositories
	authRepo := database.NewAuthRepository(database.DB)

	// Initialize services
	authService := services.NewAuthService(authRepo)

	// Initialize server
	srv := internal.NewServer(authService)

	httpServer := &http.Server{
		Addr:    ":" + os.Getenv("EB_PORT"),
		Handler: srv,
	}

	// Run server in its own goroutine
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Shutdown Gracefully
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, TIMEOUT*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
		if err := database.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing database: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
