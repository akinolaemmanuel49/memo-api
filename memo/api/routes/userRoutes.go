package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/akinolaemmanuel49/memo-api/memo/api/handlers"
	"github.com/akinolaemmanuel49/memo-api/memo/api/internal"
	"github.com/akinolaemmanuel49/memo-api/memo/api/middleware"
)

func userRoutes(app internal.Application, routes *gin.Engine) {
	userHandler := handlers.NewUserHandler(app)
	user := routes.Group("/users")
	user.Use(middleware.Authentication(app), middleware.ContextUserSoftDelete())
	{
		user.GET("", userHandler.Get)
		user.PUT("", userHandler.Update)
		user.DELETE("", userHandler.Delete)
		user.GET("/followers", userHandler.GetFollowers)
		user.GET("/following", userHandler.GetFollowing)
		user.DELETE("/avatar", userHandler.DeleteAvatar)
	}
}
