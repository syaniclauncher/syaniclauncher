<script setup>
import { computed, ref } from "vue";
import InstanceCard from "./Instances/components/InstanceCard.vue";
import InstanceDetails from "./Instances/components/InstanceDetails.vue";
import InstallationForm from "./Instances/components/InstallationForm.vue";

const STORAGE_KEY = "syanic.instances";
const SELECTED_KEY = "syanic.selectedInstance";
const instances = ref(loadInstances());
const search = ref("");
const sortBy = ref("lastPlayed");
const versionFilters = ref({ releases: true, snapshots: true, modded: true });
const layout = ref("grid");
const detailId = ref("");
const formOpen = ref(false);
const editingInstance = ref(null);
const importInput = ref(null);
const notice = ref("");
const emit = defineEmits(["navigate"]);

function loadInstances() {
    try {
        const value = JSON.parse(localStorage.getItem(STORAGE_KEY) || "[]");
        return Array.isArray(value) ? value : [];
    } catch {
        return [];
    }
}

function persist() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(instances.value));
}

const filteredInstances = computed(() => {
    const query = search.value.trim().toLowerCase();
    return instances.value
        .filter((instance) => {
            const matchesSearch = !query || `${instance.name} ${instance.version} ${instance.loader}`.toLowerCase().includes(query);
            const versionFilter = instance.versionType === "snapshot" ? "snapshots" : "releases";
            const isModded = (instance.loader || "vanilla").toLowerCase() !== "vanilla";
            return matchesSearch && versionFilters.value[versionFilter] && (!isModded || versionFilters.value.modded);
        })
        .sort((a, b) => sortBy.value === "name"
            ? (a.name || "").localeCompare(b.name || "")
            : (b.lastPlayed || 0) - (a.lastPlayed || 0));
});

const detailInstance = computed(() => instances.value.find((instance) => instance.id === detailId.value) || null);

function openCreate() {
    editingInstance.value = null;
    formOpen.value = true;
}

function openEdit(instance) {
    editingInstance.value = instance;
    formOpen.value = true;
}

function closeForm() {
    formOpen.value = false;
    editingInstance.value = null;
}

function saveInstance(data) {
    if (editingInstance.value) {
        const index = instances.value.findIndex((instance) => instance.id === editingInstance.value.id);
        if (index >= 0) instances.value[index] = { ...instances.value[index], ...data };
    } else {
        const instance = {
            ...data,
            id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            lastPlayed: 0,
            playtimeSeconds: 0,
        };
        instances.value.push(instance);
        detailId.value = instance.id;
        localStorage.setItem(SELECTED_KEY, instance.id);
    }
    persist();
    closeForm();
}

function openDetails(instance) {
    detailId.value = instance.id;
    localStorage.setItem(SELECTED_KEY, instance.id);
    window.dispatchEvent(new CustomEvent("syanic:selected-instance", { detail: instance }));
}

function backToLibrary() {
    detailId.value = "";
}

function removeInstance(instance) {
    instances.value = instances.value.filter((item) => item.id !== instance.id);
    if (localStorage.getItem(SELECTED_KEY) === instance.id) {
        localStorage.removeItem(SELECTED_KEY);
        window.dispatchEvent(new CustomEvent("syanic:selected-instance", { detail: null }));
    }
    if (detailId.value === instance.id) detailId.value = "";
    persist();
}

function duplicateInstance(instance) {
    const copy = {
        ...instance,
        id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
        name: `${instance.name} (Copy)`,
        lastPlayed: 0,
    };
    instances.value.push(copy);
    persist();
}

function duplicateFromDetails(instance) {
    duplicateInstance(instance);
    detailId.value = "";
}

function selectInstance(instance) {
    localStorage.setItem(SELECTED_KEY, instance.id);
    sessionStorage.setItem("syanic.open-instance-menu", "true");
    window.dispatchEvent(new CustomEvent("syanic:selected-instance", { detail: instance }));
    emit("navigate", "/");
    // play the game after this
}

function exportInstance(instance) {
    const blob = new Blob([JSON.stringify(instance, null, 2)], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `${instance.name || "instance"}.json`;
    link.click();
    URL.revokeObjectURL(url);
}

function importModpack(event) {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) return;
    notice.value = `Importing ${file.name} needs backend support. No instance was created.`;
}

</script>

<template>
    <InstanceDetails v-if="detailInstance"
        :instance="detailInstance"
        @back="backToLibrary"
        @play="selectInstance"
        @edit="openEdit"
        @duplicate="duplicateFromDetails"
        @delete="removeInstance"
        @export="exportInstance" />

    <main v-else class="library-page">
        <header class="library-header">
            <div class="library-controls">
                <label class="control-field search-field">
                    <span class="control-label">Search</span>
                    <span class="library-search">
                        <i class="fas fa-search" aria-hidden="true"></i>
                        <input v-model="search" type="search" placeholder="Installation name">
                    </span>
                </label>
                <label class="control-field sort-field">
                    <span class="control-label">Sort By</span>
                    <span class="sort-select-wrap">
                        <select v-model="sortBy" aria-label="Sort instances">
                            <option value="lastPlayed">Latest Played</option>
                            <option value="name">Name</option>
                        </select>
                        <i class="fas fa-chevron-down" aria-hidden="true"></i>
                    </span>
                </label>
                <fieldset class="version-filters">
                    <legend class="control-label">Versions</legend>
                    <div class="version-filter-options">
                        <label><input v-model="versionFilters.releases" type="checkbox"><span>Releases</span></label>
                        <label><input v-model="versionFilters.snapshots" type="checkbox"><span>Snapshots</span></label>
                        <label><input v-model="versionFilters.modded" type="checkbox"><span>Modded</span></label>
                    </div>
                </fieldset>
            </div>

            <div class="library-actions">
                <div class="header-actions">
                    <button type="button" class="new-instance-button" @click="openCreate">
                        <i class="fas fa-plus" aria-hidden="true"></i>
                        <span>New Installation</span>
                    </button>
                    <label class="import-button">
                        <i class="fas fa-download" aria-hidden="true"></i>
                        <span>Import Modpack (.mrpack)</span>
                        <input ref="importInput" type="file" accept=".mrpack" @change="importModpack">
                    </label>
                </div>
                <div class="library-summary">
                    <span>{{ filteredInstances.length }} {{ filteredInstances.length === 1 ? "INSTALLATION" : "INSTALLATIONS" }}</span>
                    <div class="layout-toggle" aria-label="Instance layout">
                        <button type="button" :class="{ active: layout === 'grid' }" aria-label="Grid view" @click="layout = 'grid'">
                            <i class="fas fa-th" aria-hidden="true"></i>
                        </button>
                        <button type="button" :class="{ active: layout === 'list' }" aria-label="List view" @click="layout = 'list'">
                            <i class="fas fa-list" aria-hidden="true"></i>
                        </button>
                    </div>
                </div>
            </div>
        </header>

        <div v-if="notice" class="library-notice" role="status">
            <span>{{ notice }}</span>
            <button type="button" aria-label="Dismiss notice" @click="notice = ''"><i class="fas fa-times"></i></button>
        </div>

        <section class="instance-collection" :class="`layout-${layout}`" aria-label="Instances">
            <div v-if="filteredInstances.length === 0" class="empty-library">
                <i class="fas fa-cubes" aria-hidden="true"></i>
                <strong>{{ search ? "No matching instances" : "No instances yet" }}</strong>
                <span>{{ search ? "Try a different search." : "Create an instance to get started." }}</span>
            </div>
            <InstanceCard v-for="instance in filteredInstances" :key="instance.id"
                :instance="instance"
                :layout="layout"
                @open="openDetails"
                @play="selectInstance" />
        </section>
    </main>

    <InstallationForm v-if="formOpen"
        :instance="editingInstance"
        @close="closeForm"
        @save="saveInstance" />
</template>

<style scoped>
.library-page {
    box-sizing: border-box;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: .8rem;
    overflow: auto;
    padding: 1.5rem 1.7rem 1.2rem;
    background: var(--surface-3);
}
.library-header { display: flex; flex-direction: column; border-bottom: 1px solid var(--gray-color); }
.library-controls { display: grid; grid-template-columns: minmax(0,1fr) 126px auto; align-items: end; gap: .75rem; padding-bottom: .85rem; border-bottom: 1px solid var(--gray-color); }
.control-field { min-width: 0; display: flex; flex-direction: column; gap: .3rem; }
.control-label { color: var(--text-color); font-size: .65rem; font-weight: 800; text-transform: uppercase; }
.header-actions { display: flex; align-items: center; gap: .25rem; }
.library-actions { display: flex; align-items: center; justify-content: space-between; gap: 1rem; padding: .75rem 0; }
.library-summary { display: flex; align-items: center; gap: .65rem; }
.library-summary > span { color: var(--text-dim-color); font-size: .68rem; font-weight: 800; letter-spacing: .03em; white-space: nowrap; }
.library-notice { display: flex; align-items: center; justify-content: space-between; gap: .7rem; border: 1px solid rgba(211,161,88,.25); border-radius: 5px; padding: .55rem .7rem; background: rgba(211,161,88,.07); color: #d3b27e; font-size: .75rem; }
.library-notice button { border: 0; background: transparent; color: inherit; cursor: pointer; }
.import-button, .new-instance-button {
    position: relative;
    box-sizing: border-box;
    min-height: 32px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: .55rem;
    border: 1px solid var(--gray-color);
    border-radius: 2px;
    padding: 0 1rem;
    background: #171717;
    color: #eee;
    font: inherit;
    font-size: .74rem;
    font-weight: 700;
    cursor: pointer;
}
.import-button:hover, .new-instance-button:hover { background: #222; }
.import-button input { position: absolute; inset: 0; width: 100%; opacity: 0; cursor: pointer; }
.new-instance-button { border-color: var(--gray-color); background: transparent; }
.library-search {
    box-sizing: border-box;
    min-width: 0;
    height: 34px;
    display: flex;
    align-items: center;
    gap: .55rem;
    border: 1px solid var(--accent-color);
    border-radius: 4px;
    padding: 0 .75rem;
    background: var(--surface-1);
    color: var(--text-subtle-color);
}
.library-search:focus-within { border-color: rgba(255,255,255,.32); }
.library-search input { width: 100%; border: 0; outline: 0; background: transparent; color: var(--text-color); font: inherit; }
.library-search input::placeholder { color: var(--text-dim-color); }
.sort-select-wrap, .layout-toggle {
    height: 34px;
    display: flex;
    align-items: center;
    border: 1px solid var(--gray-color);
    border-radius: 4px;
    background: var(--surface-4);
    color: var(--text-muted-color);
}
.sort-select-wrap { position: relative; }
.sort-select-wrap select { width: 100%; height: 100%; border: 0; outline: 0; appearance: none; padding: 0 1.6rem 0 .7rem; background: transparent; color: var(--text-color); font: inherit; font-size: .76rem; cursor: pointer; }
.sort-select-wrap option { background: var(--surface-4); }
.sort-select-wrap > i { position: absolute; right: .55rem; top: 50%; transform: translateY(-50%); font-size: .65rem; pointer-events: none; }
.version-filters { min-width: 0; display: flex; flex-direction: column; gap: .3rem; margin: 0; border: 0; padding: 0; }
.version-filters legend { padding: 0; }
.version-filter-options { height: 34px; display: flex; align-items: center; gap: .6rem; }
.version-filter-options label { display: inline-flex; align-items: center; gap: .3rem; color: var(--text-muted-color); font-size: .72rem; font-weight: 650; white-space: nowrap; cursor: pointer; }
.version-filter-options input { width: 14px; height: 14px; margin: 0; accent-color: var(--accent-color); }
.layout-toggle { height: 28px; gap: .1rem; padding: 2px; }
.layout-toggle button { width: 22px; height: 22px; border: 0; border-radius: 3px; background: transparent; color: var(--text-muted-color); font-size: .68rem; cursor: pointer; }
.layout-toggle button.active { background: var(--surface-6); color: var(--accent-color); }
.instance-collection { min-height: 0; display: grid; align-content: start; gap: 1rem; overflow: auto; padding: .2rem .05rem 1rem; }
.layout-grid { grid-template-columns: repeat(auto-fill, minmax(245px, 1fr)); }
.layout-list { grid-template-columns: minmax(0, 1fr); gap: 0; border-top: 1px solid var(--gray-color); }
.empty-library {
    grid-column: 1 / -1;
    min-height: 150px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: .5rem;
    color: var(--text-subtle-color);
    text-align: center;
}
.empty-library i { margin-bottom: .2rem; font-size: 1.25rem; }
.empty-library strong { color: var(--text-muted-color); font-size: .9rem; }
.empty-library span { font-size: .78rem; }
@media (max-width: 850px) {
    .library-controls { grid-template-columns: minmax(0,1fr) 126px; }
    .version-filters { grid-column: 1 / -1; }
}
@media (max-width: 640px) {
    .library-page { gap: .85rem; padding: 1rem; }
    .library-controls { grid-template-columns: minmax(0,1fr); gap: .6rem; }
    .version-filters { grid-column: auto; }
    .version-filter-options { flex-wrap: wrap; }
    .library-actions, .header-actions { flex-wrap: wrap; gap: .35rem; }
    .import-button, .new-instance-button { min-height: 32px; padding: 0 .65rem; font-size: .72rem; }
    .library-summary { margin-left: auto; }
    .layout-grid { grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); }
}
</style>
