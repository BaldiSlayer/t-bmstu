package handler

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	/*
		store := cookie.NewStore([]byte(viper.GetString("SessionSecret")))
		store.Options(sessions.Options{
			Path:     "/", // Установка пути для куки на "/"
			MaxAge:   86400,
			HttpOnly: true,
			// TODO add Secure and other need able options
		})
		router.Use(sessions.Sessions(sessionName, store))
	*/

	router.SetFuncMap(template.FuncMap{
		"nl2br": nl2br,
		"inc": func(index int) int {
			return index + 1
		},
	})
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/images", "web/static/images")
	router.Static("/styles", "web/static/styles")
	router.Static("/scripts", "web/static/scripts")

	auth := router.Group("/auth")
	{

		auth.GET("/login", h.signIn)
		auth.POST("/login", h.signIn)
		auth.GET("/registration", h.signUp)
		auth.POST("/registration", h.signUp)

		github := auth.Group("/github")
		github.GET("/login", h.githubSignUp)
		github.GET("/githubCallback", h.githubCallback)
	}

	api := router.Group("/api")
	api.Use(authMiddleware())
	{
		api.GET("/ws/contest/:contest_id/problem/:problem_id", h.handleWebSocket)
	}

	view := router.Group("/view")
	view.Use(authMiddleware())
	{
		view.GET("/add", h.add)

		forum := view.Group("/forum")
		{
			forum.GET("/")
		}

		// view.GET("/websockets", h.Htmlsome)
		view.GET("/home", h.home)
		view.GET("/timus", h.timusTaskList)

		problem := view.Group("/problem")
		{
			problem.GET("/:id", h.getTask)
			problem.POST("/:id/submit", h.submitTask)
		}

		// view.GET("/contests", h.getContests)

		contest := view.Group("/contest/:contest_id")
		{
			contest.GET("/problems", h.getContestTasks)
			// TODO сделать одинаковую функцию getTask
			contest.GET("/problem/:problem_id", h.getContestTask)
			contest.POST("/problem/:problem_id/submit", h.submitContestTask)
		}

		groups := view.Group("/group")
		{
			groups.GET("/invite/:invite_hash", h.checkInvite)
			group := groups.Group("/:group_id")
			{
				group.GET("", h.getGroupContests)
				groupContest := group.Group("/contest/:contest_id")
				{
					groupContest.GET("/tasks", h.getContestTasks)
					// TODO вести дальше до задач
				}
			}
		}
	}

	return router
}

func nl2br(s string) template.HTML {
	return template.HTML(strings.ReplaceAll(s, "\n", "<br>"))
}
