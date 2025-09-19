package config

import (
	"bufio"
	"os"
	"strings"
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
		if len(tokens) == 2{
			key := strings.TrimSpace(tokens[0])
			value := strings.TrimSpace(tokens[1])

			if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
			   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}

			os.Setenv(key, value)
		}
	}
}