package testsystems

import (
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems/codeforces"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems/timus"
)

// AllowedTestsystems - разрешенные (добавленные) тестирующие системы
var AllowedTestsystems = []TestSystem{
	&timus.Timus{Name: "timus"},
	&codeforces.Codeforces{Name: "codeforces"},
}

// TestSystem - это интерфейс класса тестирующей системы, то есть все тестирующие системы должны обладать этими функциями
type TestSystem interface {
	// GetName - получить имя тестирущей системы
	GetName() string
	// GetLanguages - получить языки на которых можно сдавать в этой тестирующей системе
	GetLanguages() []string
	// Submit - отправление решения
	Submit(login string, id string, SourceCode string, Language string, contestId int, contestTaskId int) error
	// GetProblem - получить условие задачи !!! ALERT, его надо получать по частям, см -> repository.Task
	GetProblem(taskID string) (repository.Task, error)
}
