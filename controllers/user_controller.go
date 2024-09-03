package controllers

import (
	"apimandiri/middlewares"
	"apimandiri/models"
	"apimandiri/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userController struct {
	service services.UserService
}

func NewUserController(services services.UserService) UserController {
	return &userController{services}
}

type UserResponse struct {
	ID       uint      `json:"ID"`
	Username string    `json:"Username"`
	Password string    `json:"Password"`
	Email    string    `json:"Email"`
	FullName string    `json:"FullName"`
	Buku     string    `json:"Buku"`
	CreateAt time.Time `json:"CreateAt"`
	UpdateAt time.Time `json:"UpdateAt"`
}

type UserResponseByID struct {
	ID       uint      `json:"ID"`
	Username string    `json:"Username"`
	Password string    `json:"Password"`
	Email    string    `json:"Email"`
	FullName string    `json:"FullName"`
	Buku     string    `json:"Buku"`
	CreateAt time.Time `json:"CreateAt"`
	UpdateAt time.Time `json:"UpdateAt"`
}

func (c *userController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	if user.FullName == "" || user.Username == "" || user.Password == "" || user.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "harus mengisi semua kolom"})
		return
	}
	if err := c.service.CreateUser(user); err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"pesan": "User Berhasil DIbuat"})
}

func (c *userController) GetAllUsers(ctx *gin.Context) {
	// Ambil parameter query buku dari URL
	bukuIDStr := ctx.Query("buku")

	var bukuID uint
	var err error

	// Periksa apakah bukuIDStr ada dan konversikan ke uint
	if bukuIDStr != "" {
		bukuIDUint, err := strconv.ParseUint(bukuIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID buku tidak valid"})
			return
		}
		bukuID = uint(bukuIDUint)
	}

	// Panggil service dengan bukuID, jika bukuID 0 maka akan mengambil semua data
	users, err := c.service.GetAllUsers(bukuID)
	if err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
		response := UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
			FullName: user.FullName,
			CreateAt: user.CreateAt,
			UpdateAt: user.UpdateAt,
			Buku:     "",
		}

		// Jika user memiliki buku, ambil nama bukunya
		if user.Buku != nil {
			response.Buku = user.Buku.NamaBuku
		}

		userResponses = append(userResponses, response)
	}

	ctx.JSON(http.StatusOK, userResponses)
}

func (c *userController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.service.GetUserByID(id)
	if err != nil {
		middlewares.HandleNotFoundError(ctx, err.Error())
		return
	}

	// Map Buku dan Penulis ke response
	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		FullName: user.FullName,
		Buku:     "",
		CreateAt: user.CreateAt,
		UpdateAt: user.UpdateAt,
	}
	if user.Buku != nil {
		userResponse.Buku = user.Buku.NamaBuku
	}

	ctx.JSON(http.StatusOK, userResponse)
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	_, err = c.service.GetUserByID(id)
	if err != nil {
		middlewares.HandleNotFoundError(ctx, "id tidak ditemukan")
		return
	}
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = uint(userID)

	if err := c.service.UpdateUser(user); err != nil {
		middlewares.HandleInternalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"pesan": "User telah di update"})
}

func (c *userController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.DeleteUser(id)
	if err != nil {
		// Gunakan fungsi dari error_handle.go untuk menampilkan pesan error spesifik
		if err.Error() == "user tidak dapat dihapus karena memiliki relasi dengan tabel lain" {
			middlewares.HandleInternalServerError(ctx, err.Error())
		} else {
			middlewares.HandleInternalServerError(ctx, "Terjadi kesalahan saat menghapus user")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"pesan": "User berhasil dihapus"})
}
