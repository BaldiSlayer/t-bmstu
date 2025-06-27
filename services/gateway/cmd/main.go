package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BaldiSlayer/t-bmstu/services/gateway/internal/middleware"
	"github.com/BaldiSlayer/t-bmstu/services/gateway/internal/proxy"
	"go.uber.org/zap"
)

const jwtSecretString = "super_secret"

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	jwtSecret := []byte(jwtSecretString)

	routes := []proxy.Route{
		{Prefix: "/auth/", Target: "http://auth:8081"},
		{Prefix: "/tasks", Target: "http://task:8082"},
	}

	router := proxy.NewProxyRouter(routes, logger)

	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)

	handler := authMiddleware.ServeHTTP(router)

	log.Println("Gateway started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
