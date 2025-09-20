package config

import (
	"bufio"
	"os"
	"strings"
	"time"
)

func LoadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		tokens := strings.SplitN(line, "=", 2)
		if len(tokens) == 2 {
			os.Setenv(tokens[0], tokens[1])
		}
	}
}

type Config struct {
	Schema   string
	Location *time.Location
}

func LoadConfig() *Config {
	// Pick schema based on environment
	schema := os.Getenv("LIVE_SCHEMA")
	if os.Getenv("NODE_ENV") == "test" {
		schema = os.Getenv("TEST_SCHEMA")
	}

	// Load timezone once
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic("failed to load timezone: " + err.Error())
	}

	return &Config{
		Schema:   schema,
		Location: loc,
	}
}
