package main

import (
	"log"
	"net/http"
	"os"

	"github.com/JordiWeber123/latin-dance-corner/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("API listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, server.New()); err != nil {
		log.Fatal(err)
	}
}
