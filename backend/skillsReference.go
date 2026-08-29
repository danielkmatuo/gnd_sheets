package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type AllSkills struct {
	StrSkills []string `json:"str"`
	DexSkills []string `json:"dex"`
	ConSkills []string `json:"con"`
	IntSkills []string `json:"int"`
	WisSkills []string `json:"wis"`
	ChaSkills []string `json:"cha"`
}

func skillsReferenceHandler(w http.ResponseWriter, r *http.Request) {
	skills, err := getSkillsReferenceData()
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(skills)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
}

func getSkillsReferenceData() (AllSkills, error) {
	root := findRootDir()
	if root == "" {
		return AllSkills{}, fmt.Errorf("couldnt find root")
	}

	skillsReferencePath := filepath.Join(root, "data", "reference", "skills.json")

	skillsReferenceBytes, err := os.ReadFile(skillsReferencePath)
	if err != nil {
		return AllSkills{}, err
	}

	var skills AllSkills
	err = json.Unmarshal(skillsReferenceBytes, &skills)
	if err != nil {
		return AllSkills{}, err
	}

	return skills, nil
}
