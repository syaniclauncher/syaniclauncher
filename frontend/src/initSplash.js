// checks if the backend is ready
// hides overlay after that
import { IsBackendReady } from "../wailsjs/go/main/App.js";

async function checkBackend() {
    try {
        const response = await IsBackendReady();
      if (response) {
            // document.getElementById("overlay").style.display = "none";
            // fade animation
            const overlay = document.getElementById("overlay");
            overlay.classList.add("fade-out");
            overlay.addEventListener("animationend", () => { overlay.style.display = "none"; }, { once: true });
        } else {
            setTimeout(checkBackend, 1000); //loop
        }
    } catch (error) {
        console.error(error);
        setTimeout(checkBackend, 1000); // 1s
    }
}

checkBackend();
