package usecase

import "github.com/BaldiSlayer/t-bmstu/services/task/internal/domain"

type TaskUseCase struct {
	repo domain.TaskRepository
}

func New(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) ListTasksPaginated(limit, offset int) ([]domain.Task, error) {
	return uc.repo.GetPaginated(limit, offset)
}

func (uc *TaskUseCase) GetTask(id int) (*domain.Task, error) {
	return uc.repo.GetByID(id)
}
