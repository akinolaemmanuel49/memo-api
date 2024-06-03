package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/akinolaemmanuel49/memo-api/cmd/api/handlers"
	"github.com/akinolaemmanuel49/memo-api/cmd/api/internal"
	"github.com/akinolaemmanuel49/memo-api/cmd/api/middleware"
)

func socialRoutes(app internal.Application, routes *gin.Engine) {
	socialHandler := handlers.NewSocialHandler(app)
	social := routes.Group("/social")
	social.Use(middleware.Authentication(app), middleware.ContextUserSoftDelete())
	{
		social.POST("/follow", socialHandler.Follow)
		social.POST("/unfollow", socialHandler.Unfollow)
		social.POST("/comment/:memoID", socialHandler.CreateTextComment)
		social.POST("/comment/reply/:memoID/:parentID", socialHandler.CreateTextReply)
		social.GET("/reply/:commentID/replies", socialHandler.GetReplies)
		social.GET("/comment/:memoID", socialHandler.GetComments)
	}
}
