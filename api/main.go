package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ruegerj/devops/api/handlers"
	"github.com/ruegerj/devops/api/middleware"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	host, err := resolveEnvVar("HOST")
	if err != nil {
		log.Panic(err)
	}

	port, err := resolveEnvVar("PORT")
	if err != nil {
		log.Panic(err)
	}

	jwtKey, err := resolveEnvVar("JWT_KEY")
	if err != nil {
		log.Panic(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/health", handlers.Health)
	router.Handle("/api/secret", middleware.Authenticate(jwtKey, logger, http.HandlerFunc(handlers.Secret)))

	handler := middleware.Logging(logger, router)
	addr := fmt.Sprintf("%s:%s", host, port)

	fmt.Printf("Server starting on %s...\n", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("Server failed to start:", err)
		os.Exit(1)
	}
}
