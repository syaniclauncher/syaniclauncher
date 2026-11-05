// checks if the backend is ready
// hides overlay after that
import { IsBackendReady } from "../wailsjs/go/app/App.js";

async function checkBackend() {
    try {
        if (await IsBackendReady()) {
            const overlay = document.getElementById("overlay");
            overlay.classList.add("fade-out");
            overlay.addEventListener("animationend", () => overlay.style.display = "none", { once: true });
            return;
        }
    } catch (error) {
        console.error(error);
    }

    setTimeout(checkBackend, 1000);
}

checkBackend();
