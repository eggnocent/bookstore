package controllers

import (
	"apimandiri/middlewares"
	"apimandiri/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type authController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) AuthController {
	return &authController{service}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := c.service.Authenticate(loginData.Username, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	sessionID := user.Username
	middlewares.Sessions[sessionID] = time.Now()

	ctx.JSON(http.StatusOK, gin.H{"pesan": "Login sukses"})
}

func (c *authController) Logout(ctx *gin.Context) {
	sessionID := ctx.GetHeader("Session-ID")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Session-ID dibutuhkan"})
		return
	}
	delete(middlewares.Sessions, sessionID)
	ctx.JSON(http.StatusOK, gin.H{"pesan": "Logout sukses"})
}
