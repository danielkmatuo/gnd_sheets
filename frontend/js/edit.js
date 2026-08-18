const path = window.location.pathname;
const characterName = document.querySelector("#name");
const characterClass = document.querySelector("#class");
const characterRace = document.querySelector("#race");
const characterSkills = document.querySelector("#skills");

async function getCharacterStructFromServer() {
    const id = path.match(/^\/character\/([^/]+)\/edit$/);
    const response = await fetch(`/character/${id}/data`);    
    const data = await response.json();

    return data;
}

getCharacterStructFromServer();
