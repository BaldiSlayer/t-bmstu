package handler

import (
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) getGroupContests(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("group_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	contests, err := database.GetGroupContests(groupId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.HTML(http.StatusOK, "group_contests.tmpl", gin.H{
		"Contests": contests,
	})
}

func (h *Handler) checkInvite(c *gin.Context) {
	inviteHash := c.Param("invite_hash")

	if inviteHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no such group"})
		return
	}

	exist, groupId, err := database.CheckInviteCode(inviteHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no such group"})
		return
	}

	database.AddUserToGroup(c.GetString("username"), groupId, "student")
	c.JSON(http.StatusOK, gin.H{"Success": "U are member of this group now"})
}

//
