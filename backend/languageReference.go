package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"net/http"
)

func languageReferenceHandler(w http.ResponseWriter, r *http.Request) {
	languages, err := getLanguageReferenceData()	
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(languages)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
}

func getLanguageReferenceData() (map[string][]string, error) {
	root := findRootDir()

	if root == "" {
		return map[string][]string{}, fmt.Errorf("couldnt find root")
	}

	referenceLangsBytes, err := os.ReadFile(filepath.Join(root, "data", "reference", "languages.json"))
	if err != nil {
		return map[string][]string{}, err
	}

	var allLangsMap map[string][]string
	err = json.Unmarshal(referenceLangsBytes, &allLangsMap)
	if err != nil {
		return map[string][]string{}, err
	}

	return allLangsMap, nil
}

