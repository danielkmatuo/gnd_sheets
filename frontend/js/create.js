const classSelect = document.querySelector("#class");
const classInfo = document.querySelector("#class-info");
const characterStats = document.querySelectorAll("#abilities input");
const form = document.querySelector("#character-form");
const levelSelect = document.querySelector("#level");

const characterStr = Number(document.querySelector("#str").value);
const characterDex = Number(document.querySelector("#dex").value);
const characterCon = Number(document.querySelector("#con").value);
const characterInt = Number(document.querySelector("#int").value);
const characterWis = Number(document.querySelector("#wis").value);
const characterCha = Number(document.querySelector("#cha").value);

let lastCharacterCost = 0;
let creationBonusPoints = 0;
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

//make user able to only choose skills equal to the number of skills cap
function quantitySkillsCheck (quantity, data) {
    const skillIssues = data.skill_choices.choose;
    if (quantity > skillIssues) {
        return false;
    }

    return true;
}

async function getClassInfoFromReference (selectedClass) {
    const response = await fetch(
        `/reference/classes/${selectedClass}`
    );

    if (!response.ok) {
        classInfo.textContent = "could not load class informataion from server.";
        return;
    }

    return response;
}

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
        console.log(value);
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

function validateCreationBonusPoints(bonus) {
    if (creationBonusPoints + bonus > bonusPointsCap) {
        return false;
    }

    return true;
}

//TODO: Find a way to make this code ensure the correct structure for the initial bonus points allocation (one stat +2 and another +1)
function validateAllocationCreationBonusPoints(bonus, currStat) {
    const bonusKeys = [
        "str",
        "dex",
        "con",
        "int",
        "wis",
        "cha"
    ];

    let countOnes = 0;

    if (bonus > 2 || bonus < 0) {
        alert("Invalid bonus point allocation. Can only allocate +2 or +1");
        return false;
    }

    let currStatBonus = 0;

    if (allocatedBonus[currStat] > 2) {
        alert("Invalid bonus point allocation. Same stat received both +1 and +2");
        return false;
    }

    for (let i = 0; i < bonusKeys.length; i++) {
        if (allocatedBonus[bonusKeys[i]] === 1) {
            countOnes++;
        }
    }

    if (countOnes > 1) {
        alert("Invalid bonus point allocation. Character has more than one stat with a +1");
        return false;
    }

    return true;
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

    const data = await response.json();

    var html = `
        <h2>${data.name}</h2>
        <p>Hit Die: ${data.hit_die}</p>
        <p>Saving Throws: ${data.proficiencies.saving_throws.join(", ")}</p>
        <p>Armor proficiencies: ${data.proficiencies.armor.join(", ")} </p>
        <p>Weapon proficiencies: ${data.proficiencies.weapons.join(", ")}</p>
        <p>Tools proficiencies: ${data.proficiencies.tools.join(", ")}</p>
        <fieldset>
            <legend>Choose ${data.skill_choices.choose} skills:</legend>
    `;

    var numSkills = data.skill_choices.from.length
    for (var i = 0; i < numSkills; i++) {
        html += `
        <div>
            <input type="checkbox" id="skill${i}" name="skills" value="${data.skill_choices.from[i]}">
            <label for="skill${i}">${data.skill_choices.from[i]}</label>
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
//TODO: refactor to fit the new funcs to validate the character stats
characterStats.forEach(function(statInput) {
    statInput.addEventListener("change", async function(event){
        if (event.target.value === "") {
            event.target.value = validPointBuyState[event.target.id];
            return;
        }

        let currentAttempt = {...validPointBuyState};

        currAttempt[event.target.id] = Number(event.target.value);
        const currCost = calculateTotalCost(currAttempt);
        const isValidCost = validateAbilityCost(currCost);

        if (isValidCost) {
            lastCharacterCost = currCost;
            validPointBuyState = currAttempt;
            document.querySelector("#curr-cost").textContent = "Current Cost: " + lastCharacterCost;
        }
    });
});

//gigantic mess of code to control the buttons behaviour
bonusButtons.forEach(function(button) {
    button.addEventListener("click", function(event) {
        const regex = /^(one|two|up|down).*(str|dex|con|int|wis|cha)$/;
        const match = event.target.id.match(regex);

        let currCost = 0;

        if (match[1] === "one") {
            const creationBonusPointsOk = validateCreationBonusPoints(1);
            if (!creationBonusPointsOk) {
                alert("Already spent all bonus points for this character");
                return;
            }

            creationBonusPoints++;
            document.querySelector(`#${match[2]}`).value++;

            allocatedBonus[match[2]]++;
            const bonusAllocationOk = validateAllocationCreationBonusPoints(1, match[2]);

            if (!bonusAllocationOk) {
                allocatedBonus[match[2]]--;
                creationBonusPoints--;
                document.querySelector(`#${match[2]}`).value--;
                return; 
            } 
        }
        else if (match[1] == "two") {
            const creationBonusPointsOk = validateCreationBonusPoints(2);
            if (!creationBonusPointsOk) {
                alert("Already spent all bonus points for this character");
                return;
            }

            creationBonusPoints += 2;
            document.querySelector(`#${match[2]}`).value = Number(document.querySelector(`#${match[2]}`).value) + 2;

            allocatedBonus[match[2]] += 2;
            const bonusAllocationOk = validateAllocationCreationBonusPoints(2, match[2]);

            if (!bonusAllocationOk) {
                allocatedBonus[match[2]] -= 2;
                creationBonusPoints -= 2;
                document.querySelector(`#${match[2]}`).value = Number(document.querySelector(`#${match[2]}`).value) - 2;
                return; 
            }
        }
        else if (match[1] === "up") {
            currAttempt = validPointBuyState;
            currAttempt[match[2]]++;

            currCost = calculateTotalCost(currAttempt)
            const isValidCost = validateAbilityCost(currCost);

            if (!isValidCost) {
                return;
            }

            validPointBuyState = currAttempt;
            document.querySelector(`#${match[2]}`).value++;
            lastCharacterCost = currCost;
        }
        else if (match[1] === "down") {
            currAttempt = {...validPointBuyState};
            currAttempt[match[2]]--;

            currCost = calculateTotalCost(currAttempt)
            const isValidCost = validateAbilityCost(currCost);

            if (!isValidCost) {
                return;
            }

            validPointBuyState = currAttempt;
            document.querySelector(`#${match[2]}`).value--;
            lastCharacterCost = currCost;
        }
        else {
            alert("Couldn't find buttons");
            return;
        }

        document.querySelector("#curr-cost").textContent = "Current Cost: " + lastCharacterCost;
        document.querySelector("#curr-bonus").textContent = `Available Bonus Points: ${bonusPointsCap - creationBonusPoints}`;
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

    const formData = new FormData(form);
    const character = Object.fromEntries(formData.entries());
    const statNames = ["str", "dex", "con", "int", "wis", "cha"];

    character.stats = Object.fromEntries(
        statNames.map(function(stat){
            return [stat, Number(character[stat])];
        })
    );

    for (const stat of statNames) {
        delete character[stat];
    }

    character.skills = formData.getAll("skills");
    character.level = Number(character.level);

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
