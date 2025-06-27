package inmemory

import (
	"context"
	"errors"

	"github.com/BaldiSlayer/t-bmstu/services/task/internal/domain"
)

type InMemoryRepo struct {
	tasks map[int]domain.Task
}

var _ domain.TaskRepository = (*InMemoryRepo)(nil)

func New() *InMemoryRepo {
	return &InMemoryRepo{
		tasks: map[int]domain.Task{
			1: {ID: 1, Title: "A + B", Description: "Compute sum"},
			2: {ID: 2, Title: "Max of array", Description: "Find max"},
		},
	}
}

func (r *InMemoryRepo) getAll() []domain.Task {
	tasks := make([]domain.Task, 0, len(r.tasks))

	for _, v := range r.tasks {
		tasks = append(tasks, v)
	}

	return tasks
}

func (r *InMemoryRepo) GetPaginated(ctx context.Context, limit, offset int) ([]domain.Task, error) {
	all := r.getAll()

	end := offset + limit
	if offset > len(all) {
		return []domain.Task{}, nil
	}

	if end > len(all) {
		end = len(all)
	}

	return all[offset:end], nil
}

func (r *InMemoryRepo) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	t, ok := r.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}

	return &t, nil
}
