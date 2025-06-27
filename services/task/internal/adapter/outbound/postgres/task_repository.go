package postgres

import (
	"context"
	"errors"

	"github.com/BaldiSlayer/t-bmstu/services/task/internal/domain"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) GetPaginated(ctx context.Context, limit, offset int) ([]domain.Task, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, description FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	var t domain.Task
	err := r.db.QueryRow(ctx, `SELECT id, title, description FROM tasks WHERE id=$1`, id).
		Scan(&t.ID, &t.Title, &t.Description)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("task not found")
	}

	if err != nil {
		return nil, err
	}

	return &t, nil
}
