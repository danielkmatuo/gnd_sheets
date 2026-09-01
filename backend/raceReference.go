package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

//Generic race
type Race struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Speed int `json:"speed"`
	KnownLanguages []string `json:"known_languages"`
	AdditionalLanguages int `json:"additional_languages"`
	Traits map[string]Trait `json:"traits"`
	Subraces map[string]map[string]Trait `json:"subraces"`
}

//Races with special mechanics
type Dragonborn struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Speed int `json:"speed"`
	KnownLanguages []string `json:"known_languages"`
	AdditionalLanguages int `json:"additional_languages"`
	Traits WrapperDragonbornTraits `json:"traits"`
	Ancestry map[string]DragonbornAncestry `json:"ancestry"`
}

//Generic trait
type Trait struct {
	Description string `json:"description"`

	Uses int `json:"uses"`
	Reset string `json:"reset"`

	HpTrigger int `json:"hp_trigger"`
	SetHp int `json:"set_hp"`

	ResistTypes []string `json:"resist_types"`

	Choose int `json:"choose"`
	From []string `json:"choices"`
	ChoiceType string `json:"choice_type"`
	StatBonus int `json:"stat_bonus"`

	TraitChoices map[string]Trait `json:"trait_choices"`

	BonusSpells BonusSpellsTrait `json:"bonus_spells"`	
	SpellModifier string `json:"spellcasting_stat"`

	ProficienciesGranted []string `json:"extra_proficiencies"`
	ProficienciesType string `json:"type_proficiency"`

	Bonus int `json:"bonus"`
	BonusType string `json:"bonus_type"`
}

//Traits with weird mechanics, so its better to break them down into separated structs
type BonusSpellsTrait map[string]SpellFromTrait

type SpellFromTrait struct {
	Level int `json:"level"`
	Uses int `json:"uses"`
	Reset string `json:"reset"`
}

//Wrappers for the special mechanics races with their respective traits
type WrapperDragonbornTraits struct {
	AncestryDescription Trait `json:"Draconic Ancestry"`
	BreathWeapon BreathWeaponTrait `json:"Breath Weapon"`
	Resist Trait `json:"Damage Resistance"`
}

//Traits with special mechanics
type DragonbornAncestry struct {
	Type string `json:"type"`
	BreathRange int `json:"breath_range"`
	BreathShape string `json:"breath_shape"`
	BreathSavingThrow string `json:"breath_saving_throw"`
}

type BreathWeaponTrait struct {
	Description string `json:"description"`
	Damage map[int]string `json:"damage"`
	BaselineDC int `json:"baseline_dc"`
	ModifierDC string `json:"modifier_dc"`
	Reset string `json:"reset"`
	Uses int `json:"uses"`
}

func raceReferenceHandler(w http.ResponseWriter, r *http.Request) {
	userRace := r.PathValue("race")

	if userRace == "dragonborn" {
		referenceDragonborn, err := getDragonbornReferenceData(userRace)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-type", "application/json")

		err = json.NewEncoder(w).Encode(referenceDragonborn)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		return
	}

	referenceRace, err := getRaceReferenceData(userRace)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(referenceRace)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
}

func getRaceReferenceData(userRace string) (Race, error) {
	root := findRootDir()
	if root == "" {
		return Race{}, fmt.Errorf("couldnt find root")
	}

	referenceDataPath := filepath.Join(root, "data", "reference", "races.json")

	referenceBytes, err := os.ReadFile(referenceDataPath)
	if err != nil {
		return Race{}, err
	}

	var raceMap map[string]Race
	err = json.Unmarshal(referenceBytes, &raceMap)
	if err != nil {
		return Race{}, err
	}

	race, ok := raceMap[userRace]
	if !ok {
		return Race{}, fmt.Errorf("couldnt find chosen race")
	}

	return race, nil
}

func getDragonbornReferenceData(userRace string) (Dragonborn, error) {
	root := findRootDir()
	if root == "" {
		return Dragonborn{}, fmt.Errorf("couldnt find root")
	}

	referenceDataPath := filepath.Join(root, "data", "reference", "races.json")

	referenceBytes, err := os.ReadFile(referenceDataPath)
	if err != nil {
		return Dragonborn{}, err
	}

	var dragonbornMap map[string]Dragonborn
	err = json.Unmarshal(referenceBytes, &dragonbornMap)
	if err != nil {
		return Dragonborn{}, err
	}

	dragonborn, ok := dragonbornMap[userRace]
	if !ok {
		return Dragonborn{}, fmt.Errorf("couldnt find chosen race")
	}

	return dragonborn, nil
}
