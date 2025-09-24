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
	Location *time.Location
}

func LoadConfig() *Config {
	// Load timezone once
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic("failed to load timezone: " + err.Error())
	}

	return &Config{
		Location: loc,
	}
}
