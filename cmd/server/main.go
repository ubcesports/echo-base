package main

import (
	"github.com/ubcesports/echo-base/config"
)

func main() {
	config.LoadEnv(".env")

}