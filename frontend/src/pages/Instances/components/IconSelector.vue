<script setup>
import { ref } from "vue";

const props = defineProps({
    icon: { type: String, default: "grass_block" },
    iconUrl: { type: String, default: "" },
});
const emit = defineEmits(["change"]);
const open = ref(false);

const iconPresets = [
    "acacia_leaves", "acacia_log", "acacia_planks", "ancient_debris", "andesite",
    "azalea_leaves", "bedrock", "birch_leaves", "birch_log", "birch_planks",
    "block_of_coal", "block_of_diamond", "block_of_emerald", "block_of_gold", "block_of_iron",
    "block_of_netherite", "block_of_redstone", "bookshelf", "bricks", "cake", "chest",
    "clay", "coal_ore", "cobblestone", "command_block", "copper_block", "copper_ore",
    "crafting_table", "creeper_head", "crimson_planks", "crimson_stem", "custom_package",
    "dark_oak_leaves", "dark_oak_log", "dark_oak_planks", "deepslate", "deepslate_coal_ore",
    "deepslate_copper_ore", "deepslate_diamond_ore", "deepslate_emerald_ore", "deepslate_gold_ore",
    "deepslate_iron_ore", "deepslate_lapis_lazuli_ore", "deepslate_redstone_ore", "diamond_ore",
    "diorite", "dirt", "emerald_ore", "enchanting_table", "end_stone", "ender_chest",
    "exposed_copper_block", "farmland", "flowering_azalea_leaves", "furnace", "glass",
    "glowstone", "gold_ore", "granite", "grass_block", "grass_path", "gravel", "honey_block",
    "ice", "iron_ore", "jungle_leaves", "jungle_log", "jungle_planks", "lapis_lazuli_ore",
    "lectern_book", "light_blue_glazed_terracotta", "lit_furnace", "mangrove_leaves",
    "mangrove_log", "mangrove_planks", "mycelium", "netherrack", "nether_bricks",
    "nether_quartz_ore", "nether_reactor_core", "nether_reactor_core_finished",
    "nether_reactor_core_initialized", "nether_wart_block", "oak_leaves", "oak_log",
    "oak_planks", "observer", "obsidian", "orange_glazed_terracotta", "oxidized_copper_block",
    "piston", "podzol", "pumpkin", "red_sand", "red_sandstone", "redstone_ore", "sand",
    "sandstone", "skeleton_skull", "slime_block", "snow_block", "snowy_grass_block",
    "soul_sand", "spruce_leaves", "spruce_log", "spruce_planks", "sticky_piston", "stone",
    "terracotta", "tnt", "tuff", "warped_planks", "warped_stem", "warped_wart_block",
    "water", "weathered_copper_block", "white_glazed_terracotta", "white_wool", "xmas_chest",
];

function selectedUrl() {
    return props.iconUrl || `/assets/instance_icons/${props.icon}.png`;
}

function choosePreset(icon) {
    emit("change", { icon, iconUrl: "" });
    open.value = false;
}

function uploadIcon(event) {
    const file = event.target.files?.[0];
    if (!file || !file.type.startsWith("image/")) return;
    const reader = new FileReader();
    reader.onload = () => emit("change", { icon: "__custom__", iconUrl: reader.result });
    reader.readAsDataURL(file);
}
</script>

<template>
    <div class="installation-icon-section" @click.stop>
        <button type="button" class="installation-icon-btn" aria-label="Choose installation icon" @click="open = !open">
            <img :src="selectedUrl()" alt="Icon">
            <span class="installation-icon-edit"><i class="fas fa-pen"></i></span>
        </button>

        <div v-if="open" class="icon-picker-dropdown">
            <div class="icon-picker-header">
                <div>
                    <div class="icon-picker-title">Choose Icon</div>
                    <div class="icon-picker-copy">Minecraft blocks and launcher presets</div>
                </div>
                <div class="icon-picker-preview"><img :src="selectedUrl()" alt=""></div>
            </div>
            <label class="icon-upload-btn">
                <i class="fas fa-upload" aria-hidden="true"></i>
                <span>Use local image</span>
                <input type="file" accept="image/*" @change="uploadIcon">
            </label>
            <div class="icon-picker-grid">
                <button v-for="preset in iconPresets" :key="preset" type="button" class="icon-option"
                    :class="{ 'is-selected': icon === preset && !iconUrl }"
                    :title="preset.replaceAll('_', ' ')"
                    @click="choosePreset(preset)">
                    <img :src="`/assets/instance_icons/${preset}.png`" :alt="preset">
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.installation-icon-section { position: relative; }
.installation-icon-btn {
    position: relative;
    width: 64px;
    height: 64px;
    padding: 8px;
    border: 1px solid rgba(255,255,255,.1);
    border-radius: 8px;
    background: rgba(0,0,0,.2);
    cursor: pointer;
}
.installation-icon-btn:hover { border-color: rgba(255,68,68,.5); background: rgba(255,255,255,.06); }
.installation-icon-btn img { width: 100%; height: 100%; object-fit: contain; image-rendering: pixelated; }
.installation-icon-edit {
    position: absolute;
    right: -5px;
    bottom: -5px;
    width: 22px;
    height: 22px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(255,255,255,.12);
    border-radius: 6px;
    background: var(--surface-5);
    color: #ff7777;
    font-size: .68rem;
}
.icon-picker-dropdown {
    position: absolute;
    z-index: 1000;
    left: 0;
    top: 70px;
    box-sizing: border-box;
    width: min(430px, calc(100vw - 48px));
    max-height: min(390px, calc(100vh - 220px));
    display: flex;
    flex-direction: column;
    gap: .75rem;
    overflow: hidden;
    padding: .85rem;
    border-radius: 8px;
    background: var(--surface-5);
    box-shadow: 0 8px 24px rgba(0,0,0,.5);
    animation: reveal .2s ease;
}
.icon-picker-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding-bottom: .7rem;
    border-bottom: 1px solid rgba(255,255,255,.08);
}
.icon-picker-title { color: var(--text-color); font-size: .9rem; font-weight: 700; }
.icon-picker-copy { margin-top: .15rem; color: rgba(255,255,255,.45); font-size: .75rem; }
.icon-picker-preview {
    width: 38px;
    height: 38px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 6px;
    border: 1px solid rgba(255,255,255,.08);
    border-radius: 8px;
    background: rgba(255,255,255,.04);
}
.icon-picker-preview img { width: 100%; height: 100%; object-fit: contain; image-rendering: pixelated; }
.icon-upload-btn {
    position: relative;
    min-height: 34px;
    display: flex;
    align-items: center;
    gap: .55rem;
    padding: 0 .7rem;
    border: 1px solid var(--gray-color);
    border-radius: 4px;
    color: var(--text-color);
    font-size: .78rem;
    font-weight: 700;
    cursor: pointer;
}
.icon-upload-btn input { position: absolute; inset: 0; width: 100%; opacity: 0; cursor: pointer; }
.icon-picker-grid {
    min-height: 0;
    display: grid;
    grid-template-columns: repeat(8,minmax(0,1fr));
    gap: .42rem;
    overflow-y: auto;
    padding-right: .2rem;
}
.icon-option {
    aspect-ratio: 1;
    padding: 6px;
    border: 1px solid transparent;
    border-radius: 6px;
    background: rgba(255,255,255,.035);
    cursor: pointer;
}
.icon-option:hover { border-color: rgba(255,255,255,.18); background: rgba(255,255,255,.08); }
.icon-option.is-selected { border-color: var(--accent-color); background: rgba(255,68,68,.12); }
.icon-option img { width: 100%; height: 100%; object-fit: contain; image-rendering: pixelated; }
@keyframes reveal { from { opacity: 0; transform: translateY(3px); } to { opacity: 1; transform: translateY(0); } }
</style>
