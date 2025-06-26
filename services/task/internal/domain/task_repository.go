package domain

type TaskRepository interface {
	GetPaginated(limit, offset int) ([]Task, error)
	GetByID(id int) (*Task, error)
}
