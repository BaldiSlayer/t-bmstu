package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/gin-gonic/gin"
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

	contest, err := repository.GetContestInfoById(contestId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	tasks := contest.Tasks

	task_list := []ContestTableTask{}

	for key, value := range tasks {
		_, taskParts, err := GetTaskPartsById(value.(string))

		if err != nil {
			fmt.Print(err)
		} else {
			task_list = append(task_list, ContestTableTask{ID: key, Name: taskParts.Name})
		}
	}

	sort.Sort(ByID(task_list))

	c.HTML(http.StatusOK, "constest_tasks_list.tmpl", gin.H{
		"Tasks": task_list,
	})
}

func (h *Handler) getContestTask(c *gin.Context) {
	contestId, err := strconv.Atoi(c.Param("contest_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	taskId := c.Param("task_id")
	contest, err := repository.GetContestInfoById(contestId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	tasks := contest.Tasks

	value, ok := tasks[taskId]

	if ok {
		taskInfo, taskParts, err := GetTaskPartsById(value.(string))

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

		intContestTaskId, err := strconv.Atoi(taskId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		submissions, err := repository.GetVerditctsOfContestTask(c.GetString("username"), contestId, intContestTaskId)

		c.HTML(http.StatusOK, "task.tmpl", gin.H{
			"Task":        taskParts,
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

	contestTaskId, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	contestId, err := strconv.Atoi(c.Param("contest_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	contest, err := repository.GetContestInfoById(contestId)
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
