package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/akinolaemmanuel49/memo-api/domain/repository"
	"github.com/akinolaemmanuel49/memo-api/memo/api/helpers"
	"github.com/akinolaemmanuel49/memo-api/memo/api/internal"
	"github.com/gin-gonic/gin"
)

// Authentication validates the provided access token and authenticates users.
func Authentication(app internal.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Vary", "Authorization")

		// check for the existence of the Authorization header
		authorizationHeader := ctx.GetHeader("Authorization")
		if authorizationHeader == "" {
			ctx.Header("WWW-Authenticate", "Bearer")
			helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("no Authorization header provided"))
			return
		}

		// return an error for a malformed token format
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			ctx.Header("WWW-Authenticate", "Bearer")
			helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("malformed token format"))
			return
		}

		// validate access token
		claims, err := helpers.ValidateToken(app.Config.JWTSecret, headerParts[1])
		if err != nil {
			ctx.Header("WWW-Authenticate", "Bearer")
			helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, err)
			return
		}

		// retrieve associated user
		user, err := app.Repositories.Users.GetById(claims.Subject)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				helpers.HandleErrorResponse(ctx, http.StatusUnauthorized, errors.New("user not found"))

			case errors.Is(err, repository.ErrRecordDeleted):
				helpers.HandleErrorResponse(ctx, http.StatusNotFound, errors.New("user has been deleted"))

			default:
				helpers.HandleInternalServerError(ctx, err)
			}

			return
		}

		// set user in context for further use
		helpers.ContextSetUser(ctx, user)
		ctx.Next()
	}
}
