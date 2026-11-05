<script setup>
import { ref } from "vue";
import Home from "@/pages/Home/Home.vue";
import Instances from "@/pages/Instances.vue";
import PublicServers from "@/pages/PublicServers.vue";
import Performance from "@/pages/Performance.vue";
import Settings from "@/pages/Settings.vue";

const routes = {
    "/": Home,
    "/instances": Instances,
    "/public-servers": PublicServers,
    "/performance": Performance,
    "/settings": Settings,
};

const nav = [
    { path: "/", label: "Play Game" },
    { path: "/instances", label: "Instances" },
    { path: "/public-servers", label: "Public Servers" },
    { path: "/performance", label: "Performance" },
    { path: "/settings", label: "Settings" },
];

const route = ref("/");
</script>

<template>
    <div class="app">
        <nav class="sidebar main-nav">
        <a v-for="item in nav" :key="item.path" href="#" :class="{ active: route === item.path }"
            @click.prevent="route = item.path">{{ item.label }}</a>
        </nav>
        <main>
            <component :is="routes[route] || Home" @navigate="route = $event" />
        </main>
    </div>
</template>

<style scoped>
:global(html, body, #app) { height: 100%; margin: 0; }

.app { display: flex; height: 100vh; width: 100%; }
main { display: flex; flex-direction: column; flex-grow: 1; min-width: 0; min-height: 0; }
a.active { background-color: var(--gray-color); }

.main-nav { display: flex; flex-direction: column; flex-shrink: 0; align-items: flex-start;}
.main-nav a { color: var(--text-color); display: block; text-decoration: none; padding: 0.3rem 0.5rem; border-radius: 8px;}
.main-nav a:hover { background-color: var(--gray-color); }

.sidebar { display: flex; flex-direction: column; padding: 1rem; border-right: 1px solid var(--gray-color); background-color: var(--surface-8);}
.sidebar > a { margin: 0.05rem 0; }
</style>
