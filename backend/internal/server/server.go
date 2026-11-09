// Package server wires up the HTTP routes for the Latin Dance Corner API.
package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ApiConfig struct {
	DB *sql.DB
}

// New returns the API handler with every route registered.
func New(apiCfg *ApiConfig) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/hello", handleHello)
	mux.HandleFunc("GET /api/v1/move/{id}", apiCfg.handleGetId)
	mux.HandleFunc("GET /api/v1/moves", apiCfg.handleGetAll)
	return mux
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Hello, world!"})
}

func (apiCfg *ApiConfig) handleGetId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id, must be int", http.StatusNotFound)
		return
	}

	move, err := GetMoveById(apiCfg.DB, int64(id))
	if err != nil {
		http.Error(w, "could not get move", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"move" : move})
}

func (apiCfg *ApiConfig) handleGetAll(w http.ResponseWriter, r *http.Request){
	limitParam := r.URL.Query().Get("limit")
	//prevent shadowing by declaring vars
	var limit int
	var err error
	
	if limitParam == "" {
		limit = 10
	}else {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			http.Error(w, "limit could not be parsed", http.StatusBadRequest)
			return
		}
	}

	offsetParam := r.URL.Query().Get("offset")
	var offset int
	if offsetParam == "" {
		offset = 0
	}else {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			http.Error(w, "limit could not be parsed", http.StatusBadRequest)
			return
		}
	}
	moves, err := GetMultiPaginated(apiCfg.DB, limit, offset)
	if err != nil {
		http.Error(w, "could not get resource", http.StatusInternalServerError)
		log.Printf("error in server: %s", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, moves)
}


func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writing response: %v", err)
	}
}
