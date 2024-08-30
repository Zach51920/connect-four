function getPreferredTheme() {
    if (localStorage.getItem("theme")) {
        return localStorage.getItem("theme");
    }
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "business" : "corperate";
}

function setTheme(theme) {
    document.documentElement.setAttribute("data-theme", theme);
    document.body.setAttribute("data-theme", theme);
    localStorage.setItem("theme", theme);
}

function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute("data-theme");
    const newTheme = currentTheme === "corporate" ? "business" : "corporate";
    setTheme(newTheme);
}

// Set the initial theme when the page loads
document.addEventListener("DOMContentLoaded", function () {
    setTheme(getPreferredTheme());

    // Update the checkbox state based on the current theme
    const themeToggle = document.querySelector(".theme-controller");
    if (themeToggle) {
        themeToggle.checked = getPreferredTheme() === "corporate";
    }
});

// Listen for changes in color scheme preference
window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", (e) => {
    if (!localStorage.getItem("theme")) {
        setTheme(e.matches ? "business" : "corporate");
    }
});