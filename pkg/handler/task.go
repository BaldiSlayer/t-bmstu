package handler

import (
	"github.com/Baldislayer/t-bmstu/pkg/testsystems/timus"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) timusTaskList(c *gin.Context) {
	count := c.Query("count")

	parsedCount, err := strconv.Atoi(count)
	if err != nil {
		parsedCount = 15
	}

	if parsedCount > 50 {
		parsedCount = 50
	}

	taskList, err := timus.GetTaskList(parsedCount)

	if err != nil {
		c.JSON(http.StatusBadRequest, "bad req")
	}

	c.HTML(http.StatusOK, "ts_tasks_list.tmpl", gin.H{
		"Tasks": taskList,
	})
}

func (h *Handler) submitTask(c *gin.Context) {
	var requestData struct {
		SourceCode string `json:"sourceCode"`
		Language   string `json:"language"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestData.Language == "Select language " {
		c.JSON(http.StatusBadRequest, "There are not such language")
		return
	}

	err := TaskSubmit(c.Param("id"), c.GetString("username"), requestData.SourceCode, requestData.Language,
		-1, -1)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task submitted successfully",
	})
}
