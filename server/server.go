// server/server.go
package server

import (
	"apimandiri/controllers"
	"apimandiri/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter(authController controllers.AuthController, userController controllers.UserController, bukuController controllers.BukuController, penulisController controllers.PenulisController) *gin.Engine {
	r := gin.Default()

	r.POST("/login", authController.Login)
	r.POST("/logout", authController.Logout)

	protected := r.Group("/")
	protected.Use(middlewares.JWTMiddleware())
	protected.POST("/users", userController.CreateUser)
	protected.GET("/users", userController.GetAllUsers)
	protected.GET("/users/:id", userController.GetUserByID)
	protected.PUT("/users/:id", userController.UpdateUser)
	protected.DELETE("/users/:id", userController.DeleteUser)

	// Rute untuk buku
	protected.POST("/users/:id/buku", bukuController.AddUserBook)
	protected.GET("/buku", bukuController.GetAllBooks)
	protected.GET("/buku/:id", bukuController.GetBookByID)
	protected.GET("/users/:id/buku", bukuController.GetUserBook)
	protected.PUT("/users/:id/buku", bukuController.UpdateUserBook)
	protected.PUT("/buku/:id", bukuController.UpdateBookByID)
	protected.DELETE("/users/:id/buku", bukuController.DeleteUserBook)
	protected.DELETE("/buku/:id", bukuController.DeleteBookByID)

	// Rute untuk many-to-many dengan path tambahan
	protected.POST("/manage/buku/:book_id/authors/:author_id", bukuController.AddAuthorToBook)
	protected.PUT("/manage/buku/:book_id/authors", bukuController.UpdateAuthorsForBook)               // Update penulis untuk buku
	protected.DELETE("/manage/buku/:book_id/authors/:author_id", bukuController.DeleteAuthorFromBook) // Hapus penulis dari buku

	protected.POST("/penulis", penulisController.CreatePenulis)
	protected.GET("/penulis", penulisController.GetAllPenulis)
	protected.GET("/penulis/:id", penulisController.GetPenulisByID)
	protected.PUT("/penulis/:id", penulisController.UpdatePenulis)
	protected.DELETE("/penulis/:id", penulisController.DeletePenulis)

	return r
}
