package backend

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Attributes struct {
	Str int `json:"str"`
	Dex int `json:"dex"`
	Con int `json:"con"`
	Int int `json:"int"`
	Wis int `json:"wis"`
	Cha int `json:"cha"`
}

type Character struct {
	ID string    `json:"id"`
	Name string  `json:"name"`
	Level int    `json:"level"`
	Class string `json:"class"`
	Race string  `json:"race"`
	MaxHp int `json:"max_hp"`
	CurrHp int `json:"current_hp"`
	Profs Proficiency `json:"proficiencies"`
	ToolsChoice []ToolChoice `json:"tool_proficiency_choices"`
	SkillChoice Skills `json:"skill_choices"`
	Spells Spellcasting `json:"spellcasting"`
	AC int `json:"ac"`
	Speed float64 `json:"speed"`
	Stats Attributes `json:"stats"`
}

type NewCharacter struct {
	Name string `json:"name"`
	Level int `json:"level"`
	Class string `json:"class"`
	Race string `json:"race"`
	Skills []string `json:"skills"`
	Stats Attributes `json:"stats"`
}

func charactersHandler(w http.ResponseWriter, r *http.Request) {
	root := findRootDir()
	if root == "" {
		http.Error(w, "couldnt find root", http.StatusInternalServerError)
	}

	characters, err := getAllCharacters(root)
	if err != nil {
		http.Error(w, "couldnt get all characters stored on data dir", http.StatusInternalServerError)
		return
	}
	
	tmpl, err := template.ParseFiles(filepath.Join(root, "frontend", "characters.html"))
    if err != nil {
        http.Error(
            w,
            "couldnt parse characters template",
            http.StatusInternalServerError,
        )
        return
    }

	err = tmpl.Execute(w, characters)
	if err != nil {
		http.Error(w, "couldnt fill the template with the go struct values", http.StatusInternalServerError)
		return 
	}
}

func createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	err := createCharacter(w, r)
	if err != nil {
		http.Error(w, "Failed to create a new character", http.StatusInternalServerError)
		return
	}
}

func sendCharacterDataHandler(w http.ResponseWriter, r *http.Request) {
	characterID := r.PathValue("id")
	if characterID == "" {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	err := sendCharacterData(w, characterID)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
}

//TODO: make this func receive an edited character from the frontend and then make it validate the edited character
func editCharacterDoneHandler(w http.ResponseWriter, r *http.Request) {
	characterID := r.PathValue("id")

	root := findRootDir()
	if root == "" {
		http.Error(w, "couldnt find root", http.StatusInternalServerError)
	}

	err := editCharacter(characterID, root)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	
	currentURL := path.Join("/character", characterID)
	currentURL = path.Join(currentURL, "view")

	http.Redirect(w, r, currentURL, http.StatusFound)
}

func getRandomID(root string) (string, error) {
	dirFiles, err := os.ReadDir(filepath.Join(root, "data", "characters"))
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

func sendCharacterData(w http.ResponseWriter, characterID string) error {
	root := findRootDir() 
	if root == "" {
		return fmt.Errorf("couldnt find root")
	}

	c, err := getCharacterByID(characterID, root)
	if err != nil {
		return err
	}
	
	w.Header().Set("Content-Type", "application/json")	
	err = json.NewEncoder(w).Encode(c)
	if err != nil {
		return err
	}

	return nil
}

func getCharacterInfoCreation(r *http.Request) (NewCharacter, error) {
	var c NewCharacter

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		return NewCharacter{}, err
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

func createJSON(data []byte, ID string, root string) error {
	charactersPath := filepath.Join(root, "data", "characters")
	fileName := filepath.Join(charactersPath, ID + ".json")

	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}	
	
	return nil
}

func createCharacter(w http.ResponseWriter, r *http.Request) error {
	root := findRootDir()
	if root == "" {
		return fmt.Errorf("couldnt find root")
	}

	newC, err := getCharacterInfoCreation(r)
	if err != nil {
		return err
	}

	c, err := validateNewCharacter(newC, root)
	if err != nil {
		return err
	}

	validatedC, err := validateByClassReferenceData(c)
	if err != nil {
		return err
	}

	validatedC, err = calculateDynamicallyNewCharacter(validatedC)
	if err != nil {
		return err
	}

	data, err := createCharacterJSON(validatedC)
	if err != nil {
		return err
	}
	
	err = createJSON(data, validatedC.ID, root)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

func getAllCharacters(root string) ([]Character, error) {
	files, err := os.ReadDir(filepath.Join(root, "data", "characters"))
	if err != nil {
		return []Character{}, err
	}

	var filesNames []string

	for _, file := range files {
		filePath := filepath.Join(filepath.Join(root, "data", "characters"), file.Name())
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

func getCharacterByID(characterID string, root string) (Character, error) {
	files, err := os.ReadDir(filepath.Join(root, "data", "characters"))
	if err != nil {
		return Character{}, err
	}
	
	target := ""
	exists := false

	for _, file := range files {
		if file.Name() == characterID + ".json" {
			exists = true
			filePath := filepath.Join(filepath.Join(root, "data", "characters"), file.Name())
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

func editCharacter(characterID string, root string) error {
	c, err := getCharacterByID(characterID, root)
	if err != nil {
		return err
	}

	data, err := createCharacterJSON(c)
	if err != nil {
		return fmt.Errorf("error while encoding .json")
	}

	err = createJSON(data, characterID, root)
	if err != nil {
		return err
	}

	return nil
}

//COMMENTS:
/*
*/
