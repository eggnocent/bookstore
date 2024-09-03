package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, statusCode int, message string) {
	log.Printf("Error terjadi: %s (statusCode %d)", message, statusCode)

	ctx.JSON(statusCode, gin.H{"error": message})
	ctx.Abort()
}

func HandleValidationError(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusBadRequest, message)
}

func HandleNotFoundError(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusNotFound, message)
}

func HandleAuthError(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusUnauthorized, message)
}

func HandleForbiddenError(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusForbidden, message)
}

func HandleInternalServerError(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusConflict, message)
}

func HandleMethodNotAllowedError(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusMethodNotAllowed, message)
}
