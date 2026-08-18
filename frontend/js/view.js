const path = window.location.pathname;

async function getCharacterStructFromServer() {
    const id = path.match(/^\/character\/([^/]+)\/edit$/);
    const response = await fetch(`/character/${id}/data`);    
    const data = await response.json();

    return data;
}

const characterName = document.querySelector("#name");
const characterClass = document.querySelector("#class");
const characterRace = document.querySelector("#race");
const characterSkills = document.querySelector("#skills");

async function populateViewPage() {
    window.location.href = "/character/view"

    const data = await getCharacterStructFromServer();

    characterName.value = data.name
    characterClass.value = data.class
    characterRace.value = data.race
    characterSkills.value = data.skill_choices.skills.join(", ")
}

populateViewPage();
