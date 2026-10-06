import { renderHomeNews, stopHomeNews } from "./homePage/handlerNews.js";

let currentRoute = "/";
let activePage;
let navigationId = 0;

const pages = {
    "/": {
        template: "./pages/home.html",
        onEnter: renderHomeNews,
        onLeave: stopHomeNews,
    },
    "/instances": {
        template: "./pages/instances.html",
    },
    "/public-servers": {
        template: "./pages/public-servers.html",
    },
    "/performance": {
        template: "./pages/performance.html",
    },
    "/settings": {
        template: "./pages/settings.html",
    },
    404: {
        template: "./pages/404.html",
    },
};

const handleLocation = async () => {
    const navigation = ++navigationId;
    const page = pages[currentRoute] || pages[404];
    const response = await fetch(page.template);

    if (!response.ok) {
        throw new Error(`Failed to load page ${page.template}: HTTP ${response.status}`);
    }

    const content = await response.text();
    if (navigation !== navigationId) return; // ignore old navigation if another one started
    activePage?.onLeave?.();
    document.getElementById("main-page").innerHTML = content;
    activePage = page;
    activePage.onEnter?.();
};

document.addEventListener("click", (event) => {
    const link = event.target.closest("a");
    if (!link) return;
    const path = link.getAttribute("href");
    if (!path || !path.startsWith("/")) return; // ignore external links
    event.preventDefault();
    if (path === currentRoute) return;
    currentRoute = path;
    handleLocation();
});

handleLocation();
