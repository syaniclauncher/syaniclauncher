// test: Check connection to backend
import { CheckConnection } from '../wailsjs/go/main/App.js';

async function verifyBackendConnection() {
    try {
        const response = await CheckConnection();
        console.log("Backend response:", response.trim());
    } catch (error) {
        console.error("Failed to connect to backend:", error);
    }
}

verifyBackendConnection();
