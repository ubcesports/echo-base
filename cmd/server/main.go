package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal/database/db"
)

func main() {
	config.LoadEnv(".env")
	db.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/")

}