package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// HandleErrorResponse returns an error with the specified status code and message as the response.
func HandleErrorResponse(ctx *gin.Context, statusCode int, err error) {
	ctx.AbortWithStatusJSON(statusCode, gin.H{"error": err.Error()})
}

// HandleInternalServerError logs the error and sends a generic 500 error response to the client.
func HandleInternalServerError(ctx *gin.Context, err error) {
	log.Printf("internal server error: %s", err.Error())
	HandleErrorResponse(ctx, http.StatusInternalServerError, errors.New("internal server error"))
}

// HandleValidationError is a shortcut to HandleErrorResponse with a status of 422.
// It should be used for all validation errors.
func HandleValidationError(ctx *gin.Context, err error) {
	HandleErrorResponse(ctx, http.StatusUnprocessableEntity, err)
}

// HandleLogicalDeleteError returns an error with status of 404 and the data of the deleted instance as the response.
func HandleLogicalDeleteError(ctx *gin.Context, data interface{}, err error) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error(), "data": data})
}
