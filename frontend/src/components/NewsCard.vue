
<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";

const FEED = "https://syaniclauncher.github.io/news/minecraft_news.json";

let cachedNews;
async function fetchNews() {
    try {
        return cachedNews ??= await fetch(FEED).then(response => {
            if (!response.ok) throw new Error(response.status);
            return response.json();
        });
    } catch {
        return null;
    }
}

const news = ref(); // undefined = loading, null = failed
const hovering = ref(false);
const index = ref(0);
let timer;

const item = computed(() => news.value?.[index.value]);
const date = computed(() => item.value?.date && new Date(item.value.date).toLocaleDateString());

const step = amount => {
    const length = news.value?.length;
    if (length) index.value = (index.value + amount + length) % length;
};

onMounted(async () => {
    news.value = await fetchNews();

    if (news.value?.length > 1)
        timer = setInterval(() => {
            if (!hovering.value) step(1);
        }, 7000);
});

onUnmounted(() => clearInterval(timer));
</script>

<template>
    <div v-if="news === undefined" class="home-news-empty">Loading Minecraft news...</div>
    <div v-else-if="news === null" class="home-news-empty">Could not load Minecraft news.</div>
    <div v-else-if="!news.length" class="home-news-empty">No Minecraft news available.</div>

    <article v-else class="news-card" @mouseenter="hovering = true" @mouseleave="hovering = false">
        <div class="news-card-media">
            <img v-if="item.image" class="news-card-image" :src="item.image" alt="">
        </div>
        <div class="news-card-body">
            <div class="news-card-meta">
                <template v-if="item.tag">
                    <span class="news-card-tag">{{ item.tag }}</span>
                    <span class="news-card-sep">•</span>
                </template>
                <span class="news-card-date">{{ date }}</span>
            </div>
            <h3 class="news-card-title">{{ item.title }}</h3>
            <p class="news-card-text">{{ item.text || "" }}</p>
            <button v-if="item.url" class="news-card-more" type="button" @click="BrowserOpenURL(item.url)">Read More..</button>
        </div>

        <template v-if="news.length > 1">
            <button class="news-nav news-nav--prev" type="button" aria-label="Previous article" @click="step(-1)">‹</button>
            <button class="news-nav news-nav--next" type="button" aria-label="Next article" @click="step(1)">›</button>
            <div class="news-dots">
                <button v-for="(entry, i) in news" :key="entry.title" class="news-dot"
                    :class="{ 'is-active': i === index }"
                    type="button" :aria-label="`Article ${i + 1}`" @click="index = i"></button>
            </div>
        </template>
    </article>
</template>

<style scoped>
.home-news-section { margin-top: 1.5rem; }
.home-news-section h2 { margin: 0 0 0.85rem; font-size: 1.1rem; }
.news-card { position: relative; display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1.1fr); min-height: 200px; border-radius: 10px; background: var(--surface-3); overflow: hidden; }
.news-card-media { position: relative; min-height: 200px; background: var(--surface-3); }
.news-card-image { position: absolute; inset: 0; width: 100%; height: 100%; object-fit: cover; mask-image: linear-gradient(to right, var(--image-mask-color) 70%, transparent 100%); }
.news-card-body { display: flex; flex-direction: column; justify-content: center; gap: 0.5rem; min-width: 0; padding: 1.25rem 3.25rem 2.25rem 0.5rem; }
.news-card-meta { display: flex; align-items: center; gap: 0.5rem; }
.news-card-tag { color: var(--text-muted-color); font-size: 0.7rem; font-weight: 700; text-transform: uppercase; }
.news-card-sep { color: var(--text-dim-color); font-size: 0.7rem; }
.news-card-date { color: var(--text-subtle-color); font-size: 0.75rem; }
.news-card-title { margin: 0; font-size: 1.3rem; }
.news-card-text { margin: 0; color: rgba(var(--text-color), 0.58); font-size: 0.875rem; line-height: 1.5; text-indent: 0; opacity: 0.58; }
.news-card-more { align-self: flex-start; display: inline-flex; align-items: center; gap: 0.4rem; margin-top: 0.15rem; padding: 0; border: 0; background: transparent; color: var(--accent-color); font: inherit; font-size: 0.85rem; font-weight: 600; cursor: pointer; }
.news-card-more:hover { text-decoration: underline; }
.news-nav { position: absolute; top: 50%; z-index: 1; display: inline-flex; align-items: center; justify-content: center; width: 30px; height: 30px; padding: 0; transform: translateY(-50%); border: 0; border-radius: 50%; background: var(--control-overlay-color); color: var(--text-color); font-size: 1.2rem; cursor: pointer; opacity: 0; transition: opacity 0.15s ease, background 0.15s ease; }
.news-card:hover .news-nav, .news-nav:focus-visible { opacity: 1; }
.news-nav:hover { background: var(--surface-hover-color); }
.news-nav--prev { left: 0.6rem; }
.news-nav--next { right: 0.6rem; }
.news-dots { position: absolute; right: 1.5rem; bottom: 0.85rem; display: flex; align-items: center; gap: 0.4rem; }
.news-dot { width: 7px; height: 7px; padding: 0; border: 0; border-radius: 50%; background: var(--text-disabled-color); cursor: pointer; transition: background 0.15s ease, width 0.15s ease; }
.news-dot.is-active { width: 16px; border-radius: 999px; background: var(--accent-color); }
.home-news-empty { display: flex; align-items: center; justify-content: center; min-height: 200px; border: 1px solid var(--gray-color); border-radius: 10px; background: var(--surface-alt-color); color: var(--text-subtle-color); text-align: center; }
</style>
