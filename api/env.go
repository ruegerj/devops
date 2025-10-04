package main

import (
	"fmt"
	"os"
)

func resolveEnvVar(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("missing env var: %s", key)
	}

	return value, nil
}
