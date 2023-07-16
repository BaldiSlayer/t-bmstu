package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	store := cookie.NewStore([]byte(viper.GetString("SessionSecret")))
	store.Options(sessions.Options{
		HttpOnly: true,
		// TODO add Secure and other need able options
	})
	router.Use(sessions.Sessions(sessionName, store))

	router.LoadHTMLGlob("web/templates/*")

	auth := router.Group("/auth")
	{
		auth.GET("/login", h.signUp)
		auth.GET("/callback", h.callback)
	}

	api := router.Group("/view")
	api.Use(authMiddleware())
	{
		api.GET("/problem/:id", h.getTask)
		api.POST("/problem/:id/submit", h.submitTask)
		api.GET("/contests", h.getContests)
	}

	return router
}
