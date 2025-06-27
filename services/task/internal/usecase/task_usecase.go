package usecase

import (
	"context"

	"github.com/BaldiSlayer/t-bmstu/services/task/internal/domain"
)

type TaskUseCase struct {
	repo domain.TaskRepository
}

func New(repo domain.TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) ListTasksPaginated(ctx context.Context, limit, offset int) ([]domain.Task, error) {
	return uc.repo.GetPaginated(ctx, limit, offset)
}

func (uc *TaskUseCase) GetTask(ctx context.Context, id int) (*domain.Task, error) {
	return uc.repo.GetByID(ctx, id)
}
