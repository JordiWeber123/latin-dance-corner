package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/JordiWeber123/latin-dance-corner/internal/server"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	db_link := os.Getenv("DB_CONN_STR")
	if db_link == "" {
		log.Fatal("no database, empty connection link ")
	}

	db, err := sql.Open("pgx", db_link)
	if err != nil {
		log.Fatalf("couldn't connect to database: %v\n", err)
	}

	apiCfg := server.ApiConfig{DB: db}

	log.Printf("API listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, server.New(&apiCfg)); err != nil {
		log.Fatal(err)
	}
}
