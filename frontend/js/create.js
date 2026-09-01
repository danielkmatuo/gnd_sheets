import * as api from "./api.js"
import * as creationRules from "./creationRules.js"

const classSelect = document.querySelector("#class");
const classInfo = document.querySelector("#class-info");

const raceSelect = document.querySelector("#race");

const characterStats = document.querySelectorAll("#abilities input");
const form = document.querySelector("#character-form"); //For now, I don't need this
const levelSelect = document.querySelector("#level");

const bonusButtons = document.querySelectorAll("#abilities button");
const resetBonusButton = document.querySelector("#reset-bonus");

let selectedSkills = ["something"];

let allocatedBonus = {
    "str": 0,
    "dex": 0,
    "con": 0,
    "int": 0,
    "wis": 0,
    "cha": 0
}

let validPointBuyState = {
    "str": 8,
    "dex": 8,
    "con": 8,
    "int": 8,
    "wis": 8,
    "cha": 8
};

const allSkillsData = await api.getSkillsInfoFromReference();

//render stuff on the frontend
function renderAbilities() {
    const stats = ["str", "dex", "con", "int", "wis", "cha"];
    let usedBonus = 0;

    for (const stat of stats) {
        document.querySelector(`#${stat}`).value = validPointBuyState[stat] + allocatedBonus[stat];
        usedBonus += allocatedBonus[stat];
    }

    document.querySelector("#curr-cost").textContent = `Current Cost: ${creationRules.calculateTotalCost(validPointBuyState)}`;
    document.querySelector("#curr-bonus").textContent = `Available Bonus Points: ${creationRules.calculateBonusPointsCap() - usedBonus}`;
}

function renderProficiencyBonus() {
    const characterLevel = Number(levelSelect.value);
    const proficiencyBonus = creationRules.calculateProficiencyBonus(characterLevel);

    if (proficiencyBonus < 2) {
        return
    }

    document.querySelector("#proficiency-bonus").textContent = "Proficiency Bonus: +" + proficiencyBonus;
}

function renderAllSkills(selectedSkills) {
    const proficiencyBonus = creationRules.calculateProficiencyBonus(Number(levelSelect.value));
    const skillsPanel = document.querySelector("#skills-panel");
    const mergedStats = creationRules.mergeBonusWithAbilityScores(validPointBuyState, allocatedBonus);
    const abilityModifiers = creationRules.calculateAbilityModifiers(mergedStats);
    const baselineSkillsObj = creationRules.createBaselineSkillsObj(allSkillsData, abilityModifiers);

    console.log(selectedSkills);

    let html = `
        <fieldset>
        <legend>Character Skills</legend>
    `;
    for (const key of Object.keys(baselineSkillsObj)) {
        let bonus = abilityModifiers[key];    
        for (const skill of Object.keys(baselineSkillsObj[key])) {
            let checked = false;
            for (const selected of selectedSkills) {
                if (skill === selected) {
                    checked = true; 
                    baselineSkillsObj[key][skill] = bonus + proficiencyBonus;
                    html += `
                        <p>${creationRules.capitalize(skill)}: ${creationRules.getModifierString(baselineSkillsObj[key][skill])}</p>
                    `;
                }
            }
            if (!checked) {
                baselineSkillsObj[key][skill] = bonus;
                html += `
                    <p>${creationRules.capitalize(skill)}: ${creationRules.getModifierString(baselineSkillsObj[key][skill])}</p>
                `;
            }
        }
    }

    html += `</fieldset>`;
    skillsPanel.innerHTML = html;
}

//changes data from index.html dynamically based on level passed by user
levelSelect.addEventListener("change", async function() {
    renderProficiencyBonus();
    renderAllSkills(selectedSkills);
});

//changes data from index.html dynamically based on class passed by user
//TODO: add tools list in reference data
classSelect.addEventListener("change", async function () {
    const selectedClass = classSelect.value;

    if (selectedClass === "") {
        classInfo.innerHTML = "";
        return;
    }

    const classData = await api.getClassInfoFromReference();

    var html = `
        <h2>${classData.name}</h2>
        <p>Hit Die: ${classData.hit_die}</p>
        <p>Saving Throws: ${classData.proficiencies.saving_throws.join(", ")}</p>
        <p>Armor proficiencies: ${classData.proficiencies.armor.join(", ")} </p>
        <p>Weapon proficiencies: ${classData.proficiencies.weapons.join(", ")}</p>
        <p>Tools proficiencies: ${classData.proficiencies.tools.join(", ")}</p>
        <fieldset>
            <legend>Choose ${classData.skill_choices.choose} skills:</legend>
    `;

    var numSkills = classData.skill_choices.from.length
    for (var i = 0; i < numSkills; i++) {
        html += `
        <div>
            <input type="checkbox" id="skill${i}" name="skills" value="${classData.skill_choices.from[i]}">
            <label for="skill${i}">${creationRules.capitalize(classData.skill_choices.from[i])}</label>
        </div>
        `;
    }
    
    html += `
    </fieldset>
    `;

    classInfo.innerHTML = html;

    const skillCheckboxes = document.querySelectorAll(
        'input[name="skills"]'
    );

    skillCheckboxes.forEach(function (checkbox) {
        checkbox.addEventListener("change", function (event) {
            selectedSkills = [];
            const inputedSkills = document.querySelectorAll(
                'input[name="skills"]:checked'
            );

            inputedSkills.forEach(function(skill) {
                selectedSkills.push(skill.value);
            });

            const quantity = selectedSkills.length;

            const valid = creationRules.quantitySkillsCheck(
                quantity,
                classData
            );

            if (!valid) {
                event.target.checked = false;

                alert(`You can only choose ${classData.skill_choices.choose} skills.`);
            }

            renderAllSkills(selectedSkills);
        });
    });
});

//validate character ability scores from frontend
characterStats.forEach(function(statInput) {
    statInput.addEventListener("change", async function(event){
        if (event.target.value === "") {
            event.target.value = validPointBuyState[event.target.id];
            return;
        }

        creationRules.tryPointBuyChange(event.target.id, Number(event.target.value), validPointBuyState);
        renderAbilities();
        renderAllSkills(selectedSkills);
    });
});

bonusButtons.forEach(function(button) {
    button.addEventListener("click", function(event) {
        const regex = /^(one|two|up|down).*(str|dex|con|int|wis|cha)$/;
        const match = event.target.id.match(regex);
        const stat = match[2];

        if (match[1] === "one") {
            creationRules.checkBothStates(1, stat, validPointBuyState[stat], validPointBuyState, allocatedBonus);
        }
        else if (match[1] === "two") {
            creationRules.checkBothStates(2, stat, validPointBuyState[stat], validPointBuyState, allocatedBonus);
        }
        else if (match[1] === "down") {
            creationRules.tryPointBuyChange(stat, validPointBuyState[stat] - 1, validPointBuyState, allocatedBonus);
        }
        else if (match[1] === "up") {
            creationRules.tryPointBuyChange(stat, validPointBuyState[stat] + 1, validPointBuyState, allocatedBonus);    
        }
        else {
            alert("Buttons werent found");
        }

        renderAbilities();
        renderAllSkills(selectedSkills);
    });
});

resetBonusButton.addEventListener("click", async function() {
    for (const key of Object.keys(allocatedBonus)) {
        allocatedBonus[key] = 0;
        document.querySelector(`#${key}`).value = validPointBuyState[key]; 
    }

    document.querySelector("#curr-bonus").textContent = `Available Bonus Points: ${creationRules.calculateBonusPointsCap()}`;
    renderAllSkills(selectedSkills);
});

//send form information to the go server and do the first validation of data
form.addEventListener("submit", async function(event) {
    event.preventDefault();

    const character = {
        "name": document.querySelector("#name").value,
        "level": Number(levelSelect.value), 
        "class": classSelect.value,
        "race": document.querySelector("#race").value,
        "skills": selectedSkills,
        "point_buy": validPointBuyState,
        "bonus_points": allocatedBonus
    }

    console.log(character);

    await api.sendCharacterData(character);
});
