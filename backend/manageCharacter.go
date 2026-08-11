package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Character struct {
	ID string    `json:"id"`
	Name string  `json:"name"`
	Level int    `json:"level"`
	Class string `json:"class"`
	Race string  `json:"race"`
}

func charactersHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/characters.html")
}

func createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	err := createCharacter(w, r)
	if err != nil {
		http.Error(w, "Failed to create a new character", http.StatusInternalServerError)
	}
}

func getRandomID() (string, error) {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getCharacterInfo(r *http.Request) (Character, error) {
	err := r.ParseForm()
	if err != nil {
		return Character{}, err
	}	

	characeterID, err := getRandomID()
	if err != nil {
		return Character{}, err
	}
	
	characterLevel, err := strconv.Atoi(r.FormValue("level"))
	if err != nil {
		return Character{}, err
	}

	c := Character{
		ID: characeterID,
		Name: r.FormValue("name"),
		Level: characterLevel,
		Class: r.FormValue("class"),
		Race: r.FormValue("race"),
	}

	return c, nil
}

func createCharacterJSON(c Character) ([]byte, error) {
	data, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func createJSON(data []byte, ID string) error {
	fileName := filepath.Join("../data", ID + ".json")

	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}	
	
	return nil
}

func createCharacter(w http.ResponseWriter, r *http.Request) error {
	c, err := getCharacterInfo(r)
	if err != nil {
		return err
	}

	data, err := createCharacterJSON(c)
	if err != nil {
		return err
	}
	
	err = createJSON(data, c.ID)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/characters", http.StatusFound)

	return nil
}

func getFile(charcterID string) error {
	return nil	
}
