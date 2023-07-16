package handlers

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.GET("/login")
		auth.GET("/callback")
	}

	api := router.Group("/api")
	{
		api.GET("email")
	}
}
