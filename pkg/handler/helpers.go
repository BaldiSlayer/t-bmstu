// тут будут лежать функции-помощники для обработки хэндлеров

package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems/timus"
	"strings"
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

	allowedTestsystems := []testsystems.TestSystem{
		&timus.Timus{Name: "timus"},
	}

	for _, system := range allowedTestsystems {
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

func TaskSubmit(myTaskId string, login string, SourceCode string, Language string, contest_id int, contstTaskId int) error {
	taskInfo, err := TaskInfoById(myTaskId)

	if err != nil {
		return err
	}

	taskInfo.onlineJudge.Submit(login,
		taskInfo.id,
		SourceCode,
		Language,
		contest_id,
		contstTaskId)

	return nil
}
