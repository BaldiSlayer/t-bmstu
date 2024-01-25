package handler

import (
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) profileMainPage(c *gin.Context) {
	profile, err := database.GetInfoForProfilePage(c.GetString("username"))

	if err != nil {
		// TODO return error
		return
	}

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"NickName": profile.Username,
		"Surname":  profile.LastName,
		"Name":     profile.FirstName,
		"Email":    profile.Email,
	})
}
