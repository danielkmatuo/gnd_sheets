package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	WhatTools []ToolChoice `json:"tool_proficiency_choices"`
	SkillChoice Skills `json:"skill_choices"`
	Spells Spellcasting `json:"spellcasting"`
} 

func classReferenceHandler(w http.ResponseWriter, r *http.Request) {
	className := r.PathValue("class")

	data, err := os.ReadFile("../data/reference/classes.json")
	if err != nil {
		http.Error(w, "could not read class reference data", http.StatusInternalServerError)
		return
	}

	var classes map[string]ClassInfo

	err = json.Unmarshal(data, &classes)
	if err != nil {
		http.Error(w, "could not parse class reference data", http.StatusInternalServerError)
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
	classesData, err := os.ReadFile("../data/reference/classes.json")
	if err != nil {
		return map[string]ClassInfo{}, err
	}

	var ci map[string]ClassInfo
	err = json.Unmarshal(classesData, &ci)
	if err != nil {
		return map[string]ClassInfo{}, err
	}

	fmt.Printf("decoded reference json data: \n%v", ci) 

	return ci, nil
}

func getClassInfoFromMap(ciMap map[string]ClassInfo, class string) ClassInfo {
	return ciMap[class]
}
