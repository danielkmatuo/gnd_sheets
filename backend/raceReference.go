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

	CurrHp int `json:"curr_hp"`

	ResistTypes []string `json:"resist_types"`

	Choose int `json:"choose"`
	From []string `json:"choices"`
	ChoiceType string `json:"choice_type"`
	StatBonus int `json:"stat_bonus"`

	BonusSpells BonusSpellsTrait `json:"bonus_spells"`	
	SpellModifier string `json:"spellcasting_stat"`

	ProficienciesGranted []string `json:"extra_proficiencies"`
	ProficienciesType string `json:"type_proficiency"`

	Bonus int `json:"bonus"`
	BonusType string `json:"bonus_type"`
}

//Traits with weird mechanics, so its better to break them down into separated structs
type BonusSpellsTrait struct {
	Spells map[string]SpellFromTrait `json:"bonus_spells"`
}

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
	Reset int `json:"reset"`
	Uses int `json:"uses"`
}

func raceReferenceHandler(w http.ResponseWriter, r *http.Request) {
	return	
}
