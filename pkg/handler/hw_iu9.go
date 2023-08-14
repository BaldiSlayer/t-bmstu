package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) hwIu9MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "hw_ui9_mainpage.tmpl", gin.H{})
}
