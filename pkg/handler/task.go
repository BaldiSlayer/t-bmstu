package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems/timus"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Task struct {
	// id задания на данной тестирующей системе
	id string
	// codeforces, timus
	online_judge string
}

func (h *Handler) emailGet(c *gin.Context) {
	respBody, ok := c.Get("email")
	if !ok {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	bodyBytes, err := ioutil.ReadAll(respBody.(io.Reader))
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read response body")
		return
	}

	email := string(bodyBytes)
	var emails []map[string]interface{}
	err = json.Unmarshal([]byte(email), &emails)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Email: %s", emails[0]["email"]))
}

func taskInfoById(s string) (Task, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		fmt.Println("Ошибка при декодировании:", err)
		return Task{}, err
	}
	decodedString := string(decodedBytes)

	allowedTestsystems := []string{"codeforces", "timus"} // Массив поддерживаемых систем

	for _, system := range allowedTestsystems {
		if strings.HasPrefix(decodedString, system) {
			// Нашли систему, разделяем online_judge и id
			task := Task{
				id:           decodedString[len(system):],
				online_judge: system,
			}
			return task, nil
		}
	}

	return Task{}, fmt.Errorf("неподдерживаемая система или некорректный формат")
}

func (h *Handler) getTask(c *gin.Context) {
	taskInfo, err := taskInfoById(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if taskInfo.online_judge == "timus" {
		htmlTask, err := timus.GetProblem(taskInfo.id)

		if err != nil {
			fmt.Printf("Failed to get task: %s", err)
		}

		// получить все посылки по этой задаче у этого юзера на тимусе
		submissons := repository.GetVerdicts(c.GetString("email"), taskInfo.id, taskInfo.online_judge)

		task := template.HTML(htmlTask)
		c.HTML(http.StatusOK, "task.tmpl", gin.H{
			"Statement":   task,
			"Languages":   timus.Languages,
			"Submissions": submissons,
		})
	}
}

func (h *Handler) submitTask(c *gin.Context) {
	task_info, err := taskInfoById(c.Param("id"))

	if err != nil {
		fmt.Printf("Error occured: %s", err)
	}

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

	// Примеры работы с таблицей "submissions"
	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")
	submission := repository.Submission{
		SenderLogin:    c.GetString("email"),
		TaskID:         task_info.id,
		TestingSystem:  task_info.online_judge,
		Code:           requestData.SourceCode,
		Language:       requestData.Language,
		ContestTaskID:  -1,
		ContestID:      -1,
		SubmissionTime: currentTimeString,
		SVerdictID:     "-",
	}

	repository.AddSubmission(submission)

	// err = timus.SendSubmission("342187EL", timus.Codes[requestData.Language], task_info.id, requestData.SourceCode)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"message": "Task was not submitted",
	//	})
	//	return
	//}

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
