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

type App struct {
	router        *http.ServeMux
	logger        *slog.Logger
	jwtSigningKey string
}

func (a *App) Initialize(jwtKey string) {
	a.jwtSigningKey = jwtKey
	a.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	a.router = http.NewServeMux()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	// setup global middleware stack
	handler := middleware.Logging(a.logger, a.router)

	fmt.Printf("Server starting on %s...\n", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("Server failed to start:", err)
		os.Exit(1)
	}
}

func (a *App) initializeRoutes() {
	a.router.HandleFunc("GET /health", handlers.Health)
	a.router.Handle("GET /api/secret", middleware.Authenticate(
		a.jwtSigningKey,
		a.logger,
		http.HandlerFunc(handlers.Secret)))
}
