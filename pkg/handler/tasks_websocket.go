package handler

import (
	"github.com/Baldislayer/t-bmstu/pkg/websockets"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) handleWebSocket(c *gin.Context) {
	conn, err := websockets.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Ошибка при обновлении соединения:", err)
		return
	}
	defer conn.Close()

	// Получение юзернейма из контекста Gin
	username := c.GetString("username")
	if username == "" {
		log.Println("Ошибка: не удалось получить юзернейм из контекста")
		return
	}

	// Получение параметра problem_id из пути
	problemID := c.Param("problem_id")

	val := c.Param("contest_id")
	contestID, err := strconv.Atoi(val)
	if err != nil || val == "" {
		contestID = -1
	}

	contestProblemID := -1
	if contestID != -1 {
		contestProblemID, err = strconv.Atoi(problemID)
		if err != nil {
			return
		}
	}

	// Добавляем соединение в мапу
	websockets.Mu.Lock()
	connections, ok := websockets.Connections[username]
	if !ok {
		connections = []websockets.UserConnections{}
	}
	connections = append(connections, websockets.UserConnections{
		Connection:       conn,
		ProblemId:        problemID,
		ContestId:        contestID,
		ContestProblemId: contestProblemID,
	})
	websockets.Connections[username] = connections
	websockets.Mu.Unlock()

	// Проверка на разрыв соединения
	for i, connData := range connections {
		_, _, err := connData.Connection.ReadMessage()
		if err != nil {
			// Соединение разорвано, удаляем его из списка
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}
	websockets.Mu.Lock()
	websockets.Connections[username] = connections
	websockets.Mu.Unlock()
}

func (h *Handler) Htmlsome(c *gin.Context) {
	c.HTML(http.StatusOK, "some.tmpl", gin.H{})
}
