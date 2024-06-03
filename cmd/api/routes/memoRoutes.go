package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/akinolaemmanuel49/memo-api/memo/api/handlers"
	"github.com/akinolaemmanuel49/memo-api/memo/api/internal"
	"github.com/akinolaemmanuel49/memo-api/memo/api/middleware"
)

func memoRoutes(app internal.Application, routes *gin.Engine) {
	memoHandler := handlers.NewMemoHandler(app)
	memo := routes.Group("/memo")
	memo.Use(middleware.Authentication(app), middleware.ContextUserSoftDelete())
	{
		memo.POST("/text", memoHandler.CreateTextMemo)
		memo.POST("/image", memoHandler.CreateImageMemo)
		memo.POST("/video", memoHandler.CreateVideoMemo)
		memo.POST("/audio", memoHandler.CreateAudioMemo)
		memo.GET("/:memoID", memoHandler.GetMemo)
		memo.DELETE("/:memoID", memoHandler.DeleteMemo)
		memo.POST("/like/:memoID", memoHandler.LikeMemo)
		memo.POST("/unlike/:memoID", memoHandler.UnlikeMemo)
		memo.POST("/share/:memoID", memoHandler.ShareMemo)
		memo.POST("/unshare/:memoID", memoHandler.UnshareMemo)
		memo.GET("/all", memoHandler.GetAllMemos)
		memo.GET("/feed", memoHandler.GetSubscribedMemos)
		memo.GET("/memos/:ownerID", memoHandler.GetMemosByOwnerID)
		memo.GET("/memos/me", memoHandler.GetOwnMemos)
	}
}
