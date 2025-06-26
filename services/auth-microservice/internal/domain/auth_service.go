package domain

import "context"

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password, role string) error
}
