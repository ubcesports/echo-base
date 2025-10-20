package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ubcesports/echo-base/config"
)

func main() {
	if err := run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	config.LoadEnv(".env")
	return nil
}
