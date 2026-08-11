package main

import (
	"net/http"
	"fmt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/index.html")
}

func configureServer() {
	http.HandleFunc("GET /{$}", indexHandler)
	http.HandleFunc("GET /characters",charactersHandler)
	http.HandleFunc("POST /characters/create", createCharacter)
}

func runServer(port string) error {
	fmt.Println("Starting server at http://localhost:8080")

	err := http.ListenAndServe(port, nil)
	return err
}
