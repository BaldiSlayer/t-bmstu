package codeforces

import (
	"github.com/Baldislayer/t-bmstu/pkg/repository"
)

type Codeforces struct {
	Name string
}

func (cf *Codeforces) GetName() string {
	return cf.Name
}

func (cf *Codeforces) GetLanguages() []string {
	return []string{}
}

func (cf *Codeforces) Submit(login string, id string, SourceCode string, Language string, contestId int, contestTaskId int) error {
	return nil
}

func (cf *Codeforces) GetProblem(taskID string) (repository.Task, error) {
	return repository.Task{}, nil
}
