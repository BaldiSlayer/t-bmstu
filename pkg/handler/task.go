package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getTask(c *gin.Context) {
	taskId := c.Param("id")
	taskInfo, err := TaskInfoById(taskId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	taskParts, err := GetTaskParts(taskId, &taskInfo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	submissons := repository.GetVerdicts(c.GetString("email"), taskInfo.id, taskInfo.onlineJudge.GetName())

	type Test struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	}

	var tests []Test
	if testStr, ok := taskParts.Tests["tests"].(string); ok {
		err := json.Unmarshal([]byte(testStr), &tests)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	c.HTML(http.StatusOK, "task.tmpl", gin.H{
		"Task":        taskParts,
		"Tests":       tests,
		"Languages":   taskInfo.onlineJudge.GetLanguages(),
		"Submissions": submissons,
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

	err := TaskSubmit(c.Param("id"), c.GetString("email"), requestData.SourceCode, requestData.Language,
		-1, -1)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task submitted successfully",
	})
}

func (h *Handler) getContests(c *gin.Context) {
	// пока что я буду отображать все контесты, которые у меня есть
	// потом сделать так, чтобы отображались только те, в которые он может зайти

	constests := repository.GetContests()

	c.HTML(http.StatusOK, "contests.tmpl", gin.H{
		"Contests": constests,
	})
}
