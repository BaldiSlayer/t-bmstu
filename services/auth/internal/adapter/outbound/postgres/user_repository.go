package postgres

import (
	"context"

	"github.com/BaldiSlayer/t-bmstu/services/auth/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var u domain.User

	err := ur.db.QueryRow(
		ctx,
		`SELECT WHERE username = $1`, username,
	).Scan(&u.Username, &u.Password, &u.Role)

	return &u, err
}

func (ur *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := ur.db.Exec(ctx,
		"INSERT INTO users (username, password, role) VALUES ($1, $2, $3)",
		user.Username, user.Password, user.Role,
	)

	return err
}
