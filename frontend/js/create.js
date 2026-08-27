const classSelect = document.querySelector("#class");
const classInfo = document.querySelector("#class-info");
const characterStats = document.querySelector("#abilities");
const form = document.querySelector("#character-form");
const levelSelect = document.querySelector("#level");
let creationBonusPoints = 0;
let bonusPointsCap = 3;

const bonusButtons = document.querySelectorAll("#abilities button");

let allocatedBonus = {
    "str": 0,
    "dex": 0,
    "con": 0,
    "int": 0,
    "wis": 0,
    "cha": 0
}

let lastValidScores = {
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

function validateStatsSelection(scores, scoresCost) {
    const currLevel = Number(levelSelect.value);
    if (currLevel === 1) {
        const maxCost = 27;
        var cost = 0;
        var isWithinRange = true;

        for (i = 0; i < scores.length; i++) {
            if (!(scores[i] >= 8 && scores[i] <= 15)) {
                isWithinRange = false; 
                return [true, cost, isWithinRange];
            }

            cost += scoresCost.get(scores[i]);
        }

        if (cost > maxCost) {
            return [false, cost, isWithinRange];
        }

        return [true, cost, isWithinRange];
    }
    else {
        const isValid = validateStatsSelectionAboveLvlOne(scores);

        if (!isValid) {
            return [false, 0, isValid];
        }

        return [false, 0, isValid];
    }
}

function validateCreationBonusPoints(bonus) {
    if (creationBonusPoints + bonus > bonusPointsCap) {
        return false
    }

    return true
}

//TODO: Find a way to make this code ensure the correct structure for the initial bonus points allocation (one stat +2 and another +1)
function validateAllocationBonusPoints() {
    let allocatedBonusPoints = creationBonusPoints;
    const statsVector = [
        allocatedBonus.str,
        allocatedBonus.dex,
        allocatedBonus.con,
        allocatedBonus.int,
        allocatedBonus.wis,
        allocatedBonus.cha
    ] 

    for (let i = 0; i < statsVector.length; i++) {
        if (statsVector[i] % 2 != 0) {
            statsVector[i]--;
        } 
        else {
            statsVector[i] -= 2;
        }
    }
}

function validateStatsSelectionAboveLvlOne(scores) {
    for (let i = 0; i < scores.length; i++) {
        if (!(scores[i] >= 1 && scores[i] <= 30)) {
            return false;
        }
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
characterStats.addEventListener("change", async function(event){
    if (event.target.value === "") {
        event.target.value = lastValidScores[event.target.id];
        return;
    }

    const characterStr = Number(document.querySelector("#str").value);
    const characterDex = Number(document.querySelector("#dex").value);
    const characterCon = Number(document.querySelector("#con").value);
    const characterInt = Number(document.querySelector("#int").value);
    const characterWis = Number(document.querySelector("#wis").value);
    const characterCha = Number(document.querySelector("#cha").value);
    const characterLvl = Number(form.querySelector("#level").value);

    const scores = [characterStr, characterDex, characterCon, characterInt, characterWis, characterCha];
    const costsMap = instantiateCostsMap();
    
    const [costOk, cost, isWithinRange] = validateStatsSelection(scores, costsMap)

    if (costOk && !isWithinRange) {
        alert(`Your current score (${event.target.value}) is out of range (range must be within 8 to 15).`);
        return;
    } 
    else if ((!costOk && cost != 0) && isWithinRange) {
        alert(`Your current cost (${cost}) is above 27.`);
        return;
    }
    else if ((!costOk && cost === 0) && !isWithinRange) {
        alert(`Your current score(${event.target.value}) is out of range (range must be within 1 to 30).`)
        return;
    }
    else if ((!costOk && cost === 0) && isWithinRange) {
        lastValidScores[event.target.id] = Number(event.target.value);
        document.querySelector("#curr-cost").textContent = "Current cost: not applicable for characters above level 1.";
        return;
    }
    else if (costOk && isWithinRange) {
        lastValidScores[event.target.id] = Number(event.target.value);
        document.querySelector("#curr-cost").textContent = "Current cost: " + cost;
        return;
    }
    else {
        alert("Something unexpected happened...")
        return;
    }
});

//TODO: find a way to make this event work like: user click this specific button, then add the respective bonus points
bonusButtons.forEach(function(button) {
    button.addEventListener("click", function(event) {
        const regex = /^(one|two)-bonus-(str|dex|con|int|wis|cha)$/;
        const match = event.target.id.match(regex);

        if (match[1] === "one") {
            const creationBonusPointsOk = validateCreationBonusPoints(1);
            if (!creationBonusPointsOk) {
                alert("Already spent all bonus points for this character");
                return;
            }

            creationBonusPoints++;
            switch (match[2]) {
                case "str":
                    document.querySelector("#str").value++;
                    break;
                case "dex":
                    document.querySelector("#dex").value++;
                    break;
                case "con":
                    document.querySelector("#con").value++;
                    break;
                case "int":
                    document.querySelector("#int").value++;
                    break;
                case "wis":
                    document.querySelector("#wis").value++;
                    break;
                case "cha":
                    document.querySelector("#cha").value++;
                    break;
            } 
            allocatedBonus[match[2]]++;
        }
        else if (match[1] == "two") {
            const creationBonusPointsOk = validateCreationBonusPoints(2);
            if (!creationBonusPointsOk) {
                alert("Already spent all bonus points for this character");
                return;
            }

            creationBonusPoints += 2;
            switch (match[2]) {
                case "str":
                    document.querySelector("#str").value = Number(document.querySelector("#str").value) + 2;
                    break;
                case "dex":
                    document.querySelector("#dex").value = Number(document.querySelector("#dex").value) + 2;
                    break;
                case "con":
                    document.querySelector("#con").value = Number(document.querySelector("#con").value) + 2;
                    break;
                case "int":
                    document.querySelector("#int").value = Number(document.querySelector("#int").value) + 2;
                    break;
                case "wis":
                    document.querySelector("#wis").value = Number(document.querySelector("#wis").value) + 2;
                    break;
                case "cha":
                    document.querySelector("#cha").value = Number(document.querySelector("#cha").value) + 2;
                    break;
            } 
            allocatedBonus[match[2]] += 2;
        }
        else {
            alert("Couldnt find button");
            return;
        }

        document.querySelector("#curr-bonus").textContent = `Available Bonus Points: ${bonusPointsCap - creationBonusPoints}`;
    });
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
