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

function updateTable(data) {
    const table = document.querySelector(".table tbody");
    const newRow = document.createElement("tr");

    newRow.setAttribute("data-id", data.id);

    let verdictClass = "";
    if (data.verdict === "Accepted") {
        verdictClass = "badge-success";
    } else if (data.verdict === "Waiting") {
        verdictClass = "badge-warning";
    } else if (data.verdict === "Compiling") {
        verdictClass = "badge-info";
    } else {
        verdictClass = "badge-danger";
    }

    newRow.innerHTML = `
        <td>${data.id}</td>
        <td>${data.language}</td>
        <td><span class="badge ${verdictClass}">${data.verdict}</span></td>
        <td>${data.test}</td>
        <td>${data.execution_time}</td>
        <td>${data.memory_used}</td>
    `;

    // Проверка наличия строки с таким id
    const existingRow = table.querySelector(`tr[data-id="${data.id}"]`);
    if (existingRow) {
        existingRow.innerHTML = newRow.innerHTML;
    } else if ((table.querySelector("tr:first-child td:first-child") === null) || (data.id > parseInt(table.querySelector("tr:first-child td:first-child").innerText))) {
        table.insertBefore(newRow, table.firstChild);
    }
}

const url = window.location.href;
const parts = url.split("/");
const host = window.location.host;

// Поиск индексов contest_id и problem_id в массиве parts
const contestIndex = parts.indexOf("contest");
const problemIndex = parts.indexOf("problem");

let contestId = -1;
const problemId = parts[problemIndex + 1];
if (contestIndex !== -1 && problemIndex !== -1) {
    contestId = parts[contestIndex + 1];
}

const socket = new WebSocket(`ws://${host}/api/ws/contest/${contestId}/problem/${problemId}`);


socket.onmessage = function(event) {
    const message = JSON.parse(event.data);
    // console.log(message);

    updateTable(message);
};

socket.onclose = function(event) {
    console.log(event)
};

socket.onerror = function(event) {
    console.log(event)
};

function selectItem(item) {
    let selectedMode = "text/plain"; // Режим по умолчанию для неизвестных языков

    const forChecking = item.toLowerCase()
    if (forChecking.includes("c++") || forChecking.includes("g++") | forChecking.includes("clang")) {
        selectedMode = "text/x-c++src"
    } else if (forChecking.includes("python") || forChecking.includes("pypy")) {
        selectedMode = "text/x-python"
    } else if (forChecking.includes("java")) {
        selectedMode = "text/x-java"
    }
    else if (forChecking.includes("go")) {
        selectedMode = "text/x-go"
    } else if (forChecking.includes("kotlin")) {
        selectedMode = "text/x-kotlin"
    } else if (forChecking.includes("java")) {
        selectedMode = "text/x-java"
    } else if (forChecking.includes("visual c")) {
        selectedMode = "text/x-csrc"
    } else if (forChecking.includes("pascal")) {
        selectedMode = "text/x-pascal"
    }

    codeMirrorEditor.setOption("mode", selectedMode);

    document.getElementById('dropdownMenuButton').innerText = item;
}

function sendRequest() {
    const apiUrl = window.location.href + '/submit'

    if (document.getElementById('dropdownMenuButton').innerText === 'Select an language ') {
        alert("Выбери язык программирования!")
    }

    fetch(apiUrl, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "sourceCode": codeMirrorEditor.getValue(),
            "language": document.getElementById('dropdownMenuButton').innerText,
        })
    })
        .then(response => {
            if (response.status === 200) {
                return response.json();
            } else {
                throw new Error('Request failed');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}
