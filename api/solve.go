package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Minimal handler /api/solve reached!")

	origin := r.Header.Get("Origin")
	allowedOrigins := []string{
		"https://tubes2-fe-ian.vercel.app",
		"http://localhost:3000",
	}
	for _, o := range allowedOrigins {
		if origin == o {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "API /api/solve is alive!"}
	json.NewEncoder(w).Encode(response)
}