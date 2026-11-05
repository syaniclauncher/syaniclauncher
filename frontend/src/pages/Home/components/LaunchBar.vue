<script setup>
import { ref } from "vue";

const emit = defineEmits(["navigate"]);
const open = ref(""); // id of the open picker
</script>

<template>
    <section class="home-launch-bar" @click="open = ''" @keydown.esc="open = ''">
        <div class="launch-picker launch-account-picker" :class="{ 'is-open': open === 'account' }" @click.stop>
            <div v-show="open === 'account'" class="launch-menu launch-account-menu">
                <button class="launch-menu-action add-account" type="button" @click="open = ''">
                    <span class="launch-menu-action-icon"><i class="icon icon-add" style="--icon-size: 16px;"></i></span>
                    <span>Add Account</span>
                </button>
            </div>
            <button class="launch-selection launch-account" type="button"
                @click="open = open === 'account' ? '' : 'account'">
                <img class="launch-avatar" src="/assets/portraits/question-mark.png" alt="">
                <span class="launch-selection-copy">
                    <span class="launch-selection-title">No account</span>
                    <span class="launch-selection-detail">Sign in</span>
                </span>
                <i class="icon icon-chevron" style="--icon-size: 16px; --icon-color: var(--launch-chevron-color);"></i>
            </button>
        </div>

        <div class="launch-picker launch-instance-picker" :class="{ 'is-open': open === 'instance' }" @click.stop>
            <div v-show="open === 'instance'" class="launch-menu launch-instance-menu">
                <a class="launch-menu-action manage-instances" href="#" @click.prevent="emit('navigate', '/instances'); open = ''">
                    <span class="launch-menu-action-icon manage-instances-icon"><i class="icon icon-instances" style="--icon-size: 16px;"></i></span>
                    <span>Manage Instances</span>
                </a>
            </div>
            <button class="launch-selection launch-instance" type="button"
                @click="open = open === 'instance' ? '' : 'instance'">
                <img class="launch-instance-icon" src="/assets/instance_icons/grass_block.png" alt="">
                <span class="launch-selection-copy">
                    <span class="launch-selection-title">No instance selected</span>
                    <span class="launch-selection-detail">Pick in Instances</span>
                </span>
                <i class="icon icon-chevron" style="--icon-size: 16px; --icon-color: var(--launch-chevron-color);"></i>
            </button>
        </div>

        <div class="launch-actions">
            <button class="launch-folder" type="button" title="Open game folder">
                <i class="icon icon-folder" style="--icon-size: 20px;"></i>
            </button>
            <button class="launch-game" type="button" disabled title="Game launching is not available yet">
                <i class="icon icon-rocket" style="--icon-size: 18px;"></i>
                <span>Launch Game</span>
            </button>
        </div>
    </section>
</template>

<style scoped>
.home-launch-bar { position: relative; z-index: 2; box-sizing: border-box; display: flex; align-items: center; gap: 0.7rem; min-height: 96px; padding: 1.4rem 1.25rem; border-top: 1px solid var(--launch-bar-border-color); background: var(--launch-bar-bg-color); }
.launch-picker { position: relative; flex: 0 1 min(250px, 30%); min-width: 0; }
.launch-selection { box-sizing: border-box; display: flex; align-items: center; gap: 0.5rem; width: min(250px, 30%); min-width: 0; height: 54px; padding: 0.5rem 0.75rem; border: 1px solid transparent; border-radius: 7px; background: var(--launch-card-bg-color); color: var(--text-color); text-align: left; cursor: pointer; }
.launch-picker .launch-selection { width: 100%; }
.launch-selection:hover { background: var(--launch-card-hover-bg-color); }
.launch-avatar, .launch-instance-icon { display: block; flex: 0 0 30px; width: 30px; height: 30px; border-radius: 6px; object-fit: cover;  }
.launch-avatar {image-rendering: pixelated;}
.launch-menu { position: absolute; bottom: calc(100% + 0.45rem); left: 0; z-index: 5; box-sizing: border-box; display: flex; flex-direction: column; width: 100%; max-height: min(320px, calc(100vh - 110px)); overflow-x: hidden; overflow-y: auto; border-radius: 8px; background: var(--surface-4); }
.launch-menu-item, .launch-menu-action { box-sizing: border-box; display: flex; flex: 0 0 auto; align-items: center; gap: 0.65rem; min-width: 0; padding: 0.15rem 0.75rem; border: 0; border-bottom: 1px solid var(--gray-color); background: var(--surface-4); color: var(--text-color); text-align: left; text-decoration: none; cursor: pointer; }
.launch-menu-item { display: flex; min-height: 57px; }
.launch-menu-item:hover, .launch-menu-action:hover { background: var(--surface-hover-color); }
.launch-menu-item.is-selected { background: var(--surface-7); }
.launch-menu-item-button { box-sizing: border-box; display: flex; flex: 1 1 auto; align-items: center; gap: 0.65rem; min-width: 0; padding: 0.15rem 0.75rem; border: 0; background: transparent; color: var(--text-color); text-align: left; text-decoration: none; cursor: pointer; }
.launch-menu-item-button:hover { background: transparent; }
.launch-menu-item .launch-selection-copy { gap: 0.12rem; }
.launch-menu-item .launch-selection-detail { opacity: 0.8; }
.launch-menu-item-delete { display: inline-flex; flex: 0 0 auto; align-items: center; justify-content: center; width: 28px; height: 28px; margin-right: 0.5rem; border: 0; border-radius: 6px; background: transparent; color: var(--text-muted-color); cursor: pointer; }
.launch-menu-item-delete:hover { background: rgba(255, 255, 255, 0.08); color: #ff7f7f; }
.launch-menu-badge { flex: 0 0 auto; padding: 0.2rem 0.4rem; border-radius: 4px; background: #24522a; color: #77df78; font-size: 0.65rem; }
.launch-menu-action { min-height: 36px; font-size: 0.8rem; font-weight: 600; }
.launch-menu-action:last-child { border-bottom: 0; }
.launch-menu-action-icon { display: inline-flex; flex: 0 0 16px; align-items: center; justify-content: center; width: 16px; color: var(--text-muted-color); font-size: 1.2rem; font-weight: 400; }
.manage-instances-icon { font-size: 1rem; }
.launch-picker.is-open { z-index: 6; }
.launch-selection-copy { display: flex; flex: 1; flex-direction: column; gap: 0.15rem; min-width: 0; }
.launch-selection-title, .launch-selection-detail { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.launch-selection-title { font-size: 0.875rem; font-weight: 600; }
.launch-selection-detail { color: var(--text-muted-color); font-size: 0.75rem; opacity: 0.58; }
.launch-chevron { flex: 0 0 16px; width: 16px; height: 16px; fill: none; stroke: var(--launch-chevron-color); stroke-linecap: round; stroke-linejoin: round; stroke-width: 2; }
.launch-actions { display: flex; align-items: center; gap: 0.75rem; margin-left: auto; }
.launch-folder, .launch-game { box-sizing: border-box; display: inline-flex; align-items: center; justify-content: center; height: 48px; border: 0; border-radius: 7px; color: var(--text-color); cursor: pointer; }
.launch-folder { width: 48px; background: var(--launch-card-bg-color); }
.launch-folder:hover { background: var(--launch-card-hover-bg-color); }
.launch-folder svg { width: 19px; height: 19px; fill: var(--launch-folder-icon-color); stroke: var(--launch-folder-icon-color); stroke-linecap: round; stroke-linejoin: round; stroke-width: 1.5; }
.launch-game { gap: 0.65rem; min-width: 186px; padding: 0 1.25rem; background: var(--accent-color); font: inherit; font-size: 0.95rem; font-weight: 700; transition: background 0.16s ease; }
.launch-game:disabled { opacity: 1; cursor: not-allowed; }
.launch-game svg { width: 19px; height: 19px; fill: currentColor; stroke: currentColor; stroke-linecap: round; stroke-linejoin: round; stroke-width: 1.5; }
</style>
