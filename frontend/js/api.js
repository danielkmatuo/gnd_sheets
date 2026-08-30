//fetches the reference data from the backend
async function getClassInfoFromReference(selectedClass) {
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

    return await response.json();
}

async function getLanguagesInfoFromReference() {
    const response = await fetch(
        `/reference/languages`
    );

    if (!response.ok) {
        alert("Could not load languages information from server.");
        const errorMessage = await response.text();
        console.error("Go server error: ", errorMessage)
        return;
    }

    return await response.json();
}

async function getSkillsInfoFromReference() {
    const response = await fetch(
        `/reference/skills`
    );

    if (!response.ok) {
        alert("Could not load skills information from server.");
        const errorMessage = await response.text();
        console.error("Go server error: ", errorMessage)
        return;
    }

    return await response.json();
}

async function sendCharacterData(data) {
    const response = await fetch("/characters/create", {
        method: "POST",

        headers: {
            "Content-Type": "application/json"
        },

        body: JSON.stringify(data)
    });

    if (!response.ok) {
        const errorMessage = await response.text();
        console.error("Go server error:", errorMessage);
        return;
    }
}

export {
    getClassInfoFromReference,
    getRaceInfoFromReference,
    getLanguagesInfoFromReference,
    getSkillsInfoFromReference,
    sendCharacterData
};
