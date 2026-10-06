import { BrowserOpenURL } from "../wailsjs/runtime/runtime.js";

let news;
let newsIndex = 0;
let newsTimer;

// why news loader is in the frontend? because it keep things simple
// don't argue with me, it works and it's not a big deal, just let it be
async function loadNews() {
    if (news !== undefined) return news;
    try { 
        news = await fetch("https://syaniclauncher.github.io/news/minecraft_news.json").then(response => response.json());
    } catch (error) { 
        news = null; // null = failed to load & undefined = not loaded yet btw
    }
    return news;
}

export function stopHomeNews() {
    clearInterval(newsTimer); newsTimer = undefined; // stop newsTimer
    const newsContainer = document.getElementById("homeNews");
    if (newsContainer) newsContainer.onclick = null; // remove event listener when leaving the page
}

// escape HTML special characters, move to utils later
function escapeHTML(value = "") {
    return String(value).replace(/[&<>"']/g, char => ({
        "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;",
    })[char]);
}

export async function renderHomeNews() {
    clearInterval(newsTimer);
    const newsContainer = document.getElementById("homeNews");
    if (!newsContainer) return;

    if (!news) newsContainer.innerHTML = '<div class="home-news-empty">Loading Minecraft news...</div>';
    try {
        await loadNews();
    } catch (error) {
        if (!newsContainer.isConnected) return; //checks if newscontainer is loaded in page
        console.error("Failed to fetch news:", error);
        newsContainer.innerHTML = '<div class="home-news-empty">Could not load Minecraft news.</div>';
        return;
    }

    if (!newsContainer.isConnected) return;
    if (!Array.isArray(news) || !news.length) {
        newsContainer.innerHTML = '<div class="home-news-empty">No Minecraft news available.</div>';
        return;
    }

    function render() {
        const item = news[newsIndex];
        const date = item.date
            ? new Date(item.date).toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" })
            : "";
        newsContainer.innerHTML = `
            <article class="news-card">
                <div class="news-card-media">${item.image ? `<img class="news-card-image" src="${escapeHTML(item.image)}" alt="">` : ""}</div>
                <div class="news-card-body">
                    <div class="news-card-meta">${item.tag ? `<span class="news-card-tag">${escapeHTML(item.tag)}</span><span class="news-card-sep">•</span>` : ""}
                        <span class="news-card-date">${escapeHTML(date)}</span>
                    </div>
                    <h3 class="news-card-title">${escapeHTML(item.title)}</h3>
                    <p class="news-card-text">${escapeHTML(item.text || "")}</p>
                    ${item.url ? `<button class="news-card-more" data-news-url="${escapeHTML(item.url)}" type="button">Read More..</button>` : ""}
                </div>
                ${news.length > 1 ? `
                    <button class="news-nav news-nav--prev" data-news-step="-1" aria-label="Previous article" type="button">‹</button>
                    <button class="news-nav news-nav--next" data-news-step="1" aria-label="Next article" type="button">›</button>
                    <div class="news-dots">${news.map((_, i) => `<button class="news-dot${i === newsIndex ? " is-active" : ""}" data-news-index="${i}" aria-label="Article ${i + 1}" type="button"></button>`).join("")}</div>
                ` : ""}
            </article>`;
    }

    newsContainer.onclick = event => {
        const button = event.target.closest("button");
        if (!button) return;
        if (button.dataset.newsUrl) return BrowserOpenURL(button.dataset.newsUrl);
        if (button.dataset.newsStep) newsIndex = (newsIndex + Number(button.dataset.newsStep) + news.length) % news.length;
        if (button.dataset.newsIndex !== undefined) newsIndex = Number(button.dataset.newsIndex);
        render();
    };

    render();
    if (news.length > 1) {
        newsTimer = setInterval(() => {
            if (!newsContainer.isConnected) {
                stopHomeNews(); 
                return;
            }

            if (newsContainer.querySelector(".news-card:hover")) return;
            newsIndex = (newsIndex + 1) % news.length;
            render();
        }, 7000); // 7s
    }
}
