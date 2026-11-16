<script setup>
import { computed, onMounted, onUnmounted, ref } from "vue";
import {
    AccountAddOffline,
    AccountRemove,
    AccountSetSelected,
    AccountsList,
    AccountsRefresh,
    StartMicrosoftBrowserLogin,
    StartMicrosoftDeviceLogin,
} from "@wails/go/app/App";
import { BrowserOpenURL, EventsOn } from "@wails/runtime/runtime";

const props = defineProps({ open: Boolean });
const emit = defineEmits(["toggle-picker", "close-picker"]);
const accounts = ref([]);
const selected = ref(null);
const selectedAccount = computed(() => accounts.value.find(account => account.uuid === selected.value));
const promptOpen = ref(false);
const promptStep = ref("choose");
const offlineUsername = ref("");
const deviceLogin = ref(null);
const authStatus = ref("");
const errorMessage = ref("");
const busy = ref(false);
let stopAuthListener;

async function loadAccounts() {
    const result = await AccountsList();
    accounts.value = result.accounts || [];
    selected.value = result.selected || null;
}

async function preloadAccounts() {
    try {
        await loadAccounts();
        const result = await AccountsRefresh();
        accounts.value = result.accounts || [];
        selected.value = result.selected || null;
    } catch (error) {
        errorMessage.value = error.message || String(error);
        try {
            await loadAccounts();
        } catch (loadError) {
            errorMessage.value = loadError.message || String(loadError);
        }
    }
}

function showAddPrompt() {
    emit("close-picker");
    errorMessage.value = "";
    promptStep.value = "choose";
    offlineUsername.value = "";
    deviceLogin.value = null;
    authStatus.value = "";
    promptOpen.value = true;
}

function closePrompt() {
    promptOpen.value = false;
    busy.value = false;
}

async function selectAccount(uuid) {
    try {
        const result = await AccountSetSelected(uuid);
        accounts.value = result.accounts || [];
        selected.value = result.selected || null;
        emit("close-picker");
    } catch (error) {
        errorMessage.value = error.message || String(error);
    }
}

async function removeAccount(account) {
    if (!window.confirm(`Remove ${account.username}?`)) return;
    try {
        const result = await AccountRemove(account.uuid);
        accounts.value = result.accounts || [];
        selected.value = result.selected || null;
    } catch (error) {
        errorMessage.value = error.message || String(error);
    }
}

async function addOfflineAccount() {
    busy.value = true;
    errorMessage.value = "";
    try {
        const result = await AccountAddOffline(offlineUsername.value);
        accounts.value = result.accounts || [];
        selected.value = result.selected || null;
        closePrompt();
    } catch (error) {
        errorMessage.value = error.message || String(error);
    } finally {
        busy.value = false;
    }
}

async function startBrowserLogin() {
    promptStep.value = "browser";
    authStatus.value = "Opening Microsoft sign-in in your browser…";
    errorMessage.value = "";
    busy.value = true;
    try {
        await StartMicrosoftBrowserLogin();
        authStatus.value = "Finish signing in in your browser.";
    } catch (error) {
        errorMessage.value = error.message || String(error);
        busy.value = false;
    }
}

async function startDeviceLogin() {
    promptStep.value = "device";
    authStatus.value = "Requesting a sign-in code…";
    errorMessage.value = "";
    busy.value = true;
    try {
        deviceLogin.value = await StartMicrosoftDeviceLogin();
        authStatus.value = "Enter this code in the Microsoft sign-in page.";
    } catch (error) {
        errorMessage.value = error.message || String(error);
        busy.value = false;
    }
}

async function copyDeviceCode() {
    try {
        await navigator.clipboard.writeText(deviceLogin.value.user_code);
    } catch (error) {
        errorMessage.value = error.message || "Unable to copy the code.";
    }
}

onMounted(() => {
    stopAuthListener = EventsOn("accounts:auth", async (status, message) => {
        busy.value = false;
        if (status === "success") {
            try {
                await loadAccounts();
                closePrompt();
            } catch (error) {
                errorMessage.value = error.message || String(error);
            }
        } else {
            errorMessage.value = message || "Microsoft sign-in failed.";
        }
    });
    preloadAccounts();
});

onUnmounted(() => stopAuthListener?.());
</script>

<template>
    <div class="launch-picker launch-account-picker" :class="{ 'is-open': props.open }" @click.stop>
        <div v-if="props.open" class="launch-menu launch-account-menu">
            <div v-if="!accounts.length" class="account-empty">No accounts yet.</div>
            <div v-for="account in accounts" :key="account.uuid" class="account-menu-row">
                <button class="account-menu-select" type="button" @click="selectAccount(account.uuid)">
                    <span class="account-menu-avatar">{{ account.username.slice(0, 1).toUpperCase() }}</span>
                    <span class="account-menu-copy">
                        <span class="launch-selection-title">{{ account.username }}</span>
                        <span class="launch-selection-detail">
                            {{ account.type === "msa" ? (account.expired ? "Microsoft · Sign in again" : "Microsoft") : "Offline" }}
                        </span>
                    </span>
                    <span v-if="account.uuid === selected" class="account-selected-mark">✓</span>
                </button>
                <button class="account-remove" type="button" title="Remove account"
                    @click.stop="removeAccount(account)">×</button>
            </div>
            <div v-if="errorMessage" class="account-menu-error">{{ errorMessage }}</div>
            <button class="launch-menu-action add-account" type="button" @click="showAddPrompt">
                <span class="launch-menu-action-icon"><i class="icon icon-add" style="--icon-size: 16px;"></i></span>
                <span>Add Account</span>
            </button>
        </div>
        <button class="launch-selection launch-account" type="button" @click="emit('toggle-picker')">
            <span class="launch-avatar">{{ selectedAccount ? selectedAccount.username.slice(0, 1).toUpperCase() : "?" }}</span>
            <span class="launch-selection-copy">
                <span class="launch-selection-title">{{ selectedAccount?.username || "No account" }}</span>
                <span class="launch-selection-detail">
                    {{ selectedAccount ? (selectedAccount.type === "msa" ? "Microsoft account" : "Offline account") : "Sign in" }}
                </span>
            </span>
            <i class="icon icon-chevron" style="--icon-size: 16px; --icon-color: var(--launch-chevron-color);"></i>
        </button>

        <div v-if="promptOpen" class="account-overlay" @click.self="closePrompt" @keydown.esc.stop.prevent="closePrompt">
            <div class="account-prompt" role="dialog" aria-modal="true" aria-labelledby="account-prompt-title">
                <header class="account-prompt-header">
                    <h2 id="account-prompt-title">
                        {{ promptStep === "choose" ? "Add an account" : promptStep === "offline" ? "Offline account" : "Microsoft account" }}
                    </h2>
                    <button class="account-prompt-close" type="button" aria-label="Close" @click="closePrompt">×</button>
                </header>

                <template v-if="promptStep === 'choose'">
                    <button class="account-prompt-option" type="button" @click="promptStep = 'offline'">Offline username</button>
                    <button class="account-prompt-option" type="button" @click="startBrowserLogin">Microsoft — browser sign-in</button>
                    <button class="account-prompt-option" type="button" @click="startDeviceLogin">Microsoft — device code</button>
                </template>

                <form v-else-if="promptStep === 'offline'" class="offline-account-form" @submit.prevent="addOfflineAccount">
                    <label for="offline-username">Username</label>
                    <input id="offline-username" v-model="offlineUsername" maxlength="16" placeholder="Player" autofocus>
                    <button class="account-prompt-submit" type="submit" :disabled="busy">
                        {{ busy ? "Adding…" : "Add account" }}
                    </button>
                </form>

                <div v-else class="account-auth-status">
                    <p>{{ authStatus }}</p>
                    <template v-if="deviceLogin">
                        <button class="account-prompt-option device-link" type="button"
                            @click="BrowserOpenURL(deviceLogin.verification_uri)">
                            Open {{ deviceLogin.verification_uri }}
                        </button>
                        <div class="device-code-row">
                            <code>{{ deviceLogin.user_code }}</code>
                            <button type="button" @click="copyDeviceCode">Copy code</button>
                        </div>
                    </template>
                </div>

                <p v-if="errorMessage" class="account-prompt-error" role="alert">{{ errorMessage }}</p>
                <button v-if="promptStep !== 'choose'" class="account-prompt-back" type="button"
                    :disabled="busy" @click="promptStep = 'choose'; errorMessage = ''">Back</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.launch-picker { position: relative; flex: 0 1 min(250px, 30%); min-width: 0; }
.launch-selection { box-sizing: border-box; display: flex; align-items: center; gap: 0.5rem; width: min(250px, 30%); min-width: 0; height: 54px; padding: 0.5rem 0.75rem; border: 1px solid transparent; border-radius: 7px; background: var(--launch-card-bg-color); color: var(--text-color); text-align: left; cursor: pointer; }
.launch-picker .launch-selection { width: 100%; }
.launch-selection:hover { background: var(--launch-card-hover-bg-color); }
.launch-avatar { display: inline-flex; flex: 0 0 30px; align-items: center; justify-content: center; width: 30px; height: 30px; border-radius: 6px; background: var(--surface-5); color: var(--text-color); font-weight: 700; }
.launch-menu { position: absolute; bottom: calc(100% + 0.45rem); left: 0; z-index: 5; box-sizing: border-box; display: flex; flex-direction: column; width: 100%; max-height: min(320px, calc(100vh - 110px)); overflow-x: hidden; overflow-y: auto; border-radius: 8px; background: var(--surface-4); }
.launch-menu-action { box-sizing: border-box; display: flex; flex: 0 0 auto; align-items: center; gap: 0.65rem; min-width: 0; min-height: 36px; padding: 0.15rem 0.75rem; border: 0; border-bottom: 1px solid var(--gray-color); background: var(--surface-4); color: var(--text-color); text-align: left; cursor: pointer; font-size: 0.8rem; font-weight: 600; }
.launch-menu-action:hover { background: var(--surface-hover-color); }
.launch-menu-action-icon { display: inline-flex; flex: 0 0 16px; align-items: center; justify-content: center; width: 16px; color: var(--text-muted-color); }
.launch-selection-copy { display: flex; flex: 1; flex-direction: column; gap: 0.15rem; min-width: 0; }
.launch-selection-title, .launch-selection-detail { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.launch-selection-title { font-size: 0.875rem; font-weight: 600; }
.launch-selection-detail { color: var(--text-muted-color); font-size: 0.75rem; opacity: 0.7; }
.account-menu-row { display: flex; align-items: center; min-height: 56px; border-bottom: 1px solid var(--gray-color); }
.account-menu-select { display: flex; flex: 1; align-items: center; gap: 0.65rem; min-width: 0; padding: 0.45rem 0.65rem; border: 0; background: transparent; color: var(--text-color); text-align: left; cursor: pointer; }
.account-menu-select:hover, .account-remove:hover { background: var(--surface-hover-color); }
.account-menu-avatar { display: inline-flex; flex: 0 0 30px; align-items: center; justify-content: center; width: 30px; height: 30px; border-radius: 6px; background: var(--surface-5); font-weight: 700; }
.account-menu-copy { display: flex; flex: 1; flex-direction: column; gap: 0.15rem; min-width: 0; }
.account-selected-mark { color: #77df78; }
.account-remove { width: 32px; height: 32px; margin-right: 0.45rem; border: 0; border-radius: 6px; background: transparent; color: var(--text-muted-color); font-size: 1.2rem; cursor: pointer; }
.account-empty, .account-menu-error { padding: 0.75rem; color: var(--text-muted-color); font-size: 0.8rem; }
.account-menu-error, .account-prompt-error { color: #ff8585; overflow-wrap: anywhere; }
.account-overlay { position: fixed; inset: 0; z-index: 20; display: flex; align-items: center; justify-content: center; padding: 1rem; background: rgba(0, 0, 0, 0.72); }
.account-prompt { display: flex; flex-direction: column; gap: 0.8rem; width: min(420px, 100%); padding: 1.2rem; border: 1px solid var(--gray-color); border-radius: 10px; background: var(--surface-4); box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5); }
.account-prompt-header { display: flex; align-items: center; justify-content: space-between; }
.account-prompt-header h2 { margin: 0; font-size: 1.05rem; }
.account-prompt-close { border: 0; background: transparent; color: var(--text-muted-color); font-size: 1.5rem; cursor: pointer; }
.account-prompt-option, .account-prompt-submit, .device-code-row button { min-height: 42px; padding: 0.65rem 0.8rem; border: 1px solid var(--gray-color); border-radius: 7px; background: var(--surface-5); color: var(--text-color); font: inherit; text-align: left; cursor: pointer; }
.account-prompt-option:hover, .account-prompt-submit:hover, .device-code-row button:hover { background: var(--surface-hover-color); }
.offline-account-form { display: flex; flex-direction: column; gap: 0.6rem; }
.offline-account-form input { min-height: 40px; padding: 0 0.7rem; border: 1px solid var(--gray-color); border-radius: 7px; background: var(--surface-2); color: var(--text-color); font: inherit; }
.account-prompt-submit { background: var(--accent-color); text-align: center; }
.account-prompt-submit:disabled, .account-prompt-back:disabled { opacity: 0.6; cursor: default; }
.account-auth-status p { color: var(--text-muted-color); line-height: 1.5; }
.device-link { width: 100%; overflow-wrap: anywhere; }
.device-code-row { display: flex; align-items: center; justify-content: space-between; gap: 0.6rem; margin-top: 0.65rem; }
.device-code-row code { font-size: 1.25rem; letter-spacing: 0.12em; }
.account-prompt-back { align-self: flex-start; padding: 0.4rem 0.6rem; border: 0; background: transparent; color: var(--text-muted-color); cursor: pointer; }
</style>
