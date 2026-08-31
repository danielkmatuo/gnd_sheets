package backend

import(
	"strconv"	
	"strings"
	"fmt"
	"math"
)

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

func calculateBonusProficiency(level int) (int, error) {
	if level < 1 || level > 20 {
		return -10, fmt.Errorf("Invalid level. Must be in range ")
	}

	const progressionStep = 4
	baselineBonus := 2
	progressionAnchor := 2

	for i := 5; i <= level; i++ {
		if i == progressionAnchor {
			baselineBonus++
		} else if (i - progressionAnchor == progressionStep) {
			progressionAnchor = i
			baselineBonus++
		}
	}

	return baselineBonus, nil
}

func calculateDynamicallyNewCharacter(c Character) (Character, error) {
	referenceMap, err := getClassReferenceData()
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

	statsBonus := calculateStatBonusStruct(c.Stats)

	proficiencyBonus, err := calculateBonusProficiency(c.Level)
	if err != nil {
		return Character{}, err
	}

	c.CurrHp = referenceMaxHp
	c.AC = 10 + calculateStatBonus(c.Stats.Dex)
	c.Speed = 9.0
	c.AbilitiesModifiers = statsBonus
	c.ProficiencyBonus = proficiencyBonus

	return c, nil
}
