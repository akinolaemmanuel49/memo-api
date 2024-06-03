package helpers

import (
	"github.com/akinolaemmanuel49/memo-api/domain/models"
	"github.com/gin-gonic/gin"
)

type ContextKey string

const UserContextKey = ContextKey("user")

// ContextSetUser saves the given user data in the request context.
func ContextSetUser(ctx *gin.Context, user models.User) {
	ctx.Set(string(UserContextKey), user)
}

// ContextGetUser returns the current user data stored in the request context.
func ContextGetUser(ctx *gin.Context) models.User {
	user, ok := ctx.Value(string(UserContextKey)).(models.User)
	if !ok {
		return models.User{}
	}
	return user
}
