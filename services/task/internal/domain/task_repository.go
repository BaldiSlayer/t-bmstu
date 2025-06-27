package domain

import "context"

type TaskRepository interface {
	GetPaginated(ctx context.Context, limit, offset int) ([]Task, error)
	GetByID(ctx context.Context, id int) (*Task, error)
}
