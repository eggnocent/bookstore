// middlewares/jwt_middleware.go
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware is a middleware to protect routes and verify JWT
func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := VerifyJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			ctx.Abort()
			return
		}

		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}
