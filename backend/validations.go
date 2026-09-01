package backend

import (
	"fmt"
	"math"
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

	stats, bonusPoints, err := mergePointBuyAndBonus(newCharacter)
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
	c.BonusPoints = bonusPoints

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
	classInfoMap, err := getClassReferenceData()	
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
	referenceMap, err := getClassReferenceData()
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

//TODO: change this whenever I get into feats implementation
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

func mergePointBuyAndBonus(c NewCharacter) (Attributes, Attributes, error) {
	validPointBuy, err := validatePointBuy(c)
	if err != nil {
		return Attributes{}, Attributes{}, err
	}

	validBonus, err := validateBonusPoints(c)
	if err != nil {
		return Attributes{}, Attributes{}, err
	}

	if !validPointBuy && !validBonus {
		return Attributes{}, Attributes{}, fmt.Errorf("Something went wrong in mergePointBuyAndBonus()")
	}

	stats := Attributes{
		Str: c.PointBuy.Str + c.BonusPoints.Str,
		Dex: c.PointBuy.Dex + c.BonusPoints.Dex,
		Con: c.PointBuy.Con + c.BonusPoints.Con,
		Int: c.PointBuy.Int + c.BonusPoints.Int,
		Wis: c.PointBuy.Wis + c.BonusPoints.Wis,
		Cha: c.PointBuy.Cha + c.BonusPoints.Cha,
	}

	statsRangeOk, err := validateStatsRange(stats)
	if err != nil {
		return Attributes{}, Attributes{}, err 
	} else if !statsRangeOk {
		return Attributes{}, Attributes{}, fmt.Errorf("Something went wrong in mergePointBuyAndBonus()")
	}

	return stats, c.BonusPoints, nil
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
