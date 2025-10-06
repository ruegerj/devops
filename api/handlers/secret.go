package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func Secret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]any{
		"number":  42,
		"message": "Life, Universe and everything",
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal("Failed to serialize response:", err)
	}
}
