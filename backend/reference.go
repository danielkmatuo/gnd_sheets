package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type ToolChoice struct {
	Quantity int `json:"choose"`
	Categories []string `json:"from_categories"`
}

type Proficiency struct {
	Armor []string `json:"armor"`
	Weapon []string `json:"weapons"`
	Tools []string `json:"tools"`
	Saving []string `json:"saving_throws"`
}

type Skills struct {
	Cap int `json:"choose"`
	Possibilities []string `json:"from"`
	Chosen []string `json:"skills"`
}

type Spellcasting struct {
	Ability string `json:"ability"`
	StartsAt int `json:"starts_at_level"`
}

type ClassInfo struct {
	Name string `json:"name"`
	HitDie string `json:"hit_die"`
	Profs Proficiency `json:"proficiencies"`
	ToolsChoice []ToolChoice `json:"tool_proficiency_choices"`
	SkillChoice Skills `json:"skill_choices"`
	Spells Spellcasting `json:"spellcasting"`
} 

func classReferenceHandler(w http.ResponseWriter, r *http.Request) {
	className := r.PathValue("class")

	classes, err := getReferenceData()
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	class, ok := classes[className]
	if !ok {
		http.Error(w, "class not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(class)
	if err != nil {
		http.Error(w, "could not send class data", http.StatusInternalServerError)
		return
	}
}

func getReferenceData() (map[string]ClassInfo, error) {
	root := findRootDir() 
	if root == "" {
		return map[string]ClassInfo{}, fmt.Errorf("couldnt find root")
	}

	classesData, err := os.ReadFile(filepath.Join(root, "data", "reference", "classes.json"))
	if err != nil {
		return map[string]ClassInfo{}, err
	}

	var ci map[string]ClassInfo
	err = json.Unmarshal(classesData, &ci)
	if err != nil {
		return map[string]ClassInfo{}, err
	}

	return ci, nil
}
