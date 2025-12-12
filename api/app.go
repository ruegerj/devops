package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ruegerj/devops/api/handlers"
	"github.com/ruegerj/devops/api/middleware"
	"github.com/ruegerj/devops/api/telemetry"
)

type App struct {
	router        *http.ServeMux
	logger        *slog.Logger
	jwtSigningKey string
}

func (a *App) Initialize(jwtKey string, isTelemetryEnabled bool) {
	a.jwtSigningKey = jwtKey
	a.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	a.router = http.NewServeMux()

	a.initializeRoutes()

	if isTelemetryEnabled {
		a.setupTelemtryCollection()
	}
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

func (a *App) setupTelemtryCollection() {
	registry := prometheus.NewRegistry()

	// custom metrics
	registry.MustRegister(
		telemetry.SuccessfulAuthCounter,
		telemetry.FailedAuthCounter,
		telemetry.UnlockCounter,
	)

	// general runtime metrics
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	a.router.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
}
