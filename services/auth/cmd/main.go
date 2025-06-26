package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BaldiSlayer/t-bmstu/services/auth/internal/adapter/outbound/postgres"
	"github.com/BaldiSlayer/t-bmstu/services/auth/internal/usecase"
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

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Panic("failed to connect to the database", zap.Error(err))

		return
	}

	userRepo := postgres.New(pool)

	tokenManager := NewJWTSigner([]byte(jwtSecretString))
	useCase := usecase.NewAuth(userRepo, tokenManager)

	handler := adapter_http.NewAuthHandler(useCase, logger)

	http.HandleFunc("/auth/login", handler.Login)
	http.HandleFunc("/auth/register", handler.Register)

	logger.Info("auth service started on :8081...")

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		logger.Error("error while ListenAndServe", zap.Error(err))
	}
}
