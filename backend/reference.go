package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ToolChoice struct {
	Quantity int `json:choose`
	Categories []string `json:from_categories`
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../frontend/index.html")
	if err != nil {
		http.Error(w, "failed parising character.html file", http.StatusInternalServerError)
		return
	}
	
	class := r.FormValue("class")

	ciMap, err := getReferenceData()
	if err != nil {
		http.Error(w, "failed fetching data from the reference classes.json file", http.StatusInternalServerError)
		return
	}

	ci := getClassInfoFromMap(ciMap, class)

	err = tmpl.Execute(w, ci)
	if err != nil {
		http.Error(w, "failed filling the template with class info", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/characters", http.StatusFound)
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

func parseHitDie(value string) (int, error) {
	parsedString, foundPrefix := strings.CutPrefix(value, "d")
	if !foundPrefix {
		return 0, fmt.Errorf("wrong hit die format, must start with d")
	}

	parsedValue, err := strconv.Atoi(parsedString)
	if err != nil {
		return 0, err
	}

	return parsedValue, nil
}

func getClassInfoFromMap(ciMap map[string]ClassInfo, class string) ClassInfo {
	return ciMap[class]
}

func validateReferenceData(class string, referenceMap map[string]ClassInfo) (Character, error) {
	

	return Character{}, nil
}
