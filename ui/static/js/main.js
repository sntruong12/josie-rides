var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
    var link = navLinks[i]
    if (link.getAttribute('href') == window.location.pathname) {
        link.classList.add("live");
        break;
    }
}

// timezone handling in create ride
document.addEventListener("DOMContentLoaded", () => {
    const tzSelect = document.getElementById("create_ride_timezone");
    if (!tzSelect) {
        return;
    }

    const userTz = Intl.DateTimeFormat().resolvedOptions().timeZone || "UTC";

    if (typeof Intl.supportedValuesOf === "function") {
        const timezones = Intl.supportedValuesOf("timeZone");
        tzSelect.innerHTML = "";

        const fragment = document.createDocumentFragment();
        timezones.forEach(tz => {
            const option = document.createElement("option");
            option.value = tz;
            option.textContent = tz.replace(/_/g, " ");
            if (tz === userTz) {
                option.selected = true;
            }
            fragment.appendChild(option);
        });

        tzSelect.appendChild(fragment);
    } else if (userTz !== "UTC") {
        const option = document.createElement("option");
        option.value = userTz;
        option.textContent = userTz.replace(/_/g, " ");
        option.selected = true;
        tzSelect.appendChild(option);
    }
});