package main

import (
	"net/http"
	"fmt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/index.html")
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/js/app.js")
}

func configureServer() {
	http.HandleFunc("GET /{$}", indexHandler)
	http.HandleFunc("GET /js/reference", jsHandler)

	http.HandleFunc("GET /reference/classes/{class}", classReferenceHandler)
	http.HandleFunc("GET /characters", charactersHandler)
	http.HandleFunc("POST /characters/create", createCharacterHandler)
	http.HandleFunc("GET /character/{id}/view", viewCharacterHandler)
	http.HandleFunc("GET /character/{id}/edit", editCharacterHandler)
	http.HandleFunc("POST /character/{id}/edit/done", editCharacterDoneHandler)
}

func runServer(port string) error {
	fmt.Println("Starting server at http://localhost:8080")

	return http.ListenAndServe(port, nil)
}
