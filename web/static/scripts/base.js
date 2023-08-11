function toggleSidebar() {
    const sidebar = document.getElementById('sidebar');
    const content = document.getElementById('content');
    sidebar.classList.toggle('collapsed');
    content.classList.toggle('collapsed');
}

function toggleTheme() {
    const themeToggle = document.getElementById('theme-toggle');
    const moonIcon = themeToggle.querySelector('.bi-moon-fill');
    const sunIcon = themeToggle.querySelector('.bi-sun-fill');
    const sidebarLabel = themeToggle.querySelector('.sidebar-label');


    if (sunIcon.classList.contains('d-none')) {
        moonIcon.classList.add('d-none');
        sunIcon.classList.remove('d-none');
        sidebarLabel.innerText = "День";
    } else {
        sunIcon.classList.add('d-none');
        moonIcon.classList.remove('d-none');
        sidebarLabel.innerText = "Ночь";
    }
}
