package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/akinolaemmanuel49/memo-api/cmd/api/helpers"
	"github.com/akinolaemmanuel49/memo-api/cmd/api/internal"
)

func Router(app internal.Application) *gin.Engine {
	gin.EnableJsonDecoderDisallowUnknownFields()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           helpers.MaxAge,
	}))

	// set routes
	authRoutes(app, router)
	userRoutes(app, router)
	socialRoutes(app, router)
	memoRoutes(app, router)
	return router
}
