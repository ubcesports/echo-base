package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ubcesports/echo-base/config"
	"github.com/ubcesports/echo-base/internal/database"
	"github.com/ubcesports/echo-base/internal/handlers"
	"github.com/ubcesports/echo-base/internal/middleware"
)