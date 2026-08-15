//changes data from index.html dynamically
const classSelect = document.querySelector("#class");
const classInfo = document.querySelector("#class-info");

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

    classInfo.innerHTML = `
        <h2>${data.name}</h2>
        <p>Hit Die: ${data.hit_die}</p>
        <p>Saving Throws: ${data.proficiencies.saving_throws.join(", ")}</p>
        <p>Choose ${data.skill_choices.choose} skills:</p>
        <p>${data.skill_choices.from.join(", ")}</p>
    `;
});

//send form information to the go server
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
