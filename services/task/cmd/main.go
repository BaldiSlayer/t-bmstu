package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	adapter_http "github.com/BaldiSlayer/t-bmstu/services/task/internal/adapter/inbound/http"
	"github.com/BaldiSlayer/t-bmstu/services/task/internal/adapter/outbound/inmemory"
	"github.com/BaldiSlayer/t-bmstu/services/task/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

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

	_ = pool

	userRepo := inmemory.New()

	uc := usecase.New(userRepo)

	handler := adapter_http.NewTaskHandler(uc, logger)

	http.HandleFunc("/tasks", handler.ListTasks)
	http.HandleFunc("/task/{id}", handler.GetTaskByID)

	logger.Info("auth service started on :8082...")

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		logger.Error("error while ListenAndServe", zap.Error(err))
	}
}
