<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from "vue";
import {
    AccountAddOffline,
    AccountRemove,
    AccountSetSelected,
    AccountsList,
    AccountsRefresh,
    CancelMicrosoftLogin,
    StartMicrosoftBrowserLogin,
    StartMicrosoftDeviceLogin,
} from "@wails/go/app/App";
import { BrowserOpenURL, EventsOn } from "@wails/runtime/runtime";

const props = defineProps({ open: Boolean });
const emit = defineEmits(["toggle-picker", "close-picker"]);
const accounts = ref([]);
const selected = ref(null);
const menuError = ref("");
const pendingRemovalUUID = ref(null);
const removingAccount = ref(false);
const selectedAccount = computed(() => accounts.value.find(account => account.uuid === selected.value));
const promptOpen = ref(false);
const promptStep = ref("choose");
const offlineUsername = ref("");
const deviceLogin = ref(null);
const browserLoginURL = ref("");
const authStatus = ref("");
const errorMessage = ref("");
const busy = ref(false);
const qrContainer = ref(null);
const codeCopied = ref(false);
let authAttempt = 0;
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
        menuError.value = error.message || String(error);
        try {
            await loadAccounts();
        } catch (loadError) {
            menuError.value = loadError.message || String(loadError);
        }
    }
}

function showAddPrompt() {
    emit("close-picker");
    pendingRemovalUUID.value = null;
    errorMessage.value = "";
    promptStep.value = "choose";
    offlineUsername.value = "";
    deviceLogin.value = null;
    browserLoginURL.value = "";
    authStatus.value = "";
    promptOpen.value = true;
}

function closePrompt() {
    authAttempt++;
    CancelMicrosoftLogin();
    promptOpen.value = false;
    busy.value = false;
    errorMessage.value = "";
}

async function selectAccount(uuid) {
    try {
        const result = await AccountSetSelected(uuid);
        accounts.value = result.accounts || [];
        selected.value = result.selected || null;
        pendingRemovalUUID.value = null;
        emit("close-picker");
    } catch (error) {
        menuError.value = error.message || String(error);
    }
}

function toggleRemoveConfirmation(uuid) {
    pendingRemovalUUID.value = pendingRemovalUUID.value === uuid ? null : uuid;
}

async function removeAccount(account) {
    removingAccount.value = true;
    try {
        const result = await AccountRemove(account.uuid);
        accounts.value = result.accounts || [];
        selected.value = result.selected || null;
        pendingRemovalUUID.value = null;
    } catch (error) {
        menuError.value = error.message || String(error);
    } finally {
        removingAccount.value = false;
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
    const attempt = ++authAttempt;
    promptStep.value = "browser";
    deviceLogin.value = null;
    browserLoginURL.value = "";
    authStatus.value = "Opening Microsoft sign-in in your browser…";
    errorMessage.value = "";
    busy.value = true;
    try {
        browserLoginURL.value = await StartMicrosoftBrowserLogin();
        if (attempt !== authAttempt || !promptOpen.value) return;
        authStatus.value = "Finish signing in in your browser.";
    } catch (error) {
        if (attempt !== authAttempt || !promptOpen.value) return;
        errorMessage.value = error.message || String(error);
        busy.value = false;
    }
}

async function startDeviceLogin() {
    const attempt = ++authAttempt;
    promptStep.value = "device";
    authStatus.value = "Requesting a sign-in code…";
    errorMessage.value = "";
    deviceLogin.value = null;
    browserLoginURL.value = "";
    busy.value = true;
    try {
        const login = await StartMicrosoftDeviceLogin();
        if (attempt !== authAttempt || !promptOpen.value) return;
        deviceLogin.value = login;
        authStatus.value = "Waiting for you to finish signing in…";
    } catch (error) {
        if (attempt !== authAttempt || !promptOpen.value) return;
        errorMessage.value = error.message || String(error);
        busy.value = false;
    }
}

async function copyDeviceCode() {
    try {
        await navigator.clipboard.writeText(deviceLogin.value.user_code);
        codeCopied.value = true;
        setTimeout(() => { codeCopied.value = false; }, 1200);
    } catch (error) {
        errorMessage.value = error.message || "Unable to copy the code.";
    }
}

watch(deviceLogin, async login => {
    codeCopied.value = false;
    if (!login) return;
    await nextTick();
    if (!qrContainer.value || !window.QRCode) {
        errorMessage.value = "The sign-in QR code could not be loaded.";
        return;
    }
    qrContainer.value.replaceChildren();
    new window.QRCode(qrContainer.value, {
        text: login.verification_uri,
        width: 132,
        height: 132,
        correctLevel: window.QRCode.CorrectLevel.M,
    });
});

onMounted(() => {
    stopAuthListener = EventsOn("accounts:auth", async (status, message) => {
        busy.value = false;
        if (status === "success") {
            try {
                await loadAccounts();
                if (promptOpen.value) closePrompt();
            } catch (error) {
                if (promptOpen.value) errorMessage.value = error.message || String(error);
            }
        } else if (promptOpen.value) {
            errorMessage.value = message || "Microsoft sign-in failed.";
        }
    });
    preloadAccounts();
});

onUnmounted(() => {
    stopAuthListener?.();
    CancelMicrosoftLogin();
});
</script>

<template>
    <div class="launch-picker launch-account-picker" :class="{ 'is-open': props.open }" @click.stop>
        <div v-if="props.open" class="launch-menu launch-account-menu">
            <div v-if="!accounts.length" class="account-empty">No accounts yet.</div>
            <div v-for="account in accounts" :key="account.uuid" class="account-menu-row"
                :class="{ 'is-selected': account.uuid === selected }">
                <div v-if="pendingRemovalUUID === account.uuid" class="account-confirm">
                    <span>Remove {{ account.username }}?</span>
                    <div class="account-confirm-actions">
                        <button class="account-confirm-no" type="button" :disabled="removingAccount"
                            @click="pendingRemovalUUID = null">No</button>
                        <button class="account-confirm-yes" type="button" :disabled="removingAccount"
                            @click="removeAccount(account)">{{ removingAccount ? "Removing…" : "Yes" }}</button>
                    </div>
                </div>
                <template v-else>
                    <button class="account-menu-select" type="button" @click="selectAccount(account.uuid)">
                        <span class="account-type-icon">
                            <img class="account-menu-avatar"
                                :src="account.type === 'msa' ? '/assets/portraits/steve.png' : '/assets/portraits/offline-steve.png'"
                                alt="">
                        </span>
                        <span class="account-menu-copy">
                            <span class="launch-selection-title">{{ account.username }}</span>
                            <span class="launch-selection-detail">
                                {{ account.type === "msa" ? (account.expired ? "Auth Expired (Login again)" : "Microsoft Account") : "Offline Account" }}
                            </span>
                        </span>
                        <span v-if="account.uuid === selected" class="account-selected-tag">Selected</span>
                    </button>
                    <button class="account-remove" type="button" aria-label="Remove account"
                        :disabled="removingAccount" @click.stop="toggleRemoveConfirmation(account.uuid)">
                        <i class="fa-solid fa-trash-can" aria-hidden="true"></i>
                    </button>
                </template>
            </div>
            <div v-if="menuError" class="account-menu-error">{{ menuError }}</div>
            <button class="launch-menu-action add-account" type="button" @click="showAddPrompt">
                <span class="launch-menu-action-icon"><i class="icon icon-add" style="--icon-size: 16px;"></i></span>
                <span>Add Account</span>
            </button>
        </div>
        <button class="launch-selection launch-account" type="button" @click="emit('toggle-picker')">
            <img class="launch-avatar" alt=""
                :src="!selectedAccount ? '/assets/portraits/question-mark.png' : selectedAccount.type === 'msa' ? '/assets/portraits/steve.png' : '/assets/portraits/offline-steve.png'">
            <span class="launch-selection-copy">
                <span class="launch-selection-title">{{ selectedAccount?.username || "No account" }}</span>
                <span class="launch-selection-detail">
                    {{ selectedAccount ? (selectedAccount.type === "msa" ? "Microsoft account" : "Offline account") : "Sign in" }}
                </span>
            </span>
            <i class="icon icon-chevron" style="--icon-size: 16px; --icon-color: var(--launch-chevron-color);"></i>
        </button>

        <div v-if="promptOpen" class="account-overlay" @click.self="closePrompt"
            @keydown.esc.stop.prevent="closePrompt">
            <div class="account-popover" role="dialog" aria-modal="true" aria-labelledby="account-prompt-title">
                <header class="account-popover-head">
                    <span id="account-prompt-title">
                        {{ promptStep === "choose" ? "Add an account" : promptStep === "offline" ? "Offline account" : "Sign in with Microsoft" }}
                    </span>
                    <button class="quickwebview-nav-button account-prompt-close" type="button" aria-label="Close"
                        @click="closePrompt">
                        <i class="fa-solid fa-xmark"></i>
                    </button>
                </header>

                <template v-if="promptStep === 'choose'">
                    <button class="signin-option" type="button" @click="startBrowserLogin">
                        <i class="fa-brands fa-microsoft"></i>
                        <span>Sign In with Microsoft</span>
                    </button>
                    <button class="signin-option" type="button" @click="startDeviceLogin">
                        <i class="fa-solid fa-mobile-screen-button"></i>
                        <span>Sign In with Microsoft (Device Code)</span>
                    </button>
                    <button class="signin-option" type="button" @click="promptStep = 'offline'">
                        <i class="fa-solid fa-user"></i>
                        <span>Sign In with Username (Offline)</span>
                    </button>
                </template>

                <form v-else-if="promptStep === 'offline'" class="offline-account-form" @submit.prevent="addOfflineAccount">
                    <div class="offline-field">
                        <div class="offline-input-wrap">
                            <i class="fa-solid fa-user" aria-hidden="true"></i>
                            <input id="offline-username" v-model="offlineUsername" maxlength="16"
                                placeholder="Enter Username" autocomplete="off" spellcheck="false" autofocus>
                        </div>
                        <small>Note: Official servers require Microsoft authentication.</small>
                    </div>
                    <div class="offline-form-actions">
                        <button class="offline-back" type="button" @click="promptStep = 'choose'; errorMessage = ''">
                            Back
                        </button>
                        <button class="account-prompt-submit" type="submit" :disabled="busy">
                            {{ busy ? "Adding…" : "Add account" }}
                        </button>
                    </div>
                </form>

                <div v-else-if="promptStep === 'device'" class="account-ms">
                    <div v-if="deviceLogin" class="ms-body">
                        <div class="ms-qr-housing">
                            <div ref="qrContainer" class="ms-qr" role="img" aria-label="Microsoft sign-in QR code"></div>
                        </div>
                        <div class="ms-instructions">
                            <p>Open <a class="ms-link" href="#"
                                    @click.prevent="BrowserOpenURL(deviceLogin.verification_uri)">{{ deviceLogin.verification_uri }}</a> on your browser or scan the QR code, then enter:</p>
                            <div class="ms-code-row">
                                <code>{{ deviceLogin.user_code }}</code>
                                <button class="quickwebview-nav-button" type="button"
                                    :title="codeCopied ? 'Copied' : 'Copy code'" @click="copyDeviceCode">
                                    <i :class="codeCopied ? 'fa-solid fa-check' : 'fa-solid fa-copy'"></i>
                                </button>
                            </div>
                            <div class="ms-hint">{{ authStatus }}</div>
                        </div>
                    </div>
                    <div v-else class="ms-status">{{ authStatus }}</div>
                </div>
                <div v-else-if="promptStep === 'browser'" class="account-ms">
                    <div class="ms-status">
                        <template v-if="browserLoginURL">
                            Finish signing in in your browser. <br>If it didn't open,
                            <a class="ms-link" :href="browserLoginURL"
                                @click.prevent="BrowserOpenURL(browserLoginURL)">click here</a>.
                        </template>
                        <template v-else>{{ authStatus }}</template>
                    </div>
                </div>

                <p v-if="errorMessage" class="account-prompt-error" role="alert">{{ errorMessage }}</p>
            </div>
        </div>
    </div>
</template>

<style scoped>
.launch-picker { position: relative; flex: 0 1 min(250px, 30%); min-width: 0; }
.launch-selection { box-sizing: border-box; display: flex; align-items: center; gap: 0.5rem; width: min(250px, 30%); min-width: 0; height: 54px; padding: 0.5rem 0.75rem; border: 1px solid transparent; border-radius: 7px; background: var(--launch-card-bg-color); color: var(--text-color); text-align: left; cursor: pointer; }
.launch-picker .launch-selection { width: 100%; }
.launch-selection:hover { background: var(--launch-card-hover-bg-color); }
.launch-avatar { display: block; flex: 0 0 30px; width: 30px; height: 30px; border-radius: 6px; background: var(--surface-5); object-fit: cover; image-rendering: pixelated; }
.launch-menu { position: absolute; bottom: calc(100% + 0.45rem); left: 0; z-index: 5; box-sizing: border-box; display: flex; flex-direction: column; width: 100%; max-height: min(320px, calc(100vh - 110px)); overflow-x: hidden; overflow-y: auto; border-radius: 8px; background: var(--surface-4); }
.launch-menu-action { box-sizing: border-box; display: flex; flex: 0 0 auto; align-items: center; gap: 0.65rem; min-width: 0; min-height: 36px; padding: 0.15rem 0.75rem; border: 0; border-bottom: 1px solid var(--gray-color); background: var(--surface-4); color: var(--text-color); text-align: left; cursor: pointer; font-size: 0.8rem; font-weight: 600; }
.launch-menu-action:hover { background: var(--surface-hover-color); }
.launch-menu-action-icon { display: inline-flex; flex: 0 0 16px; align-items: center; justify-content: center; width: 16px; color: var(--text-muted-color); }
.launch-selection-copy { display: flex; flex: 1; flex-direction: column; gap: 0.15rem; min-width: 0; }
.launch-selection-title, .launch-selection-detail { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.launch-selection-title { font-size: 0.875rem; font-weight: 600; }
.launch-selection-detail { color: var(--text-muted-color); font-size: 0.75rem; opacity: 0.7; }
.account-menu-row { display: flex; align-items: center; gap: 0.35rem; padding: 0.6rem 0.9rem; border-bottom: 1px solid rgba(255, 255, 255, 0.05); background: rgba(255, 255, 255, 0.03); cursor: pointer; transition: background 0.14s ease; }
.account-menu-row:hover { background: rgba(255, 68, 68, 0.1); }
.account-menu-row.is-selected { background: rgba(255, 68, 68, 0.08); }
.account-menu-select { display: flex; flex: 1; align-items: center; gap: 0.65rem; min-width: 0; padding: 0; border: 0; background: transparent; color: var(--text-color); text-align: left; cursor: pointer; }
.account-type-icon { display: inline-flex; flex: 0 0 28px; align-items: center; justify-content: center; width: 28px; height: 28px; overflow: hidden; border-radius: 6px; background: rgba(255, 255, 255, 0.05); }
.account-menu-avatar { width: 100%; height: 100%; object-fit: contain; image-rendering: pixelated; }
.account-menu-copy { display: flex; flex: 1; flex-direction: column; min-width: 0; }
.account-menu-copy .launch-selection-title { font-size: 0.92rem; font-weight: 500; }
.account-menu-copy .launch-selection-detail { margin-top: 2px; font-size: 0.74rem; opacity: 1; }
.account-selected-tag { flex: 0 0 auto; margin-left: 0.15rem; padding: 0.22rem 0.42rem; border-radius: 4px; background: rgba(57, 153, 77, 0.2); color: #8ee39a; font-size: 0.65rem; font-weight: 600; }
.account-remove { display: inline-flex; flex: 0 0 28px; align-items: center; justify-content: center; width: 28px; height: 28px; border: 0; border-radius: 6px; background: transparent; color: rgba(255, 255, 255, 0.45); cursor: pointer; transition: background 0.14s ease, color 0.14s ease; }
.account-remove:hover { background: rgba(255, 68, 68, 0.16); color: #ff6b6b; }
.account-remove:disabled { opacity: 0.5; cursor: default; }
.account-confirm { display: flex; flex: 1; align-items: center; justify-content: space-between; gap: 0.6rem; min-width: 0; color: var(--text-color); font-size: 0.82rem; }
.account-confirm-actions { display: flex; flex: 0 0 auto; gap: 0.4rem; }
.account-confirm-no, .account-confirm-yes { height: 26px; min-width: 46px; border: 1px solid rgba(255, 255, 255, 0.14); border-radius: 6px; background: transparent; color: var(--text-color); font-size: 0.76rem; cursor: pointer; }
.account-confirm-yes { border-color: #ff4444; background: #ff4444; }
.account-confirm-yes:hover { background: #ff5e5e; }
.account-confirm-no:disabled, .account-confirm-yes:disabled { opacity: 0.55; cursor: default; }
.launch-menu-action.add-account { position: sticky; bottom: 0; z-index: 1; gap: 0.6rem; padding: 0.45rem 0.9rem; border-top: 1px solid rgba(255, 255, 255, 0.08); border-bottom: 0; background: var(--surface-4); }
.launch-menu-action.add-account:hover { background: rgba(255, 68, 68, 0.1); }
.account-empty, .account-menu-error { padding: 0.75rem; color: var(--text-muted-color); font-size: 0.8rem; }
.account-menu-error, .account-prompt-error { color: #ff8585; overflow-wrap: anywhere; }
.account-overlay { box-sizing: border-box; position: fixed; inset: 0; z-index: 1200; display: flex; align-items: center; justify-content: center; padding: 1.5rem; background: rgba(15, 15, 15, 0.92); }
.account-popover { box-sizing: border-box; display: flex; flex-direction: column; gap: 0.4rem; width: min(400px, 100%); max-height: 80vh; overflow: auto; padding: 0.9rem; border-radius: 6px; background: var(--surface-5); box-shadow: 0 12px 28px rgba(0, 0, 0, 0.5); }
.account-popover-head { display: flex; align-items: center; justify-content: space-between; color: var(--text-color); font-size: 0.95rem; font-weight: 700; }
.quickwebview-nav-button { box-sizing: border-box; display: inline-flex; flex: 0 0 32px; align-items: center; justify-content: center; width: 32px; height: 32px; border: 1px solid rgba(255, 255, 255, 0.08); border-radius: 7px; background: rgba(255, 255, 255, 0.04); color: rgba(255, 255, 255, 0.78); cursor: pointer; }
.quickwebview-nav-button:hover { border-color: rgba(255, 68, 68, 0.35); background: rgba(255, 68, 68, 0.12); color: var(--text-color); }
.quickwebview-nav-button:active { background: rgba(255, 68, 68, 0.18); }
.account-prompt-close, .account-prompt-close:hover, .account-prompt-close:active { border-color: transparent; background: transparent; }
.signin-option { display: flex; align-items: center; gap: 0.6rem; min-height: 42px; padding: 0 0.9rem; border: 0; border-radius: 6px; background: rgba(255, 255, 255, 0.03); color: var(--text-color); font-size: 0.86rem; font-weight: 600; cursor: pointer; transition: border-color 0.14s ease, color 0.14s ease; }
.signin-option:hover { color: #ff6b6b; }
.account-prompt-submit { min-height: 32px; padding: 0.4rem 0.7rem; border: 1px solid var(--gray-color); border-radius: 7px; background: var(--accent-color); color: var(--text-color); font: inherit; font-size: 0.8rem; text-align: center; cursor: pointer; }
.account-prompt-submit:hover { background: var(--surface-hover-color); }
.offline-account-form { display: flex; flex-direction: column; gap: 1rem; }
.offline-field { display: flex; flex-direction: column; gap: 0.45rem; }
.offline-field label { font-size: 0.82rem; font-weight: 600; }
.offline-field small { color: var(--text-muted-color); font-size: 0.74rem; }
.offline-input-wrap { display: flex; align-items: center; gap: 0.6rem; min-height: 42px; padding: 0 0.7rem; border: 1px solid var(--gray-color); border-radius: 7px; background: var(--surface-2); color: var(--text-muted-color); }
.offline-input-wrap:focus-within { border-color: var(--accent-color); }
.offline-input-wrap input { flex: 1; min-width: 0; padding: 0.65rem 0; border: 0; outline: 0; background: transparent; color: var(--text-color); font: inherit; }
.offline-form-actions { display: flex; justify-content: flex-end; gap: 0.55rem; }
.offline-back { min-height: 32px; padding: 0 0.65rem; border: 0; border-radius: 7px; background: transparent; color: var(--text-muted-color); font: inherit; font-size: 0.8rem; cursor: pointer; }
.offline-back:hover { background: var(--surface-hover-color); color: var(--text-color); }
.account-prompt-submit:disabled { opacity: 0.6; cursor: default; }
.ms-status { padding: 0.4rem 0; color: rgba(255, 255, 255, 0.7); font-size: 0.86rem; }
.ms-body { display: flex; align-items: flex-start; gap: 0.9rem; }
.ms-qr-housing { box-sizing: border-box; display: flex; flex: 0 0 148px; align-items: center; justify-content: center; padding: 6px; border-radius: 8px; background: #fff; }
.ms-qr { box-sizing: border-box; display: flex; flex: 0 0 132px; align-items: center; justify-content: center; width: 132px; height: 132px; }
.ms-qr img, .ms-qr canvas { display: block; width: 100%; height: 100%; }
.ms-instructions { min-width: 0; color: rgba(255, 255, 255, 0.75); font-size: 0.84rem; line-height: 1.4; }
.ms-instructions p { margin: 0; }
.ms-link { color: #ff6b6b; text-decoration: none; word-break: break-all; }
.ms-link:hover { text-decoration: underline; }
.ms-code-row { display: flex; align-items: center; gap: 0.5rem; margin: 0.5rem 0; }
.ms-code-row code { flex: 1; padding: 0.4rem 0.6rem; border-radius: 6px; background: rgba(255, 255, 255, 0.06); color: #fff; font-size: 1.05rem; font-weight: 700; letter-spacing: 0.12em; text-align: center; }
.ms-code-row button { flex: 0 0 32px; min-height: 32px; padding: 0; }
.ms-hint { color: rgba(255, 255, 255, 0.45); font-size: 0.76rem; }

@media (max-width: 360px) {
    .ms-body { flex-direction: column; align-items: center; }
    .ms-instructions { width: 100%; }
}
</style>
