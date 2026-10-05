// SPA Handler :)
const route = (event) => {
    event = event || window.event;
    event.preventDefault();
    window.history.pushState({}, "", event.target.href);
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
    const path = window.location.pathname;
    const route = routes[path] || routes[404];
    const html = await fetch(route).then((data) => data.text());
    document.getElementById("main-page").innerHTML = html;
};

window.onpopstate = handleLocation;
window.route = route;

handleLocation();
