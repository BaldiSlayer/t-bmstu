package handler

import (
	"encoding/json"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

// в домашней странице будут показываться группы в которых пользователь является кем-либо
func (h *Handler) home(c *gin.Context) {
	role := c.GetString("role")

	switch role {
	case "student":
		{
			groups, err := repository.GetUserGroups(c.GetString("username"))

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}

			c.HTML(http.StatusOK, "groups.tmpl", gin.H{
				"Groups": groups,
			})
		}
	case "teacher":
		{
			c.JSON(http.StatusOK, gin.H{"msg": "Hello, Teacher"})
		}
	case "admin":
		{
			c.JSON(http.StatusOK, gin.H{"msg": "Hello, admin"})
		}
	default:
		{
			c.JSON(http.StatusBadRequest, gin.H{"error": "hacking attempt"})
		}
	}
}

func (h *Handler) add(c *gin.Context) {
	members := []json.RawMessage{
		json.RawMessage(`{"username": "sh", "role": "student"}`),
	}

	repository.AddGroupWithMembers(repository.Group{
		Title:    "smth 2",
		Students: []string{"sh"},
	},
		members)
}
