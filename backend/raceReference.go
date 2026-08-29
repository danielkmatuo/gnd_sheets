package backend

import "net/http"

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

type DragonbornAncestry struct {
	Type string `json:"type"`
	BreathRange int `json:"breath_range"`
	BreathShape string `json:"breath_shape"`
	BreathSavingThrow string `json:"breat_saving_throw"`
}

type Trait struct {

}

func raceReferenceHandler(w http.ResponseWriter, r *http.Request) {
	return	
}
