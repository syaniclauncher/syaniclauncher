// SPA Handler :)
let currentRoute = "/";

const route = (event) => {
    event = event || window.event;
    event.preventDefault();
    const link = event.currentTarget || event.target;
    const path = link.getAttribute("href");
    if (!path || path === currentRoute) return;
    currentRoute = path;
    handleLocation();
};

const routes = {
    404: './pages/404.html',
    "/": './pages/home.html',
    "/instances": './pages/instances.html',
    "/public-servers": './pages/public-servers.html',
    "/performance": './pages/performance.html',
    "/settings": './pages/settings.html',
}

const handleLocation = async () => {
    const route = routes[currentRoute] || routes[404];
    const html = await fetch(route).then((data) => data.text());
    document.getElementById("main-page").innerHTML = html;
};

window.route = route;

handleLocation();
