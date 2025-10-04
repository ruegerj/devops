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

	router := http.NewServeMux()
	router.HandleFunc("/health", handlers.Health)

	host, err := resolveEnvVar("HOST")
	if err != nil {
		log.Panic(err)
	}

	port, err := resolveEnvVar("PORT")
	if err != nil {
		log.Panic(err)
	}

	handler := middleware.Logging(logger, router)
	addr := fmt.Sprintf("%s:%s", host, port)

	fmt.Printf("Server starting on %s...\n", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("Server failed to start:", err)
		os.Exit(1)
	}
}
