package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			ctx.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			ctx.Abort()
			return
		}
		token, err := VerifyJWT(tokenString)
		if err != nil || token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid atau kedaluarsa"})
			ctx.Abort()
			return
		}
	}
}
