package tasks_websocket

import (
	"encoding/json"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
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

var Connections = make(map[string]*websocket.Conn) // Мапа для хранения соединений по юзернейму
var Mu sync.Mutex                                  // Мьютекс для безопасного доступа к мапе connections

// Функция для отправки объекта Submission по веб-сокету по нику
func SendMessageToUser(username string, submission_full repository.Submission) {
	Mu.Lock()
	conn, ok := Connections[username]
	Mu.Unlock()

	if !ok {
		log.Printf("Ошибка: соединение не найдено для пользователя %s", username)
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
		ID:            submission_full.ID,
		Verdict:       submission_full.Verdict,
		Language:      submission_full.Language,
		ExecutionTime: submission_full.ExecutionTime,
		MemoryUsed:    submission_full.MemoryUsed,
		Test:          submission_full.Test,
	}

	// Преобразуем объект Submission в JSON
	message, err := json.Marshal(submission)
	if err != nil {
		log.Printf("Ошибка при преобразовании Submission в JSON: %v", err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, message)
}
