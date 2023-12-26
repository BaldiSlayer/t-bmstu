package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) profileMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"NickName": "BaldiSlayer",
		"Surname":  "Лисов",
		"Name":     "Алексей",
		"Email":    "lisov.a2005@yandex.ru",
	})
}
