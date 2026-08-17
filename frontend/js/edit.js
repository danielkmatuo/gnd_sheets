const path = window.location.pathname;
const characterName = document.querySelector("#name");
const characterClass = document.querySelector("#class");
const characterRace = document.querySelector("#race");
const characterSkills = document.querySelector("#skills");

async function getCharacterStructFromServer() {
    parts = path.split("/");
    id = parts[2];
    const response = await fetch("/character/${id}/edit");    
    const data = await response.json();

    characterName.value = data.name;
    characterClass.value = data.class;
    characterRace.value = data.race;
    characterSkills.value = data.skill_choices.skills.join(", ");
}
