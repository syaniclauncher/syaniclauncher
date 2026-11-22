<script setup>
import { computed, ref } from "vue";

const props = defineProps({ instance: { type: Object, required: true } });
const emit = defineEmits(["back", "play", "edit", "duplicate", "delete", "export"]);
const activeTab = ref("Content");
const contentType = ref("Mods");
const search = ref("");
const moreOpen = ref(false);
const categories = ["Mods", "Resource Packs", "Shaders", "Data Packs"];
const tabs = [
    { name: "Content", icon: "fa-layer-group" },
    { name: "Worlds", icon: "fa-globe" },
    { name: "Gallery", icon: "fa-image" },
    { name: "Files", icon: "fa-folder" },
];

const loaderLabel = computed(() => {
    const labels = { vanilla: "Vanilla", forge: "Forge", fabric: "Fabric", quilt: "Quilt", neoforge: "NeoForge" };
    return labels[(props.instance.loader || "vanilla").toLowerCase()] || "Vanilla";
});
const iconUrl = computed(() => props.instance.iconUrl || `/assets/instance_icons/${props.instance.icon || "grass_block"}.png`);
const lastPlayed = computed(() => {
    if (!props.instance.lastPlayed) return "Never played";
    const days = Math.floor((Date.now() / 1000 - props.instance.lastPlayed) / 86400);
    if (days <= 0) return "Today";
    if (days === 1) return "Yesterday";
    return `${days} days ago`;
});

function setTab(tab) {
    activeTab.value = tab;
    search.value = "";
}
</script>

<template>
    <main class="instance-page">
        <button type="button" class="back-button" @click="emit('back')">
            <i class="fas fa-arrow-left" aria-hidden="true"></i>
            <span>Back to Library</span>
        </button>

        <section class="instance-overview">
            <header class="overview-header">
                <div class="instance-title">
                    <div class="overview-icon"><img :src="iconUrl" :alt="instance.name"></div>
                    <div class="overview-copy">
                        <div class="version-chip"><i class="fas fa-cube" aria-hidden="true"></i> {{ loaderLabel }} <span>|</span> MC {{ instance.version || "Unknown" }}</div>
                        <h1>{{ instance.name }}</h1>
                    </div>
                </div>
                <button type="button" class="play-button" @click="emit('play', instance)">
                    <i class="fas fa-play" aria-hidden="true"></i>
                    <span>Play</span>
                </button>
            </header>
            <div class="overview-meta">
                <div class="last-played-block">
                    <span class="meta-label">Last Played</span>
                    <span><i class="far fa-calendar" aria-hidden="true"></i> {{ lastPlayed }}</span>
                </div>
                <button type="button" class="overview-action" @click="emit('edit', instance)">
                    <i class="fas fa-cog" aria-hidden="true"></i>
                    <span>Instance Settings</span>
                </button>
                <div class="more-wrap">
                    <button type="button" class="overview-action more-button" :aria-expanded="moreOpen" @click="moreOpen = !moreOpen">
                        <i class="fas fa-ellipsis-v" aria-hidden="true"></i>
                        <span>More</span>
                    </button>
                    <div v-if="moreOpen" class="more-menu">
                        <button type="button" disabled><i class="far fa-folder"></i> Open Folder</button>
                        <button type="button" @click="emit('duplicate', instance); moreOpen = false"><i class="far fa-copy"></i> Duplicate</button>
                        <button type="button" @click="emit('export', instance); moreOpen = false"><i class="fas fa-upload"></i> Export</button>
                        <button type="button" disabled><i class="fas fa-external-link-alt"></i> Create Shortcut</button>
                        <button type="button" class="delete-item" @click="emit('delete', instance); moreOpen = false"><i class="far fa-trash-alt"></i> Delete</button>
                    </div>
                </div>
            </div>
        </section>

        <nav class="instance-tabs" aria-label="Instance sections">
            <button v-for="tab in tabs" :key="tab.name" type="button" :class="{ active: activeTab === tab.name }" @click="setTab(tab.name)">
                <i class="fas" :class="tab.icon" aria-hidden="true"></i>
                <span>{{ tab.name }}</span>
            </button>
            <span class="tabs-spacer"></span>
            <button type="button" class="display-options" title="Display options" disabled><i class="fas fa-sliders-h"></i></button>
        </nav>

        <section v-if="activeTab === 'Content'" class="content-panel">
            <div class="content-toolbar">
                <div class="content-categories" aria-label="Content category">
                    <button v-for="category in categories" :key="category" type="button"
                        :class="{ active: contentType === category }" @click="contentType = category">
                        <i class="fas" :class="category === 'Mods' ? 'fa-cube' : category === 'Resource Packs' ? 'fa-palette' : category === 'Shaders' ? 'fa-sun' : 'fa-box-open'" aria-hidden="true"></i>
                        <span>{{ category }}</span>
                    </button>
                </div>
                <label class="content-search">
                    <i class="fas fa-search" aria-hidden="true"></i>
                    <input v-model="search" type="search" :placeholder="`Search installed ${contentType.toLowerCase()}...`">
                </label>
                <button type="button" class="refresh-button" title="Refresh" disabled><i class="fas fa-sync-alt"></i></button>
                <button type="button" class="add-content-button" title="Content management is not available in this frontend build" disabled>
                    <i class="fas fa-plus" aria-hidden="true"></i><span>Add</span>
                </button>
            </div>
            <div class="content-empty">
                <div class="empty-icon"><i class="fas fa-box-open" aria-hidden="true"></i></div>
                <strong>No {{ contentType.toLowerCase() }} installed</strong>
                <span>This instance does not have any {{ contentType.toLowerCase() }} yet.</span>
                <button type="button" disabled><i class="fas fa-plus" aria-hidden="true"></i> Browse {{ contentType.toLowerCase() }}</button>
            </div>
        </section>

        <section v-else class="content-panel simple-empty">
            <div class="content-empty">
                <div class="empty-icon"><i class="fas" :class="activeTab === 'Worlds' ? 'fa-globe' : activeTab === 'Gallery' ? 'fa-image' : 'fa-folder-open'" aria-hidden="true"></i></div>
                <strong>No {{ activeTab.toLowerCase() }} yet</strong>
                <span>{{ activeTab }} for this instance will appear here.</span>
            </div>
        </section>
    </main>
</template>

<style scoped>
.instance-page {
    box-sizing: border-box;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: .75rem;
    overflow: auto;
    padding: .8rem 1rem 1rem;
    background: var(--surface-3);
}
.back-button {
    align-self: flex-start;
    display: inline-flex;
    align-items: center;
    gap: .5rem;
    border: 0;
    padding: .4rem .15rem;
    background: transparent;
    color: var(--text-muted-color);
    font: inherit;
    font-size: .82rem;
    cursor: pointer;
}
.back-button:hover { color: var(--text-color); }
.instance-overview { overflow: visible; border: 1px solid rgba(255,255,255,.07); border-radius: 8px; background: var(--surface-4); }
.overview-header { min-height: 100px; display: flex; align-items: center; justify-content: space-between; gap: 1rem; padding: 1rem 1.2rem; }
.instance-title { min-width: 0; display: flex; align-items: center; gap: .9rem; }
.overview-icon { width: 66px; height: 66px; flex: 0 0 66px; display: grid; place-items: center; border: 1px solid rgba(255,255,255,.06); border-radius: 8px; background: var(--surface-5); }
.overview-icon img { width: 48px; height: 48px; object-fit: contain; image-rendering: pixelated; }
.overview-copy { min-width: 0; }
.version-chip { width: fit-content; display: flex; align-items: center; gap: .4rem; margin-bottom: .35rem; border: 1px solid rgba(255,255,255,.08); border-radius: 12px; padding: .15rem .45rem; color: var(--text-muted-color); font-size: .67rem; white-space: nowrap; }
.version-chip i { color: #6ea985; }
.version-chip span { color: var(--text-dim-color); }
.overview-copy h1 { overflow: hidden; margin: 0; font-size: 1.35rem; text-overflow: ellipsis; white-space: nowrap; }
.play-button { min-width: 150px; height: 40px; display: inline-flex; align-items: center; justify-content: center; gap: .5rem; border: 1px solid #db514c; border-radius: 6px; background: var(--accent-color); color: white; font: inherit; font-size: .84rem; font-weight: 750; cursor: pointer; }
.play-button:hover { background: #e84b46; }
.overview-meta { min-height: 52px; display: flex; border-top: 1px solid rgba(255,255,255,.06); }
.last-played-block { min-width: 0; flex: 1; display: flex; flex-direction: column; justify-content: center; gap: .25rem; padding: .45rem 1.2rem; border-right: 1px solid rgba(255,255,255,.06); color: var(--text-muted-color); font-size: .75rem; }
.meta-label { color: var(--text-dim-color); font-size: .62rem; letter-spacing: .12em; text-transform: uppercase; }
.overview-action { min-width: 96px; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: .25rem; border: 0; border-right: 1px solid rgba(255,255,255,.06); padding: .3rem .6rem; background: transparent; color: var(--text-muted-color); font: inherit; font-size: .67rem; cursor: pointer; }
.overview-action:hover { color: var(--text-color); background: rgba(255,255,255,.025); }
.overview-action i { font-size: .85rem; }
.more-wrap { position: relative; display: flex; }
.more-button { min-width: 52px; border-right: 0; }
.more-menu { position: absolute; z-index: 4; top: calc(100% + 2px); right: 0; width: 175px; padding: .25rem 0; border: 1px solid var(--gray-color); border-radius: 5px; background: var(--surface-4); box-shadow: 0 6px 16px rgba(0,0,0,.35); }
.more-menu button { width: 100%; min-height: 33px; display: flex; align-items: center; gap: .55rem; border: 0; padding: 0 .7rem; background: transparent; color: var(--text-muted-color); text-align: left; font: inherit; font-size: .73rem; cursor: pointer; }
.more-menu button:hover:not(:disabled) { background: var(--surface-5); color: var(--text-color); }
.more-menu button i { width: 14px; }
.more-menu button:disabled { opacity: .4; cursor: default; }
.more-menu .delete-item { color: #ff7777; }
.instance-tabs { min-height: 42px; display: flex; align-items: stretch; gap: .3rem; border-bottom: 1px solid rgba(255,255,255,.08); }
.instance-tabs button { min-width: 94px; display: inline-flex; align-items: center; justify-content: center; gap: .45rem; border: 0; border-radius: 5px 5px 0 0; padding: 0 .8rem; background: transparent; color: var(--text-muted-color); font: inherit; font-size: .78rem; cursor: pointer; }
.instance-tabs button:hover { color: var(--text-color); }
.instance-tabs button.active { border-bottom: 2px solid var(--accent-color); background: rgba(255,68,68,.08); color: var(--accent-color); }
.instance-tabs .tabs-spacer { flex: 1; }
.instance-tabs button.display-options { min-width: 38px; color: var(--text-dim-color); }
.content-panel { min-height: 0; display: flex; flex-direction: column; flex: 1; }
.content-toolbar { min-height: 42px; display: flex; align-items: center; gap: .45rem; }
.content-categories { display: flex; align-items: center; gap: .2rem; }
.content-categories button { height: 30px; display: inline-flex; align-items: center; gap: .35rem; border: 1px solid transparent; border-radius: 5px; padding: 0 .55rem; background: transparent; color: var(--text-dim-color); font: inherit; font-size: .68rem; cursor: pointer; }
.content-categories button.active { border-color: rgba(255,68,68,.38); background: rgba(255,68,68,.1); color: var(--accent-color); }
.content-categories button:hover:not(.active) { background: var(--surface-4); color: var(--text-muted-color); }
.content-search { height: 30px; min-width: 100px; flex: 1; display: flex; align-items: center; gap: .45rem; border: 1px solid rgba(255,255,255,.08); border-radius: 4px; padding: 0 .55rem; background: var(--surface-1); color: var(--text-dim-color); }
.content-search input { width: 100%; border: 0; outline: 0; background: transparent; color: var(--text-color); font: inherit; font-size: .7rem; }
.content-search input::placeholder { color: var(--text-dim-color); }
.refresh-button, .add-content-button { height: 30px; border: 1px solid var(--gray-color); border-radius: 4px; background: var(--surface-4); color: var(--text-muted-color); }
.refresh-button { width: 31px; }
.add-content-button { display: flex; align-items: center; gap: .35rem; padding: 0 .65rem; background: var(--accent-color); color: white; font: inherit; font-size: .72rem; }
.refresh-button:disabled, .add-content-button:disabled { opacity: .5; cursor: not-allowed; }
.content-empty { flex: 1; min-height: 180px; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: .45rem; border: 1px solid rgba(255,255,255,.06); border-radius: 7px; background: var(--surface-4); text-align: center; }
.empty-icon { width: 42px; height: 42px; display: grid; place-items: center; margin-bottom: .25rem; border-radius: 12px; background: var(--surface-5); color: var(--text-dim-color); }
.content-empty strong { font-size: .84rem; }
.content-empty > span { color: var(--text-dim-color); font-size: .72rem; }
.content-empty > button { min-height: 30px; margin-top: .55rem; border: 1px solid var(--gray-color); border-radius: 5px; padding: 0 .65rem; background: var(--surface-4); color: var(--text-muted-color); font: inherit; font-size: .7rem; }
.content-empty > button:disabled { cursor: not-allowed; }
.simple-empty { flex: 1; }
@media (max-width: 760px) {
    .instance-page { padding: .6rem; }
    .overview-header { padding: .7rem; }
    .overview-meta .last-played-block { padding: .4rem .55rem; }
    .overview-action { min-width: 58px; padding: .25rem; font-size: .58rem; }
    .play-button { min-width: 84px; }
    .instance-tabs { overflow-x: auto; }
    .instance-tabs button { min-width: 76px; }
    .content-toolbar { flex-wrap: wrap; padding: .35rem 0; }
    .content-categories { width: 100%; overflow-x: auto; }
}
</style>
