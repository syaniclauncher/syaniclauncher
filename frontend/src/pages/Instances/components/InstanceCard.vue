<script setup>
const props = defineProps({
    instance: { type: Object, required: true },
    layout: { type: String, default: "grid" },
});
const emit = defineEmits(["open", "play"]);

function iconUrl() {
    return props.instance.iconUrl || `/assets/instance_icons/${props.instance.icon || "grass_block"}.png`;
}

function loaderLabel() {
    const labels = { vanilla: "Vanilla", forge: "Forge", fabric: "Fabric", quilt: "Quilt", neoforge: "NeoForge" };
    return labels[(props.instance.loader || "vanilla").toLowerCase()] || "Vanilla";
}

function lastPlayedLabel() {
    if (!props.instance.lastPlayed) return "Never played";
    const days = Math.floor((Date.now() / 1000 - props.instance.lastPlayed) / 86400);
    if (days <= 0) return "Today";
    if (days === 1) return "Yesterday";
    return `${days} days ago`;
}

function listLastPlayedLabel() {
    const label = lastPlayedLabel();
    return label === "Never played" ? label : `Played ${label.toLowerCase()}`;
}

</script>

<template>
    <article class="instance-card" :class="`view-${layout}`">
        <button type="button" class="instance-card-main" @click="emit('open', instance)">
            <div class="instance-card-top">
                <span class="instance-icon"><img :src="iconUrl()" :alt="instance.name"></span>
                <span class="last-played"><i class="far fa-calendar" aria-hidden="true"></i> {{ lastPlayedLabel() }}</span>
            </div>
            <div class="instance-name">{{ instance.name }}</div>
            <div class="instance-version">{{ loaderLabel() }} {{ instance.version || "Version not set" }}</div>
        </button>
        <div class="card-actions">
            <span v-if="layout === 'list'" class="list-last-played"
                :class="{ 'never-played': !instance.lastPlayed }">{{ listLastPlayedLabel() }}</span>
            <button type="button" class="card-play-button" @click="emit('play', instance)">
                <i class="fas fa-play" aria-hidden="true"></i>
                <span>Play</span>
            </button>
            <button type="button" class="card-action-button folder-action" title="Opening the instance folder requires backend support" disabled>
                <i class="fas fa-folder-open" aria-hidden="true"></i>
            </button>
            <button type="button" class="card-action-button details-action" title="Instance details and settings" @click="emit('open', instance)">
                <i class="fas fa-cog" aria-hidden="true"></i>
            </button>
        </div>
    </article>
</template>

<style scoped>
.instance-card {
    position: relative;
    box-sizing: border-box;
    min-width: 0;
    min-height: 168px;
    display: flex;
    flex-direction: column;
    border-radius: 9px;
    background: #101010;
    box-shadow: inset 0 1px rgba(255,255,255,.025);
}
.instance-card:hover { background: #131313; }
.instance-card-main {
    min-width: 0;
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: stretch;
    border: 0;
    padding: .75rem .8rem .65rem;
    background: transparent;
    color: var(--text-color);
    text-align: left;
    font: inherit;
    cursor: pointer;
}
.instance-card-main:focus-visible { outline: 1px solid var(--accent-color); outline-offset: -2px; border-radius: inherit; }
.instance-card-top { display: flex; align-items: flex-start; justify-content: space-between; gap: .5rem; }
.instance-icon {
    width: 42px;
    height: 42px;
    display: grid;
    place-items: center;
    border-radius: 10px;
    background: var(--surface-5);
}
.instance-icon img { width: 30px; height: 30px; object-fit: contain; image-rendering: pixelated; }
.last-played { display: flex; align-items: center; gap: .35rem; padding-top: .18rem; color: var(--text-dim-color); font-size: .72rem; white-space: nowrap; }
.instance-name { overflow: hidden; margin-top: .45rem; font-size: .9rem; font-weight: 750; text-overflow: ellipsis; white-space: nowrap; }
.instance-version { overflow: hidden; padding-top: .15rem; color: var(--text-dim-color); font-size: .68rem; text-overflow: ellipsis; white-space: nowrap; }
.card-actions { display: flex; align-items: center; gap: .25rem; padding: 0 .55rem .55rem; }
.card-play-button { min-width: 66px; height: 34px; display: inline-flex; align-items: center; justify-content: center; gap: .3rem; margin-right: auto; border: 1px solid var(--accent-color); border-radius: 6px; padding: 0 .45rem; background: var(--accent-color); color: #fff; font: inherit; font-size: .76rem; font-weight: 700; cursor: pointer; }
.card-play-button:hover { background: #e84b46; }
.card-action-button { width: 27px; height: 27px; border: 1px solid transparent; border-radius: 5px; background: var(--surface-5); color: var(--text-muted-color); cursor: pointer; }
.card-action-button:hover { border-color: var(--gray-color); color: var(--text-color); }
.view-list { min-height: 67px; display: grid; grid-template-columns: minmax(0,1fr) auto; align-items: center; border: 0; border-bottom: 1px solid var(--gray-color); border-radius: 0; background: transparent; box-shadow: none; }
.view-list .instance-card-main { display: grid; grid-template-columns: 42px minmax(0,1fr); grid-template-rows: auto auto; align-items: center; column-gap: .7rem; padding: .45rem .7rem; }
.view-list .instance-card-top { grid-column: 1; grid-row: 1 / span 2; align-items: center; justify-content: flex-start; }
.view-list .last-played { display: none; }
.view-list .instance-icon { width: 36px; height: 36px; background: transparent; }
.view-list .instance-icon img { width: 34px; height: 34px; }
.view-list .instance-name { grid-column: 2; grid-row: 1; align-self: end; margin: 0; font-size: .9rem; }
.view-list .instance-version { grid-column: 2; grid-row: 2; align-self: start; margin: 0; padding: 0; }
.view-list .card-actions { align-self: center; gap: .6rem; padding: .4rem .85rem .4rem 0; }
.list-last-played { color: var(--text-dim-color); font-size: .68rem; white-space: nowrap; }
.view-list .card-action-button { width: 32px; height: 32px; border: 1px solid var(--gray-color); border-radius: 6px; background: transparent; }
.folder-action:disabled { opacity: .45; cursor: not-allowed; }
@media (max-width: 640px) {
    .view-list .instance-card-main { grid-template-columns: 42px minmax(0,1fr); column-gap: .55rem; }
    .view-list .card-actions { gap: .35rem; padding-right: .5rem; }
    .view-list .card-action-button { width: 29px; height: 29px; }
}
@media (hover: hover) {
    .card-play-button { visibility: hidden; opacity: 0; transition: opacity .12s ease; }
    .instance-card:hover .card-play-button,
    .instance-card:focus-within .card-play-button { visibility: visible; opacity: 1; }
    .view-list .never-played { visibility: hidden; opacity: 0; transition: opacity .12s ease; }
    .view-list:hover .never-played,
    .view-list:focus-within .never-played { visibility: visible; opacity: 1; }
}
</style>
