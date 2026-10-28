// Package server wires up the HTTP routes for the Latin Dance Corner API.
package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// New returns the API handler with every route registered.
func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/hello", handleHello)
	return mux
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Hello, world!"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writing response: %v", err)
	}
}
