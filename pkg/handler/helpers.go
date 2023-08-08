// тут будут лежать функции-помощники для обработки хэндлеров

package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/Baldislayer/t-bmstu/pkg/tasks_websocket"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems"
	"strings"
	"time"
)

type TaskInfo struct {
	id          string
	onlineJudge testsystems.TestSystem
}

func TaskInfoById(s string) (TaskInfo, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		fmt.Println("Ошибка при декодировании:", err)
		return TaskInfo{}, err
	}
	decodedString := string(decodedBytes)

	for _, system := range testsystems.AllowedTestsystems {
		if strings.HasPrefix(decodedString, system.GetName()) {
			// Нашли систему, разделяем online_judge и id
			task := TaskInfo{
				id:          decodedString[len(system.GetName()):],
				onlineJudge: system,
			}
			return task, nil
		}
	}

	return TaskInfo{}, fmt.Errorf("неподдерживаемая система или некорректный формат")
}

func GetTaskParts(taskId string, taskInfo *TaskInfo) (repository.Task, error) {
	exist, taskParts, err := repository.TaskExist(taskId)

	if !exist {
		taskParts, err = taskInfo.onlineJudge.GetProblem(taskInfo.id)

		if err != nil {
			return repository.Task{}, err
		}

		taskParts.ID = taskId
		// надо добавить в базу данных
		repository.AddProblem(taskParts)
	}

	return taskParts, nil
}

func GetTaskPartsById(task_id string) (TaskInfo, repository.Task, error) {
	taskInfo, err := TaskInfoById(task_id)

	if err != nil {
		return TaskInfo{}, repository.Task{}, err
	}

	taskParts, err := GetTaskParts(task_id, &taskInfo)

	return taskInfo, taskParts, err
}

func TaskSubmit(myTaskId string, login string, SourceCode string, Language string, contestId int, contestTaskId int) error {
	taskInfo, err := TaskInfoById(myTaskId)

	if err != nil {
		return err
	}

	currentTime := time.Now()
	currentTimeString := currentTime.Format("2006-01-02 15:04:05")
	submission := repository.Submission{
		SenderLogin:      login,
		TaskID:           taskInfo.id,
		TestingSystem:    taskInfo.onlineJudge.GetName(),
		Code:             []byte(SourceCode),
		Language:         Language,
		ContestTaskID:    contestTaskId,
		ContestID:        contestId,
		SubmissionTime:   currentTimeString,
		Verdict:          "Waiting",
		ExecutionTime:    "-",
		MemoryUsed:       "-",
		Test:             "-",
		SubmissionNumber: "-",
	}

	id, err := repository.AddSubmission(submission)
	if err != nil {
		return err
	}

	submission.ID = id
	go tasks_websocket.SendMessageToUser(submission.SenderLogin, submission)

	return err
}
