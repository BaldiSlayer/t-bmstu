package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/BaldiSlayer/t-bmstu/services/auth/internal/domain"
)

type TokenManager interface {
	Generate(username, role string) (string, error)
}

type Auth struct {
	userRepo     domain.UserRepository
	tokenManager TokenManager
}

func NewAuth(userRepo domain.UserRepository, tokenManager TokenManager) *Auth {
	return &Auth{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (a *Auth) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	if user.Password != password {
		return "", errors.New("invalid password")
	}

	return a.tokenManager.Generate(user.Username, user.Role)
}

func (a *Auth) Register(ctx context.Context, username, password, role string) error {
	existing, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to find user by username: %w", err)
	}

	if existing != nil {
		return errors.New("user already exists")
	}

	user := &domain.User{
		Username: username,
		Password: password,
		Role:     role,
	}

	err = a.userRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
