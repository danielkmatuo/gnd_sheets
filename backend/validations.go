package backend

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func validateByClassReferenceData(c Character) (Character, error) {
	classInfoMap, err := getReferenceData()	
	if err != nil {
		return Character{}, err
	}

	referenceCharacter := classInfoMap[c.Class]
	
	if len(c.SkillChoice.Chosen) != referenceCharacter.SkillChoice.Cap {
		return Character{}, fmt.Errorf("character must have exact %d skills", referenceCharacter.SkillChoice.Cap)
	} 

	if c.Level == 1 {
		c.MaxHp, err = parseHitDie(referenceCharacter.HitDie)
		c.MaxHp = c.MaxHp + calculateStatBonus(c.Stats.Con)
		if err != nil {
			return Character{}, err
		}
	} else {
		hitDieValue, err := parseHitDie(referenceCharacter.HitDie)
		if err != nil {
			return Character{}, err
		}

		conBonus := calculateStatBonus(c.Stats.Con)
		hitDieStep := int(math.Ceil(float64(hitDieValue / 2))) + 1 //unnecessary math, but worth keeping explict for now
		c.MaxHp = hitDieValue + hitDieStep * (c.Level - 1) + conBonus * (c.Level)
	}

	//left out c.CurrHp on purpose, must be filled during character creation or editing, since I'm reusing this func for both
	c.Profs.Armor = referenceCharacter.Profs.Armor
	c.Profs.Weapon = referenceCharacter.Profs.Weapon
	c.Profs.Tools = referenceCharacter.Profs.Tools
	c.Profs.Saving = referenceCharacter.Profs.Saving
	c.ToolsChoice = referenceCharacter.ToolsChoice
	c.SkillChoice.Cap = referenceCharacter.SkillChoice.Cap
	c.SkillChoice.Possibilities = referenceCharacter.SkillChoice.Possibilities
	c.Spells = referenceCharacter.Spells
	//also left out AC and speed, which also should be filled following a race and equipment based validations

	return c, nil
}

func validateNewCharacter(newCharacter NewCharacter, root string) (Character, error) {
	c := Character{}

	characterID, err := getRandomID(root)
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

	statsValid, err := validateCharacterAbilityScores(newCharacter)
	if !statsValid {
		return Character{}, err
	}

	c.ID = characterID
	c.Name = newCharacter.Name
	c.Class = newCharacter.Class
	c.Race = newCharacter.Race
	c.Level = newCharacter.Level
	c.SkillChoice.Chosen = newCharacter.Skills
	c.Stats = newCharacter.Stats

	validatedByReference, err := validateByClassReferenceData(c)
	if err != nil {
		return Character{}, err
	}
	
	dynamicallyCalculated, err := calculateDynamicallyNewCharacter(validatedByReference)
	if err != nil {
		return Character{}, err
	}

	return dynamicallyCalculated, nil
}

func validateExistingCharacter(c Character) (Character, error){
	var validated Character
	referenceMap, err := getReferenceData()
	if err != nil {
		return Character{}, err
	}

	reference := referenceMap[c.Class]

	if len(c.SkillChoice.Chosen) != reference.SkillChoice.Cap {
		return Character{}, fmt.Errorf("character must have exact %d", reference.SkillChoice.Cap)
	} 

	if c.Level == 1 {
		validated.MaxHp, err = parseHitDie(reference.HitDie)
		if err != nil {
			return Character{}, err
		}
	} else {
		hitDieValue, err := parseHitDie(reference.HitDie)
		if err != nil {
			return Character{}, err
		}

		hitDieStep := int(math.Floor(float64(hitDieValue / 2)))
		validated.MaxHp = hitDieValue + hitDieStep * (c.Level - 1)
	}

	//passing values from c to validated
	validated.AC = c.AC
	validated.Speed = c.Speed
	validated.Level = c.Level
	validated.Class = c.Class
	validated.ID = c.ID
	validated.Name = c.Name
	validated.CurrHp = c.CurrHp
	validated.Race = c.Race

	//left out c.CurrHp on purpose, must be filled during character creation or editing, since I'm reusing this func for both
	validated.Profs.Armor = reference.Profs.Armor
	validated.Profs.Weapon = reference.Profs.Weapon
	validated.Profs.Tools = reference.Profs.Tools
	validated.Profs.Saving = reference.Profs.Saving
	validated.ToolsChoice = reference.ToolsChoice
	validated.SkillChoice.Cap = reference.SkillChoice.Cap
	validated.SkillChoice.Possibilities = reference.SkillChoice.Possibilities
	validated.Spells = reference.Spells
	//also left out AC and speed, which also should be filled following a race and equipment based validations

	return validated, nil
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

func validateAbilityCost(stats Attributes) bool {
	statsValues := [6]int{stats.Str, stats.Dex, stats.Con, stats.Int, stats.Wis, stats.Cha}
	currCost := calculateCurrCost(statsValues)

	if currCost > 27 || currCost < 0 {
		return false
	}

	return true
}

func calculateCurrCost(stats [6]int) int {
	costMap := make(map[int]int)
	allowedStatsValues := [8]int{8, 9, 10, 11, 12, 13, 14, 15}
	allowedCostsValues := [8]int{0, 1, 2, 3, 4, 5, 7, 9}

	for i, stat := range allowedStatsValues{
		costMap[stat] = allowedCostsValues[i]	
	}

	currCost := 0

	for _, value := range stats {
		if value >= 8 && value <= 15 {
			currCost += costMap[value]	
		} else {
			return -1
		}
	}

	return currCost
}

func calculateStatBonus(stat int) int {
	const baselineStatValue int = 10
	var bonus int

	if stat % 2 == 0 {
		if stat >= baselineStatValue {
			//if str = 12, then stat - baseline = 12 - 10  = 2
			//then 2/2 = +1 = str bonus
			bonus = (stat - baselineStatValue) / 2
		} else {
			//if str = 8, then baseline - stat = 10 - 8 = 2
			//then 2*-1 = -2 -> -2/2 = -1 = str bonus
			bonus = ((baselineStatValue - stat) * (-1)) / 2
		}

		return bonus
	}

	if stat < baselineStatValue {
		//add +1 to normalize the score -> if str = 7, then baseline - stat + 1 = 10 - 7 + 1 = 10 - 6 = 4
		//then 4 * -1 = -4 -> -4/2 = -2 = str bonus
		bonus = ((baselineStatValue - stat + 1) * (-1)) / 2	
		return bonus
	}

	//add -1 to normalize the score -> if str = 13, then stat - baseline - 1 = 13 - 10 - 1 = 13 - 11 = 2
	//then 2/2 = +1 = str bonus
	bonus = (stat - baselineStatValue - 1) / 2

	return bonus
}

func createAnchor(stats Attributes) ([6]int, [6]int) {
	finalState := [6]int{stats.Str, stats.Dex, stats.Con, stats.Int, stats.Wis, stats.Cha}

	var absDiffVec [6]int
	var anchor [6]int

	for i, stat := range finalState {
		absDiffVec[i] = stat - 8	
		if absDiffVec[i] < 5 {
			anchor[i] = stat
		} else {
			anchor[i] = 13
		}
	}

	currCost := calculateCurrCost(anchor)
	if currCost > 27 {
		anchor[0] = 12
		anchor[1] = 12
		anchor[2] = 12
	} else if currCost < 0 {
		return [6]int{-1}, [6]int{-1}
	}

	return anchor, absDiffVec
}

func getStatesDiff(currState [6]int, finalState [6]int) [6]int {
	var currDiffVec [6]int
	for i, statCap := range finalState {
		currDiffVec[i] = currState[i] - statCap
	}

	return currDiffVec
}

func validateRangeCharacterAbilityScores(level int, c NewCharacter) bool {
	if level > 1 {
		if c.Stats.Str < 8 || c.Stats.Str > 20 {
			return false
		} else if c.Stats.Dex < 8 || c.Stats.Dex > 20 {
			return false
		} else if c.Stats.Con < 8 || c.Stats.Con > 20 {
			return false
		} else if c.Stats.Int < 8 || c.Stats.Int > 20 {
			return false
		} else if c.Stats.Wis < 8 || c.Stats.Wis > 20 {
			return false
		} else if c.Stats.Cha < 8 || c.Stats.Cha > 20 {
			return false
		} 
	}

	costOk := validateAbilityCost(c.Stats)
	if !costOk {
		return false
	} else if c.Stats.Str < 8 || c.Stats.Str > 15 {
		return false
	} else if c.Stats.Dex < 8 || c.Stats.Dex > 15 {
		return false
	} else if c.Stats.Con < 8 || c.Stats.Con > 15 {
		return false
	} else if c.Stats.Int < 8 || c.Stats.Int > 15 {
		return false
	} else if c.Stats.Wis < 8 || c.Stats.Wis > 15 {
		return false
	} else if c.Stats.Cha < 8 || c.Stats.Cha > 15 {
		return false
	}
	
	return true
}

func validateCharacterAbilityScores(c NewCharacter) (bool, error) {
	validateRangeAndCost := validateRangeCharacterAbilityScores(c.Level, c)
	if !validateRangeAndCost {
		return false, fmt.Errorf("Character doesnt follow character creation level rules or doesnt follow cost buy rules")
	}

	var currState [6]int
	anchor, absDiffVec := createAnchor(c.Stats)
	if anchor[0] < 0 {
		return false, fmt.Errorf("Invalid value for stat in the cost buy rules")
	}
	bonusPoints := defineBonusPoints(c.Level)
	statsCap := [6]int{
		c.Stats.Str, 
		c.Stats.Dex, 
		c.Stats.Con, 
		c.Stats.Int, 
		c.Stats.Wis, 
		c.Stats.Cha,
	}

	for {
		currCost := calculateCurrCost(anchor)
		currState = anchor
		currDiffVec := getStatesDiff(currState, statsCap)
		currDiffSum := 0

		for _, diff := range currDiffVec {
			currDiffSum += diff	
		} 

		if currCost > 27 || currCost < 0 {
			return false, fmt.Errorf("Character has invalid attributes values")
		} else if currCost <= 27 && currDiffSum == bonusPoints {
			break
		}

		minAbsDiff := 100
		pointer := 0
		for i, diff := range currDiffVec {
			if diff <= minAbsDiff && diff > 0 {
				minAbsDiff = diff
				pointer = i
			}
		}

		anchor[pointer]++
		absDiffVec[pointer]--
	}

	currState = anchor
	bonusPointsDiff := getStatesDiff(currState, statsCap) 

	const maxIterations int = 100

	for i := 0; i < maxIterations; i++ {
		if currState == statsCap {
			return true, nil
		}	

		for i, diff := range bonusPointsDiff {
			if diff > 0 && bonusPoints > 0{
				bonusPointsDiff[i]--
				bonusPoints--
				currState[i]++
			} else if diff > 0 && bonusPoints <= 0 {
				return false, fmt.Errorf("Character has invalid stats, given that absDiffVec still has non zero values but ran out of bonus points")
			}
		}
	}

	return false, fmt.Errorf("Character has invalid stats, even after achieving max iterations cap")
}

func defineBonusPoints(level int) int {
	extraCost := 3
	for i := range level {
		if (i + 1) % 4 == 0 {
			extraCost += 2
		}
	}

	return extraCost
}

func calculateDynamicallyNewCharacter(c Character) (Character, error) {
	referenceMap, err := getReferenceData()
	if err != nil {
		return Character{}, err
	}

	reference := referenceMap[c.Class]
	referenceHitDieValue, err := parseHitDie(reference.HitDie)
	if err != nil {
		return Character{}, err
	}

	conBonus := calculateStatBonus(c.Stats.Con)
	referenceHitDieStep := int(math.Ceil(float64(referenceHitDieValue / 2))) + 1 //clearly unnecessary math, but still worth keeping for now 
	var referenceMaxHp int

	if c.Level == 1 {
		referenceMaxHp = referenceHitDieValue + conBonus
		if referenceMaxHp != c.MaxHp {
			return Character{}, fmt.Errorf("character has invalid max hp value")
		}
	}
	
	areEqual := false

	referenceMaxHp = referenceHitDieValue + referenceHitDieStep * (c.Level - 1) + conBonus * (c.Level)
	if referenceMaxHp == c.MaxHp {
		areEqual = true
	} 
	
	if !areEqual {
		return Character{}, fmt.Errorf("character has invalid max hp value")
	}

	c.CurrHp = referenceMaxHp
	c.AC = 10 + calculateStatBonus(c.Stats.Dex)
	c.Speed = 9.0

	return c, nil
}
