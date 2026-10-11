import { AddMicrosoftAccount } from "../../wailsjs/go/app/App.js";
import { selectHomeLaunchOption } from "./handlerLaunchOptions.js";

let launchBar;

function closeMenus() {
    if (!launchBar) return;

    launchBar.querySelectorAll(".launch-picker").forEach((picker) => {
        picker.classList.remove("is-open");
        picker.querySelector(".launch-menu").hidden = true;
    });
}

function updateSelectedOption(option, picker) {
    const trigger = picker.querySelector(".launch-selection");
    const isInstance = picker.classList.contains("launch-instance-picker");

    trigger.querySelector(".launch-selection-title").textContent = option.dataset.title;
    trigger.querySelector(".launch-selection-detail").textContent = option.dataset.detail;
    trigger.querySelector("img").src = option.dataset.image;

    picker.querySelectorAll(".launch-menu-item").forEach((item) => {
        const selected = item === option;
        item.classList.toggle("is-selected", selected);

        const badge = item.querySelector(".launch-menu-badge");
        if (badge) badge.remove();

        if (selected && isInstance) {
            const selectedBadge = document.createElement("span");
            selectedBadge.className = "launch-menu-badge";
            selectedBadge.textContent = "Selected";
            item.append(selectedBadge);
        }
    });
}

function handleLaunchBarClick(event) {
    const trigger = event.target.closest(".launch-selection");
    if (trigger) {
        const picker = trigger.closest(".launch-picker");
        const menu = picker.querySelector(".launch-menu");
        const shouldOpen = menu.hidden;

        closeMenus();

        if (shouldOpen) {
            picker.classList.add("is-open");
            menu.hidden = false;
        }
        return;
    }

    const option = event.target.closest(".launch-menu-item");
    if (option) {
        const picker = option.closest(".launch-picker");
        const isInstance = picker.classList.contains("launch-instance-picker");

        selectHomeLaunchOption(
            isInstance ? option.dataset.id : option.dataset.uuid,
            isInstance,
        );

        updateSelectedOption(option, picker);
        closeMenus();
        return;
    }

    if (event.target.closest(".add-account")) {
        closeMenus();

        AddMicrosoftAccount().catch((error) => {
            console.error("Failed to add Microsoft account:", error);
        });
    }
}

function handleDocumentClick(event) {
    if (!launchBar?.contains(event.target)) {
        closeMenus();
    }
}

function handleDocumentKeydown(event) {
    if (event.key !== "Escape") return;

    const openTrigger = launchBar?.querySelector(".launch-picker.is-open .launch-selection");
    closeMenus();
    openTrigger?.focus();
}

export function stopHomeLaunchMenus() {
    if (!launchBar) return;

    launchBar.removeEventListener("click", handleLaunchBarClick);
    document.removeEventListener("click", handleDocumentClick);
    document.removeEventListener("keydown", handleDocumentKeydown);
    launchBar = undefined;
}

export function setupHomeLaunchMenus() {
    launchBar = document.querySelector(".home-launch-bar");
    if (!launchBar) return;

    launchBar.addEventListener("click", handleLaunchBarClick);
    document.addEventListener("click", handleDocumentClick);
    document.addEventListener("keydown", handleDocumentKeydown);
}
