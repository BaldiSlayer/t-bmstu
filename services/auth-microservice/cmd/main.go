package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BaldiSlayer/t-bmstu/services/auth-microservice/internal/adapter/outbound/postgres"
	"github.com/BaldiSlayer/t-bmstu/services/auth-microservice/internal/app"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	adapter_http "github.com/BaldiSlayer/t-bmstu/services/auth-microservice/internal/adapter/inbound/http"
)

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
	pool, err := pgxpool.New(context.Background(), "")
	if err != nil {
		os.Exit(1)
	}

	userRepo := postgres.New(pool)

	tokenManager := NewJWTSigner([]byte("super_secret"))
	useCase := app.NewAuthUseCase(userRepo, tokenManager)

	handler := adapter_http.NewAuthHandler(useCase)

	http.HandleFunc("/login", handler.Login)

	log.Println("auth service started on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
