package middleware

import (
	"net/http"
	"reflect"

	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/akinolaemmanuel49/memo-api/domain/repository"
	"github.com/akinolaemmanuel49/memo-api/memo/api/helpers"
	"github.com/gin-gonic/gin"
)

// ContextUserSoftDelete checks if the instance of user
// in context has Deleted set to TRUE
func ContextUserSoftDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := helpers.ContextGetUser(ctx)
		if reflect.DeepEqual(user, models.User{}) {
			helpers.HandleErrorResponse(ctx, http.StatusNotFound, repository.ErrRecordNotFound)
			return
		}
		if user.Deleted {
			helpers.HandleErrorResponse(ctx, http.StatusNoContent, repository.ErrRecordDeleted)
		}
		ctx.Next()
	}
}
