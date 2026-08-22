package backend

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

func jsViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/js/view.js")
}

func editPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/edit.html")
}

func viewPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/character.html")
}

func configureServer() {
	http.HandleFunc("GET /{$}", indexHandler)
	http.HandleFunc("GET /js/create", jsCreateHandler)
	http.HandleFunc("GET /js/edit", jsEditHandler)
	http.HandleFunc("GET /js/view", jsViewHandler)

	http.HandleFunc("GET /reference/classes/{class}", classReferenceHandler)
	http.HandleFunc("GET /characters", charactersHandler)
	http.HandleFunc("POST /characters/create", createCharacterHandler)
	http.HandleFunc("GET /character/{id}/data", sendCharacterDataHandler)
	http.HandleFunc("GET /character/{id}/edit", editPageHandler)
	http.HandleFunc("GET /character/{id}/view", viewPageHandler)
	http.HandleFunc("POST /character/{id}/edit/done", editCharacterDoneHandler)
}

func runServer(port string) error {
	fmt.Println("Starting server at http://localhost:8080")

	return http.ListenAndServe(port, nil)
}
