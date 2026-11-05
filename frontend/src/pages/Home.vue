<script setup>
import NewsCard from "../components/NewsCard.vue";
import LaunchBar from "../components/LaunchBar.vue";

defineEmits(["navigate"]);
</script>

<template>
    <div class="page-container">
        <section class="home-news-section">
            <h2>Latest Minecraft News</h2>
            <NewsCard />
        </section>
    </div>
    <LaunchBar @navigate="$emit('navigate', $event)" />
</template>
