package testsystems

import "github.com/Baldislayer/t-bmstu/pkg/repository"

type TestSystem interface {
	GetName() string
	GetLanguages() []string
	Submit(login string, id string, SourceCode string, Language string, contestId int, contestTaskId int) error
	GetProblem(taskID string) (repository.Task, error)
}
