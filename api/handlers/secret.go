package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ruegerj/devops/api/telemetry"
)

func Secret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]any{
		"number":  42,
		"message": "Life, Universe and everything",
	}

	telemetry.UnlockCounter.Inc()
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal("Failed to serialize response:", err)
	}
}
