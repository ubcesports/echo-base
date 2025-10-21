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
)

func main() {
	if err := run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	config.LoadEnv(".env")
	database.Init()

	srv := internal.NewServer()

	httpServer := &http.Server{
		Addr:    ":" + os.Getenv("EB_PORT"),
		Handler: srv,
	}

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
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 11*time.Second)
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
