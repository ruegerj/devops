package main

import (
	"errors"
	"fmt"
	"os"
)

func resolveEnvVar(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", errors.New(fmt.Sprintf("Missing env var: %s", key))
	}

	return value, nil
}
