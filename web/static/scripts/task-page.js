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


const socket = new WebSocket("ws://127.0.0.1:8080/api/ws");

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
    document.getElementById('dropdownMenuButton').innerText = item;
}

function sendRequest() {
    // construct the API URL
    const apiUrl = window.location.href + '/submit'

    // send a POST request with JSON data
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
