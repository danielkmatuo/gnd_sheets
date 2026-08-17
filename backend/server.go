package main

import (
	"net/http"
	"fmt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/index.html")
}

func jsCreateHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/js/create.js")
}

func jsEditHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/js/edit.js")
}

func editPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/edit.html")
}

func configureServer() {
	http.HandleFunc("GET /{$}", indexHandler)
	http.HandleFunc("GET /js/create", jsCreateHandler)
	http.HandleFunc("GET /js/edit", jsEditHandler)

	http.HandleFunc("GET /reference/classes/{class}", classReferenceHandler)
	http.HandleFunc("GET /characters", charactersHandler)
	http.HandleFunc("POST /characters/create", createCharacterHandler)
	http.HandleFunc("GET /character/{id}/view", viewCharacterHandler)
	http.HandleFunc("GET /character/edit", editPageHandler)
	http.HandleFunc("GET /character/{id}/data", getCharacterDataHandler)
	http.HandleFunc("POST /character/{id}/edit/done", editCharacterDoneHandler)
}

func runServer(port string) error {
	fmt.Println("Starting server at http://localhost:8080")

	return http.ListenAndServe(port, nil)
}
