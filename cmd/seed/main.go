package main

import (
	"context"
	"os"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal/database"
	"github.com/ubcesports/echo-base/internal/services"
)

func main() {
	config.LoadEnv(".env")
	database.Init()
	defer database.Close()

	if len(os.Args) < 3 {
		println("please specify operation and value")
		os.Exit(1)
	}

	if os.Args[1] != "apikey" {
		println("operation not supported")
		os.Exit(1)
	}

	authRepo := database.NewAuthRepository(database.DB)
	authService := services.NewAuthService(authRepo)

	response, err := authService.GenerateAPIKey(context.Background(), os.Args[2])

	if err != nil {
		println("error while generating api key:", err.Error())
		os.Exit(1)
	}

	println("generated api key!")
	println("key id:", response.KeyId)
	println("token:", response.APIKey)
}
