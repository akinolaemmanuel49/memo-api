package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/akinolaemmanuel49/memo-api/cmd/api/handlers"
	"github.com/akinolaemmanuel49/memo-api/cmd/api/internal"
)

func authRoutes(app internal.Application, routes *gin.Engine) {
	authHandler := handlers.NewAuthHandler(app)
	auth := routes.Group("/auth")
	{
		auth.POST("/signup", authHandler.SignUp)
		auth.POST("/token", authHandler.Token)
		auth.POST("/refresh-token", authHandler.RefreshToken)
	}
}
