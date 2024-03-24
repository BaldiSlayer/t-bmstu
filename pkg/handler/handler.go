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

	router.SetFuncMap(template.FuncMap{
		"nl2br": nl2br,
		"inc": func(index int) int {
			return index + 1
		},
		"odd": func(index int) bool {
			return index%2 == 0
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
		view.GET("/createGroup", h.createGroup)
		view.POST("/createGroup", h.createGroup)

		view.GET("/kostyl", h.addContest)

		view.GET("/home", h.home)

		forum := view.Group("/forum")
		{
			forum.GET("/", h.forumMainPage)
		}

		view.GET("submission/:id", h.getSumbissionCode)

		profile := view.Group("/profile")
		{
			profile.GET("/", h.profileMainPage)
		}

		hwIu9Bmstu := view.Group("/hw_iu9_bmstu")
		{
			hwIu9Bmstu.GET("/", h.hwIu9MainPage)
		}

		view.GET("/timus", h.timusTaskList)
		view.GET("/acmp", h.acmpTaskList)

		problem := view.Group("/problem")
		{
			problem.GET("/:id", h.getTask)
			// TODO submitTask == submitContestTask
			problem.POST("/:id/submit", h.submitTask)
		}

		contest := view.Group("/contest/:contest_id")
		{
			contest.GET("/problems", h.getContestTasks)
			contest.GET("/problem/:problem_id", h.getTask)
			contest.POST("/problem/:problem_id/submit", h.submitContestTask)
		}

		view.GET("/groups", h.groups)
		groups := view.Group("/group")
		{
			groups.GET("/invite/:invite_hash", h.checkInvite)
			group := groups.Group("/:group_id")
			{
				group.GET("", h.getGroupContests)
				groupContest := group.Group("/contest/:contest_id")
				{
					groupContest.GET("/tasks", h.getContestTasks)
				}
			}
		}
	}

	return router
}

func nl2br(s string) template.HTML {
	return template.HTML(strings.ReplaceAll(s, "\n", "<br>"))
}
