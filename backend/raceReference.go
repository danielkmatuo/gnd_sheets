package backend

import "net/http"

//Generic race
type Race struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Speed int `json:"speed"`
	KnownLanguages []string `json:"known_languages"`
	AdditionalLanguages int `json:"additional_languages"`
	Traits map[string]Trait `json:"traits"`
	Subraces map[string]map[string]Trait `json:"subraces"`
	Ancestry map[string]DragonbornAncestry `json:"ancestry"`
}

//Races with special mechanics
type DragonbornAncestry struct {
	Type string `json:"type"`
	BreathRange int `json:"breath_range"`
	BreathShape string `json:"breath_shape"`
	BreathSavingThrow string `json:"breath_saving_throw"`
}

//Generic trait
type Trait struct {

}

//Traits with special mechanics
type BreathWeapon struct {
	Description string `json:"description"`
	Damage map[int]string `json:"damage"`
	BaselineDC int `json:"baseline_dc"`
	ModifierDC string `json:"modifier_dc"`
	Reset int `json:"reset"`
	Uses int `json:"uses"`
}

type RelentlessEndurance struct {
	Description string `json:"description"`
	Reset string `json:"reset"`
	Uses int `json:"uses"`
} 

type DwarvenThoughness struct {
	Description string `json:"description"`
	ExtraHp int `json:"extra_hp"`
}

func raceReferenceHandler(w http.ResponseWriter, r *http.Request) {
	return	
}
