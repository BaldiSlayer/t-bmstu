package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/BaldiSlayer/t-bmstu/services/auth/internal/adapter/outbound/postgres"
	"github.com/BaldiSlayer/t-bmstu/services/auth/internal/app"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	adapter_http "github.com/BaldiSlayer/t-bmstu/services/auth/internal/adapter/inbound/http"
)

const jwtSecretString = "super_secret"

type JWTSigner struct {
	secret []byte
}

func NewJWTSigner(secret []byte) *JWTSigner {
	return &JWTSigner{secret: secret}
}

func (j *JWTSigner) Generate(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	pool, err := pgxpool.New(context.Background(), "")
	if err != nil {
		return
	}

	userRepo := postgres.New(pool)

	tokenManager := NewJWTSigner([]byte(jwtSecretString))
	useCase := app.NewAuthUseCase(userRepo, tokenManager)

	handler := adapter_http.NewAuthHandler(useCase, logger)

	http.HandleFunc("/login", handler.Login)

	logger.Info("auth service started on :8080...")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Error("error while ListenAndServe", zap.Error(err))
	}
}
