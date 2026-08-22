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
		return Character{}, fmt.Errorf("character must have exact %d", referenceCharacter.SkillChoice.Cap)
	} 

	if c.Level == 1 {
		c.MaxHp, err = parseHitDie(referenceCharacter.HitDie)
		if err != nil {
			return Character{}, err
		}
	} else {
		hitDieValue, err := parseHitDie(referenceCharacter.HitDie)
		if err != nil {
			return Character{}, err
		}

		conBonus := calculateStatBonus(c.Stats.Con)
		hitDieStep := int(math.Ceil(float64(hitDieValue / 2)))
		c.MaxHp = hitDieValue + hitDieStep * (c.Level - 1) + conBonus * (c.Level - 1)
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

	if newCharacter.Stats.Str < 8 || newCharacter.Stats.Str > 15 {
		return Character{}, fmt.Errorf("invalid STR: must be within 8 to 15 range")
	}
	
	if newCharacter.Stats.Dex < 8 || newCharacter.Stats.Dex > 15 {
		return Character{}, fmt.Errorf("invalid DEX: must be within 8 to 15 range")
	}
	
	if newCharacter.Stats.Con < 8 || newCharacter.Stats.Con > 15 {
		return Character{}, fmt.Errorf("invalid CON: must be within 8 to 15 range")
	}
	
	if newCharacter.Stats.Int < 8 || newCharacter.Stats.Int > 15 {
		return Character{}, fmt.Errorf("invalid INT: must be within 8 to 15 range")
	}
	
	if newCharacter.Stats.Wis < 8 || newCharacter.Stats.Wis > 15 {
		return Character{}, fmt.Errorf("invalid WIS: must be within 8 to 15 range")
	}
	
	if newCharacter.Stats.Cha < 8 || newCharacter.Stats.Cha > 15 {
		return Character{}, fmt.Errorf("invalid CHA: must be within 8 to 15 range")
	}

	if newCharacter.Skills == nil {
		return Character{}, fmt.Errorf("cannot choose 0 skills")
	}

	abilityCostOk := validateAbilityCost(newCharacter.Stats)	
	if !abilityCostOk {
		return c, fmt.Errorf("ability cost cannot be negative neither above 27")
	}

	c.ID = characterID
	c.Name = newCharacter.Name
	c.Class = newCharacter.Class
	c.Race = newCharacter.Race
	c.Level = newCharacter.Level
	c.SkillChoice.Chosen = newCharacter.Skills
	c.Stats = newCharacter.Stats

	return c, nil
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
	costMap := make(map[int]int)
	allowedStatsValues := [8]int{8, 9, 10, 11, 12, 13, 14, 15}
	allowedCostsValues := [8]int{0, 1, 2, 3, 4, 5, 7, 9}

	for i, stat := range allowedStatsValues{
		costMap[stat] = allowedCostsValues[i]	
	}

	currCost := 0

	for _, value := range statsValues {
		currCost += costMap[value]	
	}

	if currCost > 27 || currCost < 0 {
		return false
	}

	return true
}

func calculateStatBonus(stat int) int {
	baselineStatValue := 10
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

	possibleConBonus := [11]int {-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5}
	referenceHitDieStep := int(math.Ceil(float64(referenceHitDieValue / 2))) 
	var referenceMaxHp int

	if c.Level == 1 {
		referenceMaxHp = referenceHitDieValue
		if referenceMaxHp != c.MaxHp {
			return Character{}, fmt.Errorf("character has invalid max hp value")
		}
	}
	
	areEqual := false

	for i := 0; i < len(possibleConBonus); i++ {
		referenceMaxHp = referenceHitDieValue + referenceHitDieStep * (c.Level - 1) + possibleConBonus[i] * (c.Level - 1)
		if referenceMaxHp == c.MaxHp {
			areEqual = true
			break
		} 
	}
	
	if !areEqual {
		return Character{}, fmt.Errorf("character has invalid max hp value")
	}

	c.CurrHp = referenceMaxHp
	c.AC = 10 + calculateStatBonus(c.Stats.Dex)
	c.Speed = 9.0

	return c, nil
}
