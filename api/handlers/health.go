package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal("Failed to serialize health response:", err)
	}
}
