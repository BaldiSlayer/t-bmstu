// тут будут лежать функции-помощники для обработки хэндлеров

package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems"
	"github.com/Baldislayer/t-bmstu/pkg/websockets"
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

func GetTaskParts(taskId string, taskInfo *TaskInfo) (database.Task, error) {
	exist, taskParts, err := database.TaskExist(taskId)

	if !exist {
		taskParts, err = taskInfo.onlineJudge.GetProblem(taskInfo.id)

		if err != nil {
			return database.Task{}, err
		}

		taskParts.ID = taskId
		// надо добавить в базу данных
		database.AddProblem(taskParts)
	}

	return taskParts, nil
}

func GetTaskPartsById(task_id string) (TaskInfo, database.Task, error) {
	taskInfo, err := TaskInfoById(task_id)

	if err != nil {
		return TaskInfo{}, database.Task{}, err
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
	submission := database.Submission{
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

	if !taskInfo.onlineJudge.CheckLanguage(Language) {
		return errors.New("no such language")
	}

	id, err := database.AddSubmission(submission)
	if err != nil {
		return err
	}

	submission.ID = id
	go websockets.SendMessageToUser(submission.SenderLogin, submission)

	return err
}
