const emptySelections = {
    account: {
        title: "No account",
        detail: "Sign in",
        image: "assets/portraits/question-mark.png",
    },
    instance: {
        title: "No instance selected",
        detail: "Pick in Instances",
        image: "assets/instance_icons/grass_block.png",
    },
};

// demo: use demos until implemented
const demoAccounts = [
    {uuid: "demo-account-steve", title: "Steve", detail: "Microsoft", image: "assets/portraits/steve.png",},
    {uuid: "demo-account-alex", title: "Alex", detail: "Offline", image: "assets/portraits/offline-steve.png",},
];

const demoInstances = [
    {id: "demo-instance-survival", title: "Survival World", detail: "1.20.4", image: "assets/instance_icons/grass_block.png",},
    {id: "demo-instance-build", title: "Build Test", detail: "Creative", image: "assets/instance_icons/crafting_table.png",},
];

let accounts = demoAccounts;
let instances = demoInstances;
let selectedAccountUUID = accounts[0]?.uuid ?? null; //auto-select first account if available for demo
let selectedInstanceID = instances[0]?.id ?? null;

//let accounts = [];
//let instances = [];
//let selectedAccountUUID = null;
//let selectedInstanceID = null;



function getOptionIdentifier(option, isInstance) {
    return isInstance ? option.id : option.uuid;
}

function getEmptySelection(isInstance) {
    return isInstance ? emptySelections.instance : emptySelections.account;
}

function createSelectedBadge() {
    const badge = document.createElement("span");
    badge.className = "launch-menu-badge";
    badge.textContent = "Selected";
    return badge;
}

function createLaunchOption(option, imageClass, isSelected, isInstance) {
    const item = document.createElement("button");
    item.className = `launch-menu-item${isSelected ? " is-selected" : ""}`;
    item.type = "button";

    if (isInstance) {
        item.dataset.id = option.id;
    } else {
        item.dataset.uuid = option.uuid;
    }

    item.dataset.title = option.title;
    item.dataset.detail = option.detail;
    item.dataset.image = option.image;

    const image = document.createElement("img");
    image.className = imageClass;
    image.src = option.image;
    image.alt = "";
    item.append(image);

    const copy = document.createElement("span");
    copy.className = "launch-selection-copy";

    const title = document.createElement("span");
    title.className = "launch-selection-title";
    title.textContent = option.title;

    const detail = document.createElement("span");
    detail.className = "launch-selection-detail";
    detail.textContent = option.detail;

    copy.append(title, detail);
    item.append(copy);

    if (isInstance && isSelected) {
        item.append(createSelectedBadge());
    }

    return item;
}

function populatePicker(picker, options, imageClass, isInstance) {
    if (!picker) return;

    const menu = picker.querySelector(".launch-menu");
    const trigger = picker.querySelector(".launch-selection");
    const action = menu.querySelector(".launch-menu-action");
    const selectedIdentifier = isInstance ? selectedInstanceID : selectedAccountUUID;
    const selectedOption = options.find(
        (option) => getOptionIdentifier(option, isInstance) === selectedIdentifier,
    );

    menu.querySelectorAll(".launch-menu-item").forEach((item) => item.remove());

    const displayOption = selectedOption ?? getEmptySelection(isInstance);
    trigger.querySelector(".launch-selection-title").textContent = displayOption.title;
    trigger.querySelector(".launch-selection-detail").textContent = displayOption.detail;
    trigger.querySelector("img").src = displayOption.image;

    options.forEach((option) => {
        const isSelected = getOptionIdentifier(option, isInstance) === selectedIdentifier;
        menu.insertBefore(createLaunchOption(option, imageClass, isSelected, isInstance), action);
    });
}

export function selectHomeLaunchOption(identifier, isInstance) {
    const optionList = isInstance ? instances : accounts;
    const optionKey = isInstance ? "id" : "uuid";

    if (!optionList.some((option) => option[optionKey] === identifier)) return;

    if (isInstance) {
        selectedInstanceID = identifier;
    } else {
        selectedAccountUUID = identifier;
    }
}

export function populateHomeLaunchOptions() {
    selectedAccountUUID = accounts.some((account) => account.uuid === selectedAccountUUID)
        ? selectedAccountUUID
        : null;
    selectedInstanceID = instances.some((instance) => instance.id === selectedInstanceID)
        ? selectedInstanceID
        : null;

    const launchBar = document.querySelector(".home-launch-bar");
    if (!launchBar) return;

    populatePicker(launchBar.querySelector(".launch-account-picker"), accounts, "launch-avatar", false);
    populatePicker(launchBar.querySelector(".launch-instance-picker"), instances, "launch-instance-icon", true);
}
