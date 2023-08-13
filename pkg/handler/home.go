package handler

import (
	"encoding/json"
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/gin-gonic/gin"
)

func (h *Handler) add(c *gin.Context) {
	members := []json.RawMessage{
		json.RawMessage(`{"username": "sh", "role": "student"}`),
	}

	database.AddGroupWithMembers(database.Group{
		Title:    "smth 2",
		Students: []string{"sh"},
	},
		members)
}
