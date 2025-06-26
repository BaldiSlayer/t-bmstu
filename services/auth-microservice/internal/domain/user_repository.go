package domain

import "context"

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
}
