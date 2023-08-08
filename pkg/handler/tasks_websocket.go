package handler

import (
	"github.com/Baldislayer/t-bmstu/pkg/websockets"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

	// Добавляем соединение в мапу
	websockets.Mu.Lock()
	websockets.Connections[username] = conn
	websockets.Mu.Unlock()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Ошибка чтения сообщения от пользователя %s: %v", username, err)
			}
			break
		}

		log.Printf("Получено сообщение от пользователя %s: %s", username, p)

		// Пример эхо-ответа (отправка обратно полученного сообщения)
		if err := conn.WriteMessage(websocket.TextMessage, p); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Ошибка отправки сообщения пользователю %s: %v", username, err)
			}
			break
		}
	}

	// Удаляем соединение из мапы при закрытии
	websockets.Mu.Lock()
	delete(websockets.Connections, username)
	websockets.Mu.Unlock()
}

func (h *Handler) Htmlsome(c *gin.Context) {
	c.HTML(http.StatusOK, "some.tmpl", gin.H{})
}
