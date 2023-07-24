package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"html/template"
	"strings"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	store := cookie.NewStore([]byte(viper.GetString("SessionSecret")))
	store.Options(sessions.Options{
		Path:     "/", // Установка пути для куки на "/"
		MaxAge:   86400,
		HttpOnly: true,
		// TODO add Secure and other need able options
	})
	router.Use(sessions.Sessions(sessionName, store))

	router.SetFuncMap(template.FuncMap{
		"nl2br": nl2br,
	})
	router.LoadHTMLGlob("web/templates/*")

	auth := router.Group("/auth")
	{
		auth.GET("/login", h.signUp)
		auth.GET("/callback", h.callback)
	}

	//api := router.Group("/api")
	// api.Use(authMiddleware())
	//{
	//	api.GET("/problem/:id")
	//}

	view := router.Group("/view")
	view.Use(authMiddleware())
	{
		view.GET("/timus", h.timusTaskList)
		view.GET("/problem/:id", h.getTask)
		view.POST("/problem/:id/submit", h.submitTask)

		view.GET("/contests", h.getContests)

		contest := view.Group("/contest/:contest_id")
		{
			contest.GET("/tasks", h.getContestTasks)
			contest.GET("/task/:task_id", h.getContestTask)
			contest.POST("/task/:task_id/submit", h.submitContestTask)
		}
	}

	return router
}

func nl2br(s string) template.HTML {
	return template.HTML(strings.ReplaceAll(s, "\n", "<br>"))
}
