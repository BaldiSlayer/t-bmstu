package http

import (
	"encoding/json"
	"net/http"

	"github.com/BaldiSlayer/t-bmstu/services/auth/internal/usecase"
	"go.uber.org/zap"
)

type AuthHandler struct {
	useCase *usecase.Auth
	logger  *zap.Logger
}

func NewAuthHandler(uc *usecase.Auth, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		useCase: uc,
		logger:  logger,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("failed to decode login request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)

		return
	}

	token, err := h.useCase.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Warn("login failed", zap.String("username", req.Username), zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnauthorized)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("failed to decode register request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.useCase.Register(r.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		h.logger.Warn("registration failed", zap.String("username", req.Username), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("user registered successfully", zap.String("username", req.Username))

	w.WriteHeader(http.StatusCreated)
}
