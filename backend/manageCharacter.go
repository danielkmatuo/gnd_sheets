package main

import (
	"net/http"
)

func charactersHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/characters.html")
}

func createCharacter(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/characters", 302)
}

func getFile(charcterId string) error {
	return nil	
}
