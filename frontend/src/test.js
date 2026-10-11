// test: Check connection to backend
import { CheckConnection } from '../wailsjs/go/app/App.js';

async function verifyBackendConnection() {
    try {
        const response = await CheckConnection();
        console.log("Backend response:", response.trim());
    } catch (error) {
        console.error("Failed to connect to backend:", error);
    }
}

verifyBackendConnection();

/*import { AddMicrosoftAccount } from '../wailsjs/go/app/App.js';
await AddMicrosoftAccount();*/
