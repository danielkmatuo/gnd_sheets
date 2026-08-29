const classSelect = document.querySelector("#class");
const classInfo = document.querySelector("#class-info");

const raceSelect = document.querySelector("#race");

const characterStats = document.querySelectorAll("#abilities input");
const form = document.querySelector("#character-form"); //For now, I don't need this
const levelSelect = document.querySelector("#level");

let bonusPointsCap = 3; //TODO: change this baseline whenever I add feats options on character creation
const bonusButtons = document.querySelectorAll("#abilities button");
const resetBonusButton = document.querySelector("#reset-bonus");

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

//fetches the reference data from the backend
async function getClassInfoFromReference (selectedClass) {
    const response = await fetch(
        `/reference/classes/${selectedClass}`
    );

    if (!response.ok) {
        alert("Could not load class information from server.");
        const errorMessage = await response.text();
        console.error("Go server error: ", errorMessage)
        return;
    }

    return response;
}

async function getRaceInfoFromReference(selectedRace) {
    const response = await fetch(
        `/reference/races/${selectedRace}`
    );

    if (!response.ok) {
        alert("Could not load race information from server.");
        const errorMessage = await response.text();
        console.error("Go server error: ", errorMessage)
        return;
    }

    return response;
}

async function getLanguagesInfoFromReference(selectedRace) {
    const response = await fetch(
        `/reference/languages`
    );

    if (!response.ok) {
        alert("Could not load languages information from server.");
        const errorMessage = await response.text();
        console.error("Go server error: ", errorMessage)
        return;
    }

    return response;
}

async function getSkillsInfoFromReference(selectedRace) {
    const response = await fetch(
        `/reference/skills`
    );

    if (!response.ok) {
        alert("Could not load skills information from server.");
        const errorMessage = await response.text();
        console.error("Go server error: ", errorMessage)
        return;
    }

    return response;
}

//make user able to only choose skills equal to the number of skills cap
function quantitySkillsCheck (quantity, data) {
    const skillIssues = data.skill_choices.choose;
    if (quantity > skillIssues) {
        return false;
    }

    return true;
}

//TODO: calculate proficiency bonus based on level selected by the user
function calculateProficiencyBonus() {
    const level = Number(levelSelect.value);
    const baselineBonus = 2;
    const progressionStep = 4; //hardcoded for now

    if (level < 1 || level > 20) {
        alert("Invalid level. Must be in range 1 to 20");
        return -1;
    }


}

//frontend validation of the character ability scores
function instantiateCostsMap() {
    const scoresCosts = new Map(); 

    scoresCosts.set(8, 0);
    scoresCosts.set(9, 1);
    scoresCosts.set(10, 2);
    scoresCosts.set(11, 3);
    scoresCosts.set(12, 4);
    scoresCosts.set(13, 5);
    scoresCosts.set(14, 7);
    scoresCosts.set(15, 9);

    return scoresCosts;
}

function calculateTotalCost(attemptedState) {
    const costMap = instantiateCostsMap();
    let currCost = 0;
   
    for (const value of Object.values(attemptedState)) {
        if (value < 8 || value > 15) {
            alert("Invalid ability scores. The point buy scores must be within range 8 to 15");
            return -1;
        }

        currCost += costMap.get(value);
    }

    return currCost;
}

function validateAbilityCost(characterCost) {
    if (characterCost > 27) {
        alert("Invalid character. Total ability scores cost is above 27 or below 0");
        return false;
    }
    else if (characterCost < 0) {
        return false;
    }

    return true;
}

function validateCreationBonusPoints(allocationAttempt) {
    let bonusSum = 0;

    for (const value of Object.values(allocationAttempt)) {
        bonusSum += value;
    }

    if (bonusSum > bonusPointsCap) {
        return false;
    }

    return true;
}

function tryAllocationBonusPoints(bonus, currStat) {
    const bonusKeys = [
        "str",
        "dex",
        "con",
        "int",
        "wis",
        "cha"
    ];
    const attemptedAllocation = {...allocatedBonus};
    attemptedAllocation[currStat] += bonus;

    let countOnes = 0;

    if (!validateCreationBonusPoints(attemptedAllocation)) {
        alert("Invalid bonus point allocation. Can only allocate points until budget limit");
        return false;
    }

    if (bonus > 2 || bonus < 0) {
        alert("Invalid bonus point allocation. Can only allocate +2 or +1");
        return false;
    }

    if (attemptedAllocation[currStat] > 2) {
        alert("Invalid bonus point allocation. Same stat received both +1 and +2");
        return false;
    }

    for (let i = 0; i < bonusKeys.length; i++) {
        if (attemptedAllocation[bonusKeys[i]] === 1) {
            countOnes++;
        }
    }

    if (countOnes > 1) {
        alert("Invalid bonus point allocation. Character has more than one stat with a +1");
        return false;
    }

    allocatedBonus = attemptedAllocation;

    return true;
}

function tryPointBuyChange(stat, newValue) {
    const attempt = {...validPointBuyState};

    attempt[stat] = newValue;

    const cost = calculateTotalCost(attempt);

    if (!validateAbilityCost(cost)) {
        return false;
    }

    validPointBuyState = attempt;

    return true;
}

function checkBothStates(bonus, stat, newValue) {
    const pointBuyAttemptOk = tryPointBuyChange(stat, newValue);

    if (pointBuyAttemptOk) {
        const allocationOk = tryAllocationBonusPoints(bonus, stat);
        if (!allocationOk) {
            return false;
        }

        return true;
    }

    return false;
}

//render stuff on the frontend
function renderAbilities() {
    const stats = ["str", "dex", "con", "int", "wis", "cha"];
    let usedBonus = 0;

    for (const stat of stats) {
        document.querySelector(`#${stat}`).value = validPointBuyState[stat] + allocatedBonus[stat];
        usedBonus += allocatedBonus[stat];
    }

    document.querySelector("#curr-cost").textContent = `Current Cost: ${calculateTotalCost(validPointBuyState)}`;
    document.querySelector("#curr-bonus").textContent = `Available Bonus Points: ${bonusPointsCap - usedBonus}`;
}

//changes data from index.html dynamically
//TODO: add tools list in reference data
classSelect.addEventListener("change", async function () {
    const selectedClass = classSelect.value;

    if (selectedClass === "") {
        classInfo.innerHTML = "";
        return;
    }

    const response = await getClassInfoFromReference(selectedClass);

    const classData = await response.json();

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

    var numSkills = data.skill_choices.from.length
    for (var i = 0; i < numSkills; i++) {
        html += `
        <div>
            <input type="checkbox" id="skill${i}" name="skills" value="${classData.skill_choices.from[i]}">
            <label for="skill${i}">${classData.skill_choices.from[i]}</label>
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
            const selectedSkills = document.querySelectorAll(
                'input[name="skills"]:checked'
            );

            const quantity = selectedSkills.length;

            const valid = quantitySkillsCheck(
                quantity,
                data
            );

            if (!valid) {
                event.target.checked = false;

                alert(`You can only choose ${data.skill_choices.choose} skills.`);
            }
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

        tryPointBuyChange(event.target.id, Number(event.target.value));
        renderAbilities();
    });
});

bonusButtons.forEach(function(button) {
    button.addEventListener("click", function(event) {
        const regex = /^(one|two|up|down).*(str|dex|con|int|wis|cha)$/;
        const match = event.target.id.match(regex);
        const stat = match[2];

        if (match[1] === "one") {
            checkBothStates(1, stat, validPointBuyState[stat]);
        }
        else if (match[1] === "two") {
            checkBothStates(2, stat, validPointBuyState[stat]);
        }
        else if (match[1] === "down") {
            tryPointBuyChange(stat, validPointBuyState[stat] - 1);
        }
        else if (match[1] === "up") {
            tryPointBuyChange(stat, validPointBuyState[stat] + 1);    
        }
        else {
            alert("Buttons werent found");
        }

        renderAbilities();
    });
});

resetBonusButton.addEventListener("click", async function() {
    creationBonusPoints = 0;
    for (const key of Object.keys(allocatedBonus)) {
        allocatedBonus[key] = 0;
        document.querySelector(`#${key}`).value = validPointBuyState[key]; 
    }

    document.querySelector("#curr-bonus").textContent = `Available Bonus Points: ${bonusPointsCap}`;
});

//send form information to the go server and do the first validation of data
form.addEventListener("submit", async function(event) {
    event.preventDefault();

    const skills = Array.from(
        document.querySelectorAll('input[name="skills"]:checked')
    ).map(input => input.value);

    character = {
        "name": document.querySelector("#name").value,
        "level": Number(document.querySelector("#level").value), 
        "class": classSelect.value,
        "race": document.querySelector("#race").value,
        "skills": skills,
        "point_buy": validPointBuyState,
        "bonus_points": allocatedBonus
    }

    console.log(character);

    const response = await fetch("/characters/create", {
        method: "POST",

        headers: {
            "Content-Type": "application/json"
        },

        body: JSON.stringify(character)
    });

    if (!response.ok) {
        const errorMessage = await response.text();
        console.error("Go server error:", errorMessage);
        return;
    }

    window.location.href = "/characters";
});
