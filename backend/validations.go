package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func checkIDExist (characterID string) (bool, error) {
	files, err := os.ReadDir("../data/characters")
	exist := false

	for _, file := range(files){
		if characterID + ".json" == file.Name() {
			characterID, err = getRandomID()
			if err != nil {
				return false, err
			}

			exist = true
			break
		}
	}

	return exist, nil
}

//TODO: refactor validateExistingCharacter to fit the new js + go workflow and to ensure data is properly handled before storage
func validateNewCharacter(newCharacter NewCharacter) (Character, error) {
	c := Character{}

	characterID, err := getRandomID()
	if err != nil {
		return Character{}, err
	}	

	if newCharacter.Name == "" {
		return Character{}, fmt.Errorf("name cannot be empty")
	}

	if newCharacter.Class == "" {
		return Character{}, fmt.Errorf("class cannot be empty")
	}

	if newCharacter.Race == "" {
		return Character{}, fmt.Errorf("race cannot be empty")
	}

	if newCharacter.Level < 1 || newCharacter.Level > 20 {
		return Character{}, fmt.Errorf("invalid level: must be within 1 to 20 range")
	}

	if newCharacter.Skills == nil {
		return Character{}, fmt.Errorf("cannot choose 0 skills")
	}

	c.ID = characterID
	c.Name = newCharacter.Name
	c.Class = newCharacter.Class
	c.Race = newCharacter.Race
	c.Level = newCharacter.Level
	c.SkillChoice.Chosen = newCharacter.Skills

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

//	level, err := strconv.Atoi(s[2])
//	if err != nil {
//		return Character{}, err
//	}


	//TODO: return correct Character{} struct
	return Character{}, nil
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

//TODO: make this func the "server validate character using business rules" step for the three step validation (1. js; 2. simple go validation; 3. validation by reference)
func validateByClassReferenceData(c Character) (Character, error) {
	

	return Character{}, nil
}
