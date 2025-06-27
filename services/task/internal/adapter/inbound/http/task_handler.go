package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BaldiSlayer/t-bmstu/services/task/internal/usecase"
	"go.uber.org/zap"
)

type TaskHandler struct {
	uc     *usecase.TaskUseCase
	logger *zap.Logger
}

func NewTaskHandler(uc *usecase.TaskUseCase, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{
		uc:     uc,
		logger: logger,
	}
}

func (t *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	tasks, err := t.uc.ListTasksPaginated(r.Context(), limit, offset)
	if err != nil {
		t.logger.Error("failed to list tasks", zap.Error(err))
		http.Error(w, "failed to fetch tasks", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (t *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = 0
	}

	task, err := t.uc.GetTask(r.Context(), id)
	if err != nil {
		t.logger.Warn("task not found", zap.Int("id", id), zap.Error(err))
		http.Error(w, "task not found", http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
