const path = window.location.pathname;

async function getCharacterStructFromServer() {
    const id = path.match(/^\/character\/([^/]+)\/view$/)[1];
    const response = await fetch(`/character/${id}/data`);    
    if (!response.ok) {
        console.log("couldnt get the character data from server")
        return
    }

    const data = await response.json();

    console.log(JSON.stringify(data, null, 2));

    return data;
}

//basic info
const characterName = document.querySelector("#name");
const characterClass = document.querySelector("#class");
const characterRace = document.querySelector("#race");
const characterLevel = document.querySelector("#level");

//combat info
const characterHpPointer = document.querySelector("#hp");
const characterAc = document.querySelector("#ac");
const characterSpeed = document.querySelector("#speed");

//proficiencies
const characterArmorProficiecy = document.querySelector("#armor");
const characterWeaponProficiency = document.querySelector("#weapons");
const characterToolsProficiency = document.querySelector("#tools");
const characterSkills = document.querySelector("#skills");
const characterSavingThrows = document.querySelector("#saving");

function populateViewPage(data) {
    characterName.textContent = "Name: " + data.name;
    characterClass.textContent = "Class: " + data.class;
    characterRace.textContent = "Race: " + data.race;
    characterLevel.textContent = "Level: " + data.level;
    
    characterHpPointer.textContent = "HP: " + data.current_hp + "/" + data.max_hp;
    characterAc.textContent = "AC: " + data.ac;
    characterSpeed.textContent = "Speed: " + data.speed;

    characterArmorProficiecy.textContent = "Armor proficiency: " + data.proficiencies.armor.join(", ");
    characterWeaponProficiency.textContent = "Weapon proficiency: " + data.proficiencies.weapons.join(", ");
    characterToolsProficiency.textContent = "Tool proficiency: " + data.proficiencies.tools.join(", ");
    characterSkills.textContent = "Skills proficiency: " + data.skill_choices.skills.join(", ");
    characterSavingThrows.textContent = "Saving throws proficiency: " + data.proficiencies.saving_throws.join(", ");
}

async function redirectToViewPage() {
    const data = await getCharacterStructFromServer();

    populateViewPage(data);
}

redirectToViewPage();
