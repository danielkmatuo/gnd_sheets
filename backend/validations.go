package main

import (
	"strconv"
	"fmt"
)

//TODO: refactor validateNewCharacter and validateExistingCharacter to fit the new js + go workflow and to ensure data is properly handled before storage
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

//	level, err := strconv.Atoi(s[2])
//	if err != nil {
//		return Character{}, err
//	}

	//TODO: Correct the Character struct for data validation
	c := Character{}

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

func validateReferenceData(class string, referenceMap map[string]ClassInfo) (Character, error) {
	

	return Character{}, nil
}
