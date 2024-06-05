package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/akinolaemmanuel49/memo-api/cmd/api/handlers"
	"github.com/akinolaemmanuel49/memo-api/cmd/api/internal"
	"github.com/akinolaemmanuel49/memo-api/cmd/api/middleware"
)

func userRoutes(app internal.Application, routes *gin.Engine) {
	userHandler := handlers.NewUserHandler(app)
	user := routes.Group("/users")
	user.Use(middleware.Authentication(app), middleware.ContextUserSoftDelete())
	{
		user.GET("", userHandler.Me)
		user.GET(":userID", userHandler.GetById)
		user.PUT("", userHandler.Update)
		user.DELETE("", userHandler.Delete)
		user.GET("/all", userHandler.GetAll)
		user.GET("/find", userHandler.Find)
		user.GET("/followers", userHandler.GetFollowers)
		user.GET("/following", userHandler.GetFollowing)
		user.DELETE("/avatar", userHandler.DeleteAvatar)
	}
}
