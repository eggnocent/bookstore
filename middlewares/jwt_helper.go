package middlewares

import (
	"apimandiri/config"
	"errors"
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

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing tida sesuai")
		}
		return config.JWTSecretkey, nil
	})
}
