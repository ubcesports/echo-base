package config

import (
	"bufio"
	"os"
	"strings"
)

func loadEnv(path string) {
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
		if len(tokens) == 2{
			os.Setenv(tokens[0], tokens[1])
		}
	}
}