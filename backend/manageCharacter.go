package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
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
	fmt.Println("getting characters info")
	characters, err := getAllCharacters()
	if err != nil {
		http.Error(w, "couldnt get all characters stored on data dir", http.StatusInternalServerError)
	}
	
	fmt.Println("parsing the characters.html file")
	tmpl, err := template.ParseFiles("../frontend/characters.html")
    if err != nil {
        http.Error(
            w,
            "couldnt parse characters template",
            http.StatusInternalServerError,
        )
        return
    }

	fmt.Println("trying to fill the template with Character struct data")
	err = tmpl.Execute(w, characters)
	if err != nil {
		http.Error(w, "couldnt fill the template with the go struct values", http.StatusInternalServerError)
		return 
	}

	fmt.Println("all done")
}

func createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	err := createCharacter(w, r)
	if err != nil {
		http.Error(w, "Failed to create a new character", http.StatusInternalServerError)
	}
}

func viewCharacterHandler(w http.ResponseWriter, r *http.Request) {
	characterID := r.PathValue("id")
	if characterID == "" {
		http.Error(w, "ID not found", http.StatusInternalServerError)
	}

	c, err := getCharacterByID(characterID)
	if err != nil {
		http.Error(w, "failed to find the character file by ID", http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFiles("../frontend/character.html")
	if err != nil {
		http.Error(w, "failed parising character.html file", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, c)
	if err != nil {
		http.Error(w, "failed to fill character info into template", http.StatusInternalServerError)
		return
	}
}

func editCharacterHandler(w http.ResponseWriter, r *http.Request) {
	characterID := r.PathValue("id")
	if characterID == "" {
		http.Error(w, "ID not found", http.StatusInternalServerError)
		return
	}

	c, err := getCharacterByID(characterID)
	if err != nil {
		http.Error(w, "couldnt get character by ID", http.StatusInternalServerError)
		return
	}
	
	tmpl, err := template.ParseFiles("../frontend/edit.html")
	if err != nil {
		http.Error(w, "wasnt able to parse html file", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, c)
	if err != nil {
		http.Error(w, "wasnt able to render html file", http.StatusInternalServerError)
		return
	}
}

func editCharacterDoneHandler(w http.ResponseWriter, r *http.Request) {
	characterID := r.PathValue("id")

	err := editCharacter(characterID, r)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	
	currentURL := filepath.Join("/character", characterID)
	currentURL = filepath.Join(currentURL, "view")

	http.Redirect(w, r, currentURL, http.StatusFound)
}

func getRandomID() (string, error) {
	dirFiles, err := os.ReadDir("../data")
	if err != nil {
		return "", err
	}

	for {
		bytes := make([]byte, 16)

		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}

		randomID := hex.EncodeToString(bytes)

		exists := false

		for _, file := range dirFiles {
			if randomID + ".json" == file.Name() {
				exists = true
				break
			}
		}

		if !exists {
			return randomID, nil
		}
	}
}

func validateNewCharacter(s []string) (Character, error) {
	if s[1] == "" {
		return Character{}, fmt.Errorf("character with empty name is invalid")
	}

	if s[3] == "" {
		return Character{}, fmt.Errorf("character with empty class is invalid")
	}

	if s[4] == "" {
		return Character{}, fmt.Errorf("character with empty race is invalid")
	}

	level, err := strconv.Atoi(s[2])
	if err != nil {
		return Character{}, err
	}

	c := Character{
		s[0],
		s[1],
		level,
		s[3],
		s[4],
	}

	return c, nil
}

func validateExistingCharacter(s []string) (Character, error){
	oldChar, err := getCharacterByID(s[0])
	if err != nil {
		return Character{}, err
	}

	if s[1] == "" {
		s[1] = oldChar.Name
	}
	if s[3] == "" {
		s[3] = oldChar.Class
	}
	if s[4] == "" {
		s[4] = oldChar.Race
	}

	level, err := strconv.Atoi(s[2])
	if err != nil {
		return Character{}, err
	}

	newChar := Character{
			s[0],
			s[1],
			level,
			s[3],
			s[4],
		}

	return newChar, nil
}

func getCharacterInfo(r *http.Request) (Character, error) {
	err := r.ParseForm()
	if err != nil {
		return Character{}, err
	}	
	
	characterLevel := r.FormValue("level")

	characterName := r.FormValue("name")
	
	characterClass := r.FormValue("class")

	characterRace := r.FormValue("race")

	characterID := r.PathValue("id")

	values := []string {
		characterID, 
		characterName,
		characterLevel,
		characterClass,
		characterRace,
	}

	if characterID == "" {
		characterID, err = getRandomID()
		if err != nil {
			return Character{}, err
		}

		values[0] = characterID

		validated, err := validateNewCharacter(values)
		if err != nil {
			return Character{}, err
		}

		return validated, nil
	}

	validated, err := validateExistingCharacter(values)

	return validated, nil
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

func getAllCharacters() ([]Character, error) {
	files, err := os.ReadDir("../data")
	if err != nil {
		return []Character{}, err
	}

	var filesNames []string

	for _, file := range files {
		filePath := filepath.Join("../data", file.Name())
		filesNames = append(filesNames, filePath)
	}
	
	var filesData [][]byte

	for _, name := range filesNames {
		data, err := os.ReadFile(name)
		if err != nil {
			return []Character{}, err
		}

		filesData = append(filesData, data)
	}

	characters := make([]Character, len(filesData))

	for i, data := range filesData {
		err := json.Unmarshal(data, &characters[i])
		if err != nil {
			return []Character{}, err
		}
	}
	
	return characters, nil
}

func getCharacterByID(characterID string) (Character, error) {
	files, err := os.ReadDir("../data")
	if err != nil {
		return Character{}, err
	}
	
	target := ""
	exists := false

	for _, file := range files {
		if file.Name() == characterID + ".json" {
			exists = true
			filePath := filepath.Join("../data", file.Name())
			target = filePath
			break
		}
	}

	if !exists {
		return Character{}, fmt.Errorf("character not found")
	}

	data, err := os.ReadFile(target)
	if err != nil {
		return Character{}, err
	}
	
	var c Character
	err = json.Unmarshal(data, &c)
	if err != nil {
		return Character{}, err
	}

	return c, nil
}

func editCharacter(characterID string, r *http.Request) error {
	c, err := getCharacterInfo(r)
	if err != nil {
		return err
	}

	data, err := createCharacterJSON(c)
	if err != nil {
		return fmt.Errorf("error while decoding .json")
	}

	err = createJSON(data, characterID)
	if err != nil {
		return err
	}

	return nil
}

//COMMENTS:
/*
*/
