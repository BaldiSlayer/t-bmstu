package handler

import (
	"encoding/json"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems/timus"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
)

func (h *Handler) timusTaskList(c *gin.Context) {
	count := c.Query("count")

	parsedCount, err := strconv.Atoi(count)
	if err != nil {
		// TODO
		parsedCount = 15
	}

	// TODO
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

func (h *Handler) getTask(c *gin.Context) {
	taskId := c.Param("id")
	taskInfo, err := TaskInfoById(taskId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskParts, err := GetTaskParts(taskId, &taskInfo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	submissons, err := repository.GetVerditctsOfContestTask(c.GetString("username"), -1, -1)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type Test struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	}

	var tests []Test
	if testStr, ok := taskParts.Tests["tests"].(string); ok {
		err := json.Unmarshal([]byte(testStr), &tests)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.HTML(http.StatusOK, "task-page.tmpl", gin.H{
		"Task":        taskParts,
		"Condition":   template.HTML(taskParts.Condition),
		"InputData":   template.HTML(taskParts.InputData),
		"OutputData":  template.HTML(taskParts.OutputData),
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

// TODO вынести в какой-то другой файл / убрать это по причине "все контесты" лежат в какой-то группе
func (h *Handler) getContests(c *gin.Context) {
	// пока что я буду отображать все контесты
	// сделать отображение только тех, в которые он может зайти

	contests := repository.GetContests()

	c.HTML(http.StatusOK, "contests.tmpl", gin.H{
		"Contests": contests,
	})
}
