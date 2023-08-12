const themeToggle = document.getElementById('theme-toggle');
const sidebarLabel = themeToggle.querySelector('.sidebar-label');
const moonIcon = themeToggle.querySelector('.bi-moon-fill');
const sunIcon = themeToggle.querySelector('.bi-sun-fill');

function toggleSidebar() {
    const sidebar = document.getElementById('sidebar');
    const content = document.getElementById('content');
    sidebar.classList.toggle('collapsed');
    content.classList.toggle('collapsed');
}

function toggleTheme() {
    body.classList.toggle('dark-mode');
    if (sunIcon.classList.contains('d-none')) {
        moonIcon.classList.add('d-none');
        sunIcon.classList.remove('d-none');
        sidebarLabel.innerText = "День";

        if (document.getElementById("Code") !== null) {
            codeMirrorEditor.setOption("theme", "material");
        }

        localStorage.setItem('theme', 'dark');
    } else {
        sunIcon.classList.add('d-none');
        moonIcon.classList.remove('d-none');
        sidebarLabel.innerText = "Ночь";

        if (document.getElementById("Code") !== null) {
            codeMirrorEditor.setOption("theme", "default");
        }

        localStorage.setItem('theme', 'light');
    }
}

const body = document.body;
const savedTheme = localStorage.getItem('theme');

// Добавляем класс, который блокирует анимацию перехода


if (savedTheme === 'dark') {
    moonIcon.classList.add('d-none');
    sunIcon.classList.remove('d-none');
    body.classList.add('dark-mode');
}

// Удаляем класс через небольшую паузу
setTimeout(() => {
    body.classList.add('set-transition');
    const sidebar = document.getElementById("sidebar");
    sidebar.classList.add('set-transition');
    }, 100);
