package backend

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

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

	stats, err := mergePointBuyAndBonus(newCharacter)
	if err != nil {
		return Character{}, err
	}

	c.ID = characterID
	c.Name = newCharacter.Name
	c.Class = newCharacter.Class
	c.Race = newCharacter.Race
	c.Level = newCharacter.Level
	c.SkillChoice.Chosen = newCharacter.Skills
	c.Stats = stats

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

func validateByClassReferenceData(c Character) (Character, error) {
	classInfoMap, err := getReferenceData()	
	if err != nil {
		return Character{}, err
	}

	referenceCharacter := classInfoMap[c.Class]
	
	skillsOk, err := validateSkillsByReference(referenceCharacter, c)
	if err != nil {
		return Character{}, err
	} else if !skillsOk {
		return Character{}, fmt.Errorf("Something unexpected happened in validateSkillsByReference")
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

func validateExistingCharacter(c Character) (Character, error){
	var validated Character
	referenceMap, err := getReferenceData()
	if err != nil {
		return Character{}, err
	}

	reference := referenceMap[c.Class]

	skillsOk, err := validateSkillsByReference(reference, c)
	if err != nil {
		return Character{}, err
	} else if !skillsOk {
		return Character{}, fmt.Errorf("Something unexpected happened in validateSkillsByReference")
	}

	if c.Level == 1 {
		validated.MaxHp, err = parseHitDie(reference.HitDie)
		validated.MaxHp += calculateStatBonus(c.Stats.Con)
		if err != nil {
			return Character{}, err
		}
	} else {
		hitDieValue, err := parseHitDie(reference.HitDie)
		if err != nil {
			return Character{}, err
		}

		hitDieStep := int(math.Floor(float64(hitDieValue / 2))) + 1
		validated.MaxHp = hitDieValue + hitDieStep * (c.Level - 1) + calculateStatBonus(c.Stats.Con)
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

func validateAbilityCost(stats [6]int) (bool, error) {
	currCost, err := calculateCurrCost(stats)
	if err != nil {
		return false, err
	}

	if currCost > 27 {
		return false, fmt.Errorf("Invalid character ability score. Current score %d is above 27", currCost)
	}

	return true, nil
}

func calculateCurrCost(stats [6]int) (int, error) {
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
			return -1, fmt.Errorf("Invalid character ability score. Score %d is out of the range 8 to 15", value)
		}
	}

	return currCost, nil
}

func validatePointBuy(c NewCharacter) (bool, error) {
	pointsBought := [6]int{
		c.PointBuy.Str,
		c.PointBuy.Dex,
		c.PointBuy.Con,
		c.PointBuy.Int,
		c.PointBuy.Wis,
		c.PointBuy.Cha,
	}

	costOk, err := validateAbilityCost(pointsBought)
	if err != nil {
		return false, err
	}
	
	if !costOk {
		return false, fmt.Errorf("Something unexpected happened in validatePointBuy()")
	}	

	return true, nil
}

func defineBonusPointsCap(level int) int {
	cap := 3
	for i := range level {
		if (i + 1) % 4 == 0 {
			cap += 2
		}
	}
	
	return cap
}

func validateBonusPoints(c NewCharacter) (bool, error) {
	characterLevel := c.Level
	bonusCap := defineBonusPointsCap(characterLevel)

	bonusPoints := [6]int{
		c.BonusPoints.Str,
		c.BonusPoints.Dex,
		c.BonusPoints.Con,
		c.BonusPoints.Int,
		c.BonusPoints.Wis,
		c.BonusPoints.Cha,
	}
	
	totalUsedBonus := 0

	for	_, bonus := range bonusPoints {
		totalUsedBonus += bonus	
	} 

	if totalUsedBonus > bonusCap {
		return false, fmt.Errorf("Invalid character ability scores. Used %d bonus points, which is higher than the %d cap", totalUsedBonus, bonusCap)
	}
	
	if bonusCap == 3 {
		countOnes := 0
		hasTwoBonus := false
		for _, bonus := range bonusPoints {
			if bonus == 1 {
				countOnes++
			} else if bonus == 2{
				hasTwoBonus = true
			}	
		}
		
		if countOnes != 1 {
			return false, fmt.Errorf("Invalid character ability scores. Used more than one +1 for a level %d character", characterLevel)
		} else if !hasTwoBonus {
			return false, fmt.Errorf("Invalid character ability scores. Did not use the +2 bonus point")
		}
	}

	return true, nil
}

func validateStatsRange(stats Attributes) (bool, error) {
	statsValues := [6]int{
		stats.Str,
		stats.Dex,
		stats.Con,
		stats.Int,
		stats.Wis,
		stats.Cha,
	}

	for _, value := range statsValues {
		if value > 30 || value < 8 {
			return false, fmt.Errorf("Invalid character ability score. Score %d is out of range 8 to 30", value)
		}
	}

	return true, nil
}

func mergePointBuyAndBonus(c NewCharacter) (Attributes, error) {
	validPointBuy, err := validatePointBuy(c)
	if err != nil {
		return Attributes{}, err
	}

	validBonus, err := validateBonusPoints(c)
	if err != nil {
		return Attributes{}, err
	}

	if !validPointBuy && !validBonus {
		return Attributes{}, fmt.Errorf("Something went wrong in mergePointBuyAndBonus()")
	}

	stats := Attributes{
		Str: c.PointBuy.Str + c.BonusPoints.Str,
		Dex: c.PointBuy.Str + c.BonusPoints.Str,
		Con: c.PointBuy.Str + c.BonusPoints.Str,
		Int: c.PointBuy.Str + c.BonusPoints.Str,
		Wis: c.PointBuy.Str + c.BonusPoints.Str,
		Cha: c.PointBuy.Str + c.BonusPoints.Str,
	}

	statsRangeOk, err := validateStatsRange(stats)
	if err != nil {
		return Attributes{}, err 
	} else if !statsRangeOk {
		return Attributes{}, fmt.Errorf("Something went wrong in mergePointBuyAndBonus()")
	}

	return stats, nil
}

func validateSkillsByReference(reference ClassInfo, c Character) (bool, error) {
	skillCap := reference.SkillChoice.Cap	
	referenceSkills := reference.SkillChoice.Possibilities
	userChoices := c.SkillChoice.Chosen

	countValidSkills := 0

	for _, chosenSkill := range userChoices {
		for _, skill := range referenceSkills {
			if chosenSkill == skill {
				countValidSkills++
			} 
		}
	}

	if countValidSkills != skillCap {
		return false, fmt.Errorf("Invalid skills. Chose a quantity of valid skills different from %d", skillCap)
	}

	return true, nil
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
