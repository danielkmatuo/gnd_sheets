const classSelect = document.querySelector("#class");
const classInfo = document.querySelector("#class-info");

//make user able to only choose skills equal to the number of skills cap
function quantitySkillsCheck (quantity, data) {
    const skillIssues = data.skill_choices.choose
    if (quantity > skillIssues) {
        return false
    }

    return true
}

//changes data from index.html dynamically
//TODO: add tools list in reference data
classSelect.addEventListener("change", async function () {
    const selectedClass = classSelect.value;

    if (selectedClass === "") {
        classInfo.innerHTML = "";
        return;
    }

    const response = await fetch(
        `/reference/classes/${selectedClass}`
    );

    if (!response.ok) {
        classInfo.textContent = "Could not load class information.";
        return;
    }

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

                alert(`You can only choose ${maxSkills} skills.`);
            }
        });
    });
});

//send form information to the go server and do the first validation of data
const form = document.querySelector("#character-form");

form.addEventListener("submit", async function(event) {
    event.preventDefault();

    const formData = new FormData(form);

    const character = Object.fromEntries(formData.entries());

    character.level = Number(character.level);

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
