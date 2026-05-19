package main

import (
	"log"
	"net/http"

	"contractgen/internal/web"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", web.IndexHandler)
	mux.HandleFunc("POST /generate", web.GenerateHandler)

	log.Println("listening on port: 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}