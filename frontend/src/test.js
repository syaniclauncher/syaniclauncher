// test: Check connection to backend
import { CheckConnection } from '@wails/go/app/App.js';

async function verifyBackendConnection() {
    try {
        const response = await CheckConnection();
        console.log("Backend response:", response.trim());
    } catch (error) {
        console.error("Failed to connect to backend:", error);
    }
    await AddMicrosoftAccount();
}

verifyBackendConnection();

/*import { AddMicrosoftAccount } from '@wails/go/app/App.js';
(async () => await AddMicrosoftAccount())();*/
