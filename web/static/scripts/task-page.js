codingTheme = "default";

if (body.classList.contains("dark-mode")) {
    codingTheme = "material";
}

// Инициализация CodeMirror для элемента с id="Code"
const codeTextArea = document.getElementById("Code");
const codeMirrorEditor = CodeMirror.fromTextArea(codeTextArea, {
    mode: "text/plain",
    lineNumbers: true,
    theme: codingTheme,
});
codeMirrorEditor.setSize(null, 500);

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

    // Получаем класс текущей первой строки таблицы
    const firstRowClass = table.querySelector("tr:first-child")?.classList[0];

    // Определяем класс для новой строки в зависимости от класса текущей первой строки
    const newRowClass = firstRowClass === "even" ? "odd" : "even";

    newRow.className = newRowClass;

    newRow.innerHTML = `
        <td>${data.id}</td>
        <td>${data.language}</td>
        <td><span class="badge ${verdictClass}">${data.verdict}</span></td>
        <td>${data.test}</td>
        <td>${data.execution_time}</td>
        <td>${data.memory_used}</td>
        <td><a href="/view/submission/${data.id}">тык</a></td>
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
        return
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

document.addEventListener('DOMContentLoaded', function() {
    const copyIcons = document.querySelectorAll('.copy-icon');
    copyIcons.forEach(icon => {
        icon.addEventListener('click', function() {
            const parentTestBox = this.closest('.test-box');
            const textToCopyElement = parentTestBox.querySelector('.test-input, .test-output');
            const textToCopy = textToCopyElement.innerHTML.replace(/<br\s*[\/]?>/gi, '\n'); // Replace <br> tags with newline characters
            const textarea = document.createElement('textarea');
            textarea.value = textToCopy;
            document.body.appendChild(textarea);
            textarea.select();
            document.execCommand('copy');
            document.body.removeChild(textarea);

            // Show Toastify notification
            Toastify({
                text: 'Text copied to clipboard',
                duration: 2000, // Notification will disappear after 2 seconds
                gravity: 'bottom', // Position the notification at the bottom
                position: 'right' // Position the notification on the right side
            }).showToast();
        });
    });
});
