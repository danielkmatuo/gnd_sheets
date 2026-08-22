const classSelect = document.querySelector("#class");
const classInfo = document.querySelector("#class-info");
const characterStats = document.querySelector("#abilities");
const form = document.querySelector("#character-form");

//make user able to only choose skills equal to the number of skills cap
function quantitySkillsCheck (quantity, data) {
    const skillIssues = data.skill_choices.choose
    if (quantity > skillIssues) {
        return false
    }

    return true
}

async function getClassInfoFromReference (selectedClass) {
    const response = await fetch(
        `/reference/classes/${selectedClass}`
    )

    if (!response.ok) {
        classInfo.textContent = "could not load class informataion from server."
        return
    }

    return response
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
    const maxCost = 27;
    var cost = 0;

    for (i = 0; i < scores.length; i++) {
        cost += scoresCost.get(scores[i]);
    }

    if (cost > maxCost) {
        return [false, cost];
    }

    return [true, cost];
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

//validate character ability scores from frontend, then send the current cost to server for further validation
characterStats.addEventListener("change", async function(){
    const characterStr = Number(document.querySelector("#str").value);
    const characterDex = Number(document.querySelector("#dex").value);
    const characterCon = Number(document.querySelector("#con").value);
    const characterInt = Number(document.querySelector("#int").value);
    const characterWis = Number(document.querySelector("#wis").value);
    const characterCha = Number(document.querySelector("#cha").value);

    const scores = [characterStr, characterDex, characterCon, characterInt, characterWis, characterCha];
    const costsMap = instantiateCostsMap();
    var [costOk, currCost] = validateStatsSelection(scores, costsMap);

    if (!costOk) {
        alert("your current score is above 27");
    }

    document.querySelector("#curr-cost").textContent = "Current cost: " + currCost;
})

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
        console.error("Could not create character");
        return;
    }

    window.location.href = "/characters";
});
