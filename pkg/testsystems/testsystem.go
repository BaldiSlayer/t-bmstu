package testsystems

import (
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/Baldislayer/t-bmstu/pkg/tasks_websocket"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems/timus"
	"sync"
)

// AllowedTestsystems - разрешенные (добавленные) тестирующие системы
var AllowedTestsystems = []TestSystem{
	&timus.Timus{Name: "timus"},
}

// TestSystem - это интерфейс класса тестирующей системы, то есть все тестирующие системы должны обладать этими функциями
type TestSystem interface {
	// GetName - получить имя тестирущей системы
	GetName() string
	// GetLanguages - получить языки на которых можно сдавать в этой тестирующей системе
	GetLanguages() []string
	// Submitter - воркер, который занимается отправлением решений, и будет запускаться в отдельной горутине
	Submitter(wg *sync.WaitGroup, ch chan<- repository.Submission)
	// GetProblem - получить условие задачи !!! ALERT, его надо получать по частям, см -> repository.Task
	GetProblem(taskID string) (repository.Task, error)

	Checker(wg *sync.WaitGroup, ch chan<- repository.Submission)
}

var wg sync.WaitGroup

func InitGorutines() error {
	submitterChannels := make(map[string]chan repository.Submission)
	checkerChannels := make(map[string]chan repository.Submission)

	for _, TestSystem := range AllowedTestsystems {
		ch1 := make(chan repository.Submission)
		submitterChannels[TestSystem.GetName()] = ch1

		ch2 := make(chan repository.Submission)
		checkerChannels[TestSystem.GetName()] = ch2

		wg.Add(2)

		go TestSystem.Submitter(&wg, ch1)
		go TestSystem.Checker(&wg, ch2)
	}

	// запустим горутины для самбиттеров
	for _, ch := range submitterChannels {
		go func(c <-chan repository.Submission) {
			for msg := range c {
				// надо обновить запись в базе данных
				err := repository.UpdateSubmissionData(msg)
				if err != nil {
					fmt.Println(err)
				}

				// передать по веб-сокету
				go tasks_websocket.SendMessageToUser(msg.SenderLogin, msg)
			}
		}(ch)
	}

	// для чекеров
	for _, ch := range checkerChannels {
		go func(c <-chan repository.Submission) {
			for msg := range c {
				// надо обновить запись в базе данных
				err := repository.UpdateSubmissionData(msg)
				if err != nil {
					fmt.Println(err)
				}

				// если проверка была окончена, то записать это в соответствующем поле у пользователя

				// передать по веб-сокету
				go tasks_websocket.SendMessageToUser(msg.SenderLogin, msg)
			}
		}(ch)
	}

	return nil
}
