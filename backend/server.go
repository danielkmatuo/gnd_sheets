package backend

import (
	"net/http"
	"fmt"
	"path/filepath"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	root := findRootDir() 
	if root == "" {
		http.Error(w, "couldnt find root", http.StatusInternalServerError)
	}

	http.ServeFile(w, r, filepath.Join(root, "frontend", "index.html"))
}

func jsCreateHandler(w http.ResponseWriter, r *http.Request) {
	root := findRootDir() 
	if root == "" {
		http.Error(w, "couldnt find root", http.StatusInternalServerError)
	}
	
	http.ServeFile(w, r, filepath.Join(root, "frontend", "js", "create.js"))
}

func jsEditHandler(w http.ResponseWriter, r *http.Request) {
	root := findRootDir() 
	if root == "" {
		http.Error(w, "couldnt find root", http.StatusInternalServerError)
	}
	
	http.ServeFile(w, r, filepath.Join(root, "frontend", "js", "edit.js"))
}

func jsViewHandler(w http.ResponseWriter, r *http.Request) {
	root := findRootDir() 
	if root == "" {
		http.Error(w, "couldnt find root", http.StatusInternalServerError)
	}
	
	http.ServeFile(w, r, filepath.Join(root, "frontend", "js", "view.js"))
}

func editPageHandler(w http.ResponseWriter, r *http.Request) {
	root := findRootDir() 
	if root == "" {
		http.Error(w, "couldnt find root", http.StatusInternalServerError)
	}
	
	http.ServeFile(w, r, filepath.Join(root, "frontend", "edit.html"))
}

func viewPageHandler(w http.ResponseWriter, r *http.Request) {
	root := findRootDir() 
	if root == "" {
		http.Error(w, "couldnt find root", http.StatusInternalServerError)
	}

	http.ServeFile(w, r, filepath.Join(root, "frontend", "character.html"))
}

func configureServer() {
	//handlers to serve static files
	http.HandleFunc("GET /{$}", indexHandler)
	http.HandleFunc("GET /js/create", jsCreateHandler)
	http.HandleFunc("GET /js/edit", jsEditHandler)
	http.HandleFunc("GET /js/view", jsViewHandler)

	//handlers to get reference data
	http.HandleFunc("GET /reference/classes/{class}", classReferenceHandler)
	http.HandleFunc("GET /reference/languages", languageReferenceHandler)
	http.HandleFunc("GET /reference/skills", skillsReferenceHandler)
	http.HandleFunc("GET /reference/races/{race}", raceReferenceHandler)

	//handlers to create a new character
	http.HandleFunc("POST /characters/create", createCharacterHandler)

	//handlers to send stored characters
	http.HandleFunc("GET /character/{id}/data", sendCharacterDataHandler)

	//handlers to view characters
	http.HandleFunc("GET /characters", charactersHandler)
	http.HandleFunc("GET /character/{id}/view", viewPageHandler)

	//handlers to edit character
	http.HandleFunc("GET /character/{id}/edit", editPageHandler)
	http.HandleFunc("POST /character/{id}/edit/done", editCharacterDoneHandler)
}

func runServer(port string) error {
	fmt.Println("Starting server at http://localhost:8080")

	return http.ListenAndServe(port, nil)
}
