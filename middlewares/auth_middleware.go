package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var Sessions = make(map[string]time.Time)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID := ctx.GetHeader("Session-ID")
		if sessionID == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Session-ID dibutuhkan"})
			ctx.Abort()
			return
		}

		lastActivity, exists := Sessions[sessionID]
		if !exists || time.Since(lastActivity) > 5*time.Minute {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi Anda habis atau tidak ada"})
			ctx.Abort()
			return
		}

		Sessions[sessionID] = time.Now()

		ctx.Next()
	}
}
