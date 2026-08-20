const path = window.location.pathname;

async function getCharacterStructFromServer() {
    const id = path.match(/^\/character\/([^/]+)\/data$/)[1];
    const response = await fetch(`/character/${id}/data`);    
    if (!response.ok) {
        console.log("couldnt get the character data from server")
        return
    }

    const data = await response.json();

    return data;
}

const characterName = document.querySelector("#name");
const characterClass = document.querySelector("#class");
const characterRace = document.querySelector("#race");
const characterSkills = document.querySelector("#skills");

function populateViewPage(data) {
    characterName.textContent = data.name
    characterClass.textContent = data.class
    characterRace.textContent = data.race
    characterSkills.textContent = data.skill_choices.skills.join(", ")
}

async function redirectToViewPage() {
    const data = await getCharacterStructFromServer();

    populateViewPage(data);
}

redirectToViewPage();
