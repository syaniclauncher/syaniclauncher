<script setup>
import { computed, ref, watch } from "vue";
import IconSelector from "./IconSelector.vue";

const props = defineProps({ instance: { type: Object, default: null } });
const emit = defineEmits(["close", "save"]);
const showMoreVersions = ref(false);
const form = ref(blankForm());

const versions = [
    { id: "1.21.4", type: "release", latest: true },
    { id: "1.21.1", type: "release", latest: false },
    { id: "1.20.1", type: "release", latest: false },
    { id: "24w14potato", type: "snapshot", latest: false },
];
const visibleVersions = computed(() => {
    const choices = showMoreVersions.value ? versions : versions.filter((version) => version.latest);
    if (form.value.version && !choices.some((version) => version.id === form.value.version)) {
        return [{ id: form.value.version, type: form.value.versionType || "release", latest: false }, ...choices];
    }
    return choices;
});
const loaders = [
    { id: "vanilla", label: "Vanilla", icon: "fa-cube" },
    { id: "fabric", label: "Fabric", icon: "fa-leaf" },
    { id: "forge", label: "Forge", icon: "fa-hammer" },
    { id: "neoforge", label: "NeoForge", icon: "fa-fire" },
    { id: "quilt", label: "Quilt", icon: "fa-shapes" },
];

function blankForm() {
    return {
        name: "",
        loader: "vanilla",
        version: "1.21.4",
        loaderVersion: "",
        versionType: "release",
        icon: "grass_block",
        iconUrl: "",
        javaExec: "",
        jvmArgs: "",
        memoryMb: 4096,
        installed: false,
    };
}

watch(() => props.instance, (instance) => {
    showMoreVersions.value = false;
    form.value = instance ? { ...blankForm(), ...instance } : blankForm();
}, { immediate: true });

function changeIcon(value) {
    form.value.icon = value.icon;
    form.value.iconUrl = value.iconUrl;
}

function changeVersion(version) {
    form.value.version = version.id;
    form.value.versionType = version.type;
}

function submit() {
    const version = versions.find((item) => item.id === form.value.version);
    const loaderName = loaders.find((loader) => loader.id === form.value.loader)?.label || "Vanilla";
    const name = form.value.name.trim() || `${loaderName} ${form.value.version}`;
    emit("save", { ...form.value, versionType: version?.type || form.value.versionType, name });
}
</script>

<template>
    <div class="form-overlay" @click.self="emit('close')">
        <form class="instance-form" @submit.prevent="submit">
            <header class="form-header">
                <h2>{{ instance ? "Instance Settings" : "New Instance" }}</h2>
                <button type="button" class="close-button" aria-label="Close" @click="emit('close')">
                    <i class="fas fa-times" aria-hidden="true"></i>
                </button>
            </header>

            <div class="form-body">
                <div class="profile-row">
                    <div class="icon-field">
                        <span class="field-label">Icon</span>
                        <IconSelector :icon="form.icon" :icon-url="form.iconUrl" @change="changeIcon" />
                    </div>
                    <label class="name-field">
                        <span class="field-label">Instance Name</span>
                        <input v-model="form.name" type="text" placeholder="Enter name or leave empty for auto">
                    </label>
                </div>

                <div class="version-heading">
                    <label for="instance-version">Minecraft Version</label>
                    <label class="version-toggle">
                        <span>Show snapshots, historical and local</span>
                        <input v-model="showMoreVersions" type="checkbox">
                        <span class="toggle-track"></span>
                    </label>
                </div>
                <label class="version-select-wrap">
                    <select id="instance-version" :value="form.version"
                        @change="changeVersion(visibleVersions.find((version) => version.id === $event.target.value))">
                        <option v-for="version in visibleVersions" :key="version.id" :value="version.id">
                            {{ version.id }}{{ version.type === "snapshot" ? " (Snapshot)" : "" }}
                        </option>
                    </select>
                    <i class="fas fa-chevron-down" aria-hidden="true"></i>
                </label>

                <fieldset class="loader-field">
                    <legend>Mod Loader</legend>
                    <div class="loader-options">
                        <button v-for="loader in loaders" :key="loader.id" type="button"
                            class="loader-option" :class="[`loader-${loader.id}`, { active: form.loader === loader.id }]"
                            @click="form.loader = loader.id">
                            <i class="fas" :class="loader.icon" aria-hidden="true"></i>
                            <span>{{ loader.label }}</span>
                        </button>
                    </div>
                </fieldset>
            </div>

            <footer class="form-footer">
                <button type="button" class="cancel-button" @click="emit('close')">Cancel</button>
                <button type="submit" class="submit-button">
                    <i class="fas" :class="instance ? 'fa-check' : 'fa-plus'" aria-hidden="true"></i>
                    <span>{{ instance ? "Save Changes" : "Create Instance" }}</span>
                </button>
            </footer>
        </form>
    </div>
</template>

<style scoped>
.form-overlay {
    position: fixed;
    inset: 0;
    z-index: 900;
    display: grid;
    place-items: center;
    padding: 1rem;
    background: rgba(0,0,0,.68);
}
.instance-form {
    box-sizing: border-box;
    width: min(700px, 100%);
    max-height: min(94vh, 620px);
    display: flex;
    flex-direction: column;
    overflow: visible;
    border: 1px solid rgba(255,255,255,.09);
    border-radius: 12px;
    padding: 1.25rem 1.4rem 1rem;
    background: var(--surface-3);
    box-shadow: 0 18px 48px rgba(0,0,0,.5);
}
.form-header { display: flex; align-items: center; justify-content: space-between; padding-bottom: .9rem; }
.form-header h2 { margin: 0; font-size: 1.1rem; font-weight: 750; }
.close-button { width: 32px; height: 32px; border: 0; border-radius: 6px; background: var(--surface-5); color: var(--text-muted-color); cursor: pointer; }
.close-button:hover { color: var(--text-color); background: var(--surface-6); }
.form-body { display: flex; flex-direction: column; gap: 1.15rem; overflow: visible; padding: .2rem .05rem .8rem; }
.profile-row { display: grid; grid-template-columns: 100px minmax(0,1fr); align-items: end; gap: 1rem; }
.icon-field, .name-field { min-width: 0; display: flex; flex-direction: column; gap: .45rem; }
.field-label, .version-heading > label, .loader-field legend { color: var(--text-color); font-size: .79rem; font-weight: 700; }
.name-field input, .version-select-wrap select {
    box-sizing: border-box;
    width: 100%;
    height: 40px;
    border: 1px solid var(--gray-color);
    border-radius: 6px;
    padding: 0 .75rem;
    outline: 0;
    background: var(--surface-1);
    color: var(--text-color);
    font: inherit;
    font-size: .8rem;
}
.name-field input:focus, .version-select-wrap:focus-within { border-color: rgba(255,255,255,.38); }
.name-field input::placeholder { color: var(--text-dim-color); }
.version-heading { display: flex; align-items: center; justify-content: space-between; gap: 1rem; }
.version-toggle { display: inline-flex; align-items: center; gap: .5rem; color: var(--text-muted-color); font-size: .72rem; cursor: pointer; }
.version-toggle input { position: absolute; width: 1px; height: 1px; opacity: 0; }
.toggle-track { position: relative; width: 33px; height: 18px; border-radius: 10px; background: #414145; transition: background .15s ease; }
.toggle-track::after { position: absolute; top: 2px; left: 2px; width: 14px; height: 14px; border-radius: 50%; background: #fff; content: ""; transition: transform .15s ease; }
.version-toggle input:checked + .toggle-track { background: var(--accent-color); }
.version-toggle input:checked + .toggle-track::after { transform: translateX(15px); }
.version-select-wrap { position: relative; display: block; }
.version-select-wrap select { appearance: none; padding-right: 2rem; cursor: pointer; }
.version-select-wrap > i { position: absolute; top: 50%; right: .85rem; transform: translateY(-50%); color: var(--text-muted-color); font-size: .7rem; pointer-events: none; }
.version-select-wrap option { background: var(--surface-4); color: var(--text-color); }
.loader-field { min-width: 0; margin: 0; border: 0; padding: 0; }
.loader-field legend { margin-bottom: .65rem; }
.loader-options { display: flex; flex-wrap: wrap; gap: .45rem; }
.loader-option {
    min-height: 34px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: .45rem;
    border: 1px solid var(--gray-color);
    border-radius: 7px;
    padding: 0 .7rem;
    background: var(--surface-4);
    color: var(--text-muted-color);
    font: inherit;
    font-size: .75rem;
    font-weight: 650;
    cursor: pointer;
}
.loader-option:hover { border-color: #555; color: var(--text-color); }
.loader-option.active { border-color: var(--accent-color); background: rgba(255,68,68,.1); color: #ff7777; }
.loader-fabric.active { border-color: #b7a954; color: #d3c76e; background: rgba(183,169,84,.08); }
.loader-forge.active { border-color: #5793e2; color: #77aaf0; background: rgba(87,147,226,.08); }
.loader-neoforge.active { border-color: #d78a42; color: #e9a15c; background: rgba(215,138,66,.08); }
.loader-quilt.active { border-color: #a56be2; color: #bd8bf0; background: rgba(165,107,226,.08); }
.form-footer { display: flex; justify-content: flex-end; gap: .55rem; padding-top: .75rem; border-top: 1px solid rgba(255,255,255,.07); }
.cancel-button, .submit-button { min-height: 37px; display: inline-flex; align-items: center; justify-content: center; gap: .45rem; border-radius: 6px; padding: 0 .9rem; font: inherit; font-size: .78rem; font-weight: 700; cursor: pointer; }
.cancel-button { border: 1px solid var(--gray-color); background: var(--surface-4); color: var(--text-muted-color); }
.cancel-button:hover { color: var(--text-color); }
.submit-button { border: 1px solid #db514c; background: var(--accent-color); color: white; }
.submit-button:hover { background: #e84b46; }
@media (max-width: 540px) {
    .instance-form { padding: .9rem; }
    .profile-row { grid-template-columns: 76px minmax(0,1fr); gap: .65rem; }
    .version-heading { align-items: flex-start; flex-direction: column; gap: .45rem; }
}
@media (max-height: 500px) {
    .form-overlay { align-items: start; overflow: auto; }
    .instance-form { max-height: none; }
}
</style>
