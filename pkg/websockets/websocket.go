package websockets

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем все Origin
	},
}

// Структура для хранения соединений пользователя и текущей страницы
type UserConnections struct {
	Connection       *websocket.Conn
	ProblemId        string
	ContestId        int
	ContestProblemId int
}

var Connections = make(map[string][]UserConnections) // Мапа для хранения соединений по юзернейму
var Mu sync.Mutex                                    // Мьютекс для безопасного доступа к мапе Connections

// Функция для отправки объекта Submission по веб-сокету по нику и странице
func SendMessageToUser(username string, submissionFull database.Submission) {
	Mu.Lock()
	userConns, ok := Connections[username]
	Mu.Unlock()

	if !ok {
		return
	}

	type SendingDataSubmission struct {
		ID            int    `json:"id"`
		Verdict       string `json:"verdict"`
		Language      string `json:"language"`
		ExecutionTime string `json:"execution_time"`
		MemoryUsed    string `json:"memory_used"`
		Test          string `json:"test"`
	}

	submission := SendingDataSubmission{
		ID:            submissionFull.ID,
		Verdict:       submissionFull.Verdict,
		Language:      submissionFull.Language,
		ExecutionTime: submissionFull.ExecutionTime,
		MemoryUsed:    submissionFull.MemoryUsed,
		Test:          submissionFull.Test,
	}

	// Преобразуем объект Submission в JSON
	message, err := json.Marshal(submission)
	if err != nil {
		log.Printf("Ошибка при преобразовании Submission в JSON: %v", err)
		return
	}

	// Отправляем сообщение только на нужной странице
	for _, conn := range userConns {
		if (conn.ContestId != -1 && conn.ContestId == submissionFull.ContestID &&
			conn.ContestProblemId == submissionFull.ContestTaskID) || (submissionFull.ContestID == -1 &&
			conn.ProblemId == base64.StdEncoding.EncodeToString([]byte(submissionFull.TestingSystem+submissionFull.TaskID))) {
			conn.Connection.WriteMessage(websocket.TextMessage, message)
		}
	}
}
