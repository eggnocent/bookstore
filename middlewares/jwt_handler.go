package middlewares

import (
	"apimandiri/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(config.JTWExpiredTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecretkey)
}
