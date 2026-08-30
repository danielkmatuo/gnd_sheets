//make user able to only choose skills equal to the number of skills cap
function quantitySkillsCheck (quantity, data) {
    const skillIssues = data.skill_choices.choose;
    if (quantity > skillIssues) {
        return false;
    }

    return true;
}

//calculate proficiency bonus based on level selected by the user
function calculateProficiencyBonus(level) {
    const progressionStep = 4;
    let baselineBonus = 2;
    let progressionAnchor = 5;

    if (level < 1 || level > 20) {
        alert("Invalid level. Must be in range 1 to 20");
        return -1;
    }

    for (let i = 5; i <= level; i++) {
        if (i === progressionAnchor) {
            baselineBonus++;
        }
        else if (i - progressionAnchor === progressionStep) {
            progressionAnchor = i;
            baselineBonus++;
        }
    }
    
    return baselineBonus;
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

//TODO: change this func whenever I start to work on feats
function calculateBonusPointsCap() {
    return 3;
}

function validateCreationBonusPoints(allocationAttempt) {
    let bonusSum = 0;

    for (const value of Object.values(allocationAttempt)) {
        bonusSum += value;
    }

    if (bonusSum > calculateBonusPointsCap()) {
        return false;
    }

    return true;
}

function tryAllocationBonusPoints(bonus, currStat, pastState) {
    const bonusKeys = [
        "str",
        "dex",
        "con",
        "int",
        "wis",
        "cha"
    ];
    const attemptedAllocation = {...pastState};
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

    pastState[currStat] = attemptedAllocation[currStat];

    return true;
}

function tryPointBuyChange(stat, newValue, pastState) {
    const attempt = {...pastState};

    attempt[stat] = newValue;

    const cost = calculateTotalCost(attempt);

    if (!validateAbilityCost(cost)) {
        return false;
    }

    pastState[stat] = attempt[stat];

    return true;
}

function checkBothStates(bonus, stat, newValue, pointBuyState, allocationState) {
    const pointBuyAttemptOk = tryPointBuyChange(stat, newValue, pointBuyState);

    if (pointBuyAttemptOk) {
        const allocationOk = tryAllocationBonusPoints(bonus, stat, allocationState);
        if (!allocationOk) {
            return false;
        }

        return true;
    }

    return false;
}

//manipulate the skills data received from the backend
function createSkillsObj(data) {
    let skillsObj = {};

    for (const key of Object.keys(data)) {
        let innerObj = {};
        for (const skill of Object.values(data[key])) {
            innerObj[skill] = 0;
        }
        skillsObj[key] = innerObj;
    }

    return skillsObj;
}

function createBaselineSkillsObj(data, modifiers) {
    let baselineSkillsObj = createSkillsObj(data);
     
    for (const key of Object.keys(modifiers)) {
        let innerObj = baselineSkillsObj[key];
        for (const skill of Object.keys(innerObj)) {
            baselineSkillsObj[key][skill] = modifiers[key]; 
        }
    }
    
    return baselineSkillsObj;
}

//calculate ability scores modifiers
function calculateAbilityModifiers(pointBuyState) {
    let modifiers = {};

    for (const key of Object.keys(pointBuyState)) {
        const stat = pointBuyState[key];
        let bonus = 0; 
        const baseline = 10;

        if (stat % 2 == 0) {
            if (stat >= baseline) {
                bonus = (stat - baseline) / 2;
            }
            else {
                bonus = ((baseline - stat) * (-1)) / 2; 
            }
        }
        else if (stat < baseline) {
            bonus = ((baseline - stat + 1) * (-1)) / 2;
        }
        else {
            bonus = (stat - baseline - 1) / 2;
        }

        modifiers[key] = bonus;
    }

    return modifiers;
}

//helpers
function getModifierString(value) {
    if (value >= 0) {
        return "+" + value;
    }

    return value;
}

function capitalize(name) {
    const arr = name.split(" ");

    if (arr.length < 2) {
        const firstLetter = name.charAt(0);
        const upperFirstLetter = firstLetter.toUpperCase();
        return name.replace(name.charAt(0), upperFirstLetter);
    }
    else {
        let newName = "";

        for (let i = 0; i < arr.length; i++) {
            if (arr[i] != "of") {
                const firstLetter = arr[i].charAt(0);
                const upperFirstLetter = firstLetter.toUpperCase();
                if (i < arr.length - 1) { 
                    newName += arr[i].replace(arr[i].charAt(0), upperFirstLetter) + " ";
                }
                else {
                    newName += arr[i].replace(arr[i].charAt(0), upperFirstLetter);
                }
            }
            else {
                newName += arr[i] + " ";
            }
        }
        return newName;
    }
}

export {
    quantitySkillsCheck,
    calculateProficiencyBonus,
    instantiateCostsMap,
    calculateTotalCost,
    calculateBonusPointsCap,
    validateAbilityCost,
    validateCreationBonusPoints,
    tryAllocationBonusPoints,
    tryPointBuyChange,
    checkBothStates,
    createBaselineSkillsObj,
    calculateAbilityModifiers,
    getModifierString,
    capitalize
};
