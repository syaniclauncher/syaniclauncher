let news = null;

async function fetchNews() {
    try {
        const response = await fetch("https://syaniclauncher.github.io/news/minecraft_news.json");
        if (!response.ok) { throw new Error(`HTTP ${response.status}`); }
        news = await response.json();
    } catch (error) {
        console.error("Failed to fetch news:", error);
        news = null;
    }
}
