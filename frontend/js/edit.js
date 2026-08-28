function extractIdFromRoute() {
    const routeRegex = /^\/character\/([^/]+)\/edit$/;
    const match = window.location.pathname.match(routeRegex);

    return match[1];
}

function returnStatBounsString(bonus) {
    if (Number(bonus) >= 0) {
        return `+${bonus}`;
    }

    return `${bonus}`;
}

async function getCharacterStructFromServer() {
    const id = extractIdFromRoute();
    const response = await fetch(`/character/${id}/data`);    
    const data = await response.json();

    return data;
}

//basic info
const characterName = document.querySelector("#name");
const characterClass = document.querySelector("#class");
const characterRace = document.querySelector("#race");
const characterLevel = document.querySelector("#level");

//stats info
const characterStr = document.querySelector("#str");
const characterDex = document.querySelector("#dex");
const characterCon = document.querySelector("#con");
const characterInt = document.querySelector("#int");
const characterWis = document.querySelector("#wis");
const characterCha = document.querySelector("#cha");

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

//redirect to edit div
const divRedirects = document.querySelector("#redirects");


function populateEditPage(data) {
    characterName.textContent = "Name: " + data.name;
    characterClass.textContent = "Class: " + data.class;
    characterRace.textContent = "Race: " + data.race;
    characterLevel.textContent = "Level: " + data.level;

    characterStr.textContent = "Strengh: " + data.stats.str + ` (${returnStatBounsString(data.stats_bonus.str)})`;
    characterDex.textContent = "Dexterity: " + data.stats.dex + ` (${returnStatBounsString(data.stats_bonus.dex)})`;
    characterCon.textContent = "Constitution: " + data.stats.con + ` (${returnStatBounsString(data.stats_bonus.con)})`;
    characterInt.textContent = "Intelligence: " + data.stats.int + ` (${returnStatBounsString(data.stats_bonus.int)})`;
    characterWis.textContent = "Wisdom: " + data.stats.wis + ` (${returnStatBounsString(data.stats_bonus.wis)})`;
    characterCha.textContent = "Charisma: " + data.stats.cha + ` (${returnStatBounsString(data.stats_bonus.cha)})`;
    
    characterHpPointer.textContent = "HP: " + data.current_hp + "/" + data.max_hp;
    characterAc.textContent = "AC: " + data.ac;
    characterSpeed.textContent = "Speed: " + data.speed;

    characterArmorProficiecy.textContent = "Armor proficiency: " + data.proficiencies.armor.join(", ");
    characterWeaponProficiency.textContent = "Weapon proficiency: " + data.proficiencies.weapons.join(", ");
    characterToolsProficiency.textContent = "Tool proficiency: " + data.proficiencies.tools.join(", ");
    characterSkills.textContent = "Skills proficiency: " + data.skill_choices.skills.join(", ");
    characterSavingThrows.textContent = "Saving throws proficiency: " + data.proficiencies.saving_throws.join(", ");

    const id = extractIdFromRoute();
    const html = `
        <a href="/character/${id}/view">Back</a>
        <a href="/characters">View all characters</a>
    `;

    divRedirects.innerHTML = html;
}

async function getDataAndPopulate() {
    const data = await getCharacterStructFromServer();

    populateEditPage(data);
}

getDataAndPopulate();
