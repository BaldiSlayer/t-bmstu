package testsystems

import (
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems/acmp"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems/timus"
	"github.com/Baldislayer/t-bmstu/pkg/websockets"
	"sync"
)

// AllowedTestsystems - разрешенные (добавленные) тестирующие системы
var AllowedTestsystems = []TestSystem{
	&timus.Timus{Name: "timus"},
	&acmp.ACMP{Name: "acmp"},
}

// TestSystem - это интерфейс класса тестирующей системы, то есть все тестирующие системы должны обладать этими функциями
type TestSystem interface {
	Init()
	// GetName - получить имя тестирущей системы
	GetName() string
	// CheckLanguage - проверяет, существует ли у данной тестирующей системы такой язык программирования
	CheckLanguage(language string) bool
	// GetLanguages - получить языки на которых можно сдавать в этой тестирующей системе
	GetLanguages() []string
	// Submitter - воркер, который занимается отправлением посылок, и будет запускаться в отдельной горутине
	Submitter(wg *sync.WaitGroup, ch chan<- database.Submission)
	// GetProblem - получить условие задачи !!! ALERT, его надо получать по частям, см -> database.Task
	GetProblem(taskID string) (database.Task, error)
	// Checker - воркер, который занимается обновлением статусов посылок
	Checker(wg *sync.WaitGroup, ch chan<- database.Submission)
}

var wg sync.WaitGroup

func InitGorutines() error {
	submitterChannels := make(map[string]chan database.Submission)
	checkerChannels := make(map[string]chan database.Submission)

	for _, TestSystem := range AllowedTestsystems {
		ch1 := make(chan database.Submission)
		submitterChannels[TestSystem.GetName()] = ch1

		ch2 := make(chan database.Submission)
		checkerChannels[TestSystem.GetName()] = ch2

		// инициализация
		TestSystem.Init()

		wg.Add(2)

		go TestSystem.Submitter(&wg, ch1)
		go TestSystem.Checker(&wg, ch2)
	}

	// запустим горутины для самбиттеров
	for _, ch := range submitterChannels {
		go func(c <-chan database.Submission) {
			for msg := range c {
				// надо обновить запись в базе данных
				err := database.UpdateSubmissionData(msg)
				if err != nil {
					fmt.Println(err)
				}

				// передать по веб-сокету
				go websockets.SendMessageToUser(msg.SenderLogin, msg)
			}
		}(ch)
	}

	// для чекеров
	for _, ch := range checkerChannels {
		go func(c <-chan database.Submission) {
			for msg := range c {
				// надо обновить запись в базе данных
				err := database.UpdateSubmissionData(msg)
				if err != nil {
					fmt.Println(err)
				}

				// если проверка была окончена, то записать это в соответствующем поле у пользователя и в контест
				if msg.Status == 2 {

				}

				// передать по веб-сокету
				go websockets.SendMessageToUser(msg.SenderLogin, msg)
			}
		}(ch)
	}

	return nil
}
