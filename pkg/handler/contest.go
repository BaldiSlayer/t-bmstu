package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"sort"
	"strconv"
)

type ContestTableTask struct {
	ID   string
	Name string
}

type ByID []ContestTableTask

func (t ByID) Len() int           { return len(t) }
func (t ByID) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByID) Less(i, j int) bool { return t[i].ID < t[j].ID }

func (h *Handler) getContestTasks(c *gin.Context) {
	contestId, err := strconv.Atoi(c.Param("contest_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	contest, err := database.GetContestInfoById(contestId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	tasks := contest.Tasks

	taskList := []ContestTableTask{}

	for key, value := range tasks {
		_, taskParts, err := GetTaskPartsById(value.(string))

		if err != nil {
			fmt.Print(err)
		} else {
			taskList = append(taskList, ContestTableTask{ID: key, Name: taskParts.Name})
		}
	}

	sort.Sort(ByID(taskList))

	c.HTML(http.StatusOK, "contest-tasks-list.tmpl", gin.H{
		"Tasks": taskList,
	})
}

func (h *Handler) getTask(c *gin.Context) {
	stringContestId := c.Param("contest_id")
	ok := true
	value := ""
	contestId := -1
	intContestTaskId := -1

	if stringContestId != "" {
		// получаем все нужное для контеста
		contestID, err := strconv.Atoi(stringContestId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		taskId := c.Param("problem_id")
		contest, err := database.GetContestInfoById(contestID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		tasks := contest.Tasks

		problemId, exist := tasks[taskId]
		value = problemId.(string)
		ok = exist

		intContestTaskId, err = strconv.Atoi(taskId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contestId = contestID
	} else {
		value = c.Param("id")
	}

	if ok {
		taskInfo, taskParts, err := GetTaskPartsById(value)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

		parts, err := TaskInfoById(value)
		submissions := database.GetVerdicts(c.GetString("username"), parts.id, parts.onlineJudge.GetName(), contestId, intContestTaskId)

		c.HTML(http.StatusOK, "task-page.tmpl", gin.H{
			"Task":        taskParts,
			"Condition":   template.HTML(taskParts.Condition),
			"InputData":   template.HTML(taskParts.InputData),
			"OutputData":  template.HTML(taskParts.OutputData),
			"Tests":       tests,
			"Languages":   taskInfo.onlineJudge.GetLanguages(),
			"Submissions": submissions,
		})
	} else {
		c.JSON(http.StatusBadRequest, "No such task")
	}
}

func (h *Handler) submitContestTask(c *gin.Context) {
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

	contestTaskId, err := strconv.Atoi(c.Param("problem_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	contestId, err := strconv.Atoi(c.Param("contest_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	contest, err := database.GetContestInfoById(contestId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	tasks := contest.Tasks

	value, ok := tasks[strconv.Itoa(contestTaskId)]

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no such task"})
		return
	}

	err = TaskSubmit(value.(string), c.GetString("username"), requestData.SourceCode, requestData.Language,
		contestId, contestTaskId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task submitted successfully",
	})
}
