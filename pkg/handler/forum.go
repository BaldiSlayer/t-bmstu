package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) forumMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "forum.tmpl", gin.H{})
}
