package main

import (
	"apimandiri/config"
	"apimandiri/controllers"
	"apimandiri/models"
	"apimandiri/repositories"
	"apimandiri/server"
	"apimandiri/services"
	"log"

	"gorm.io/gorm"
)

func main() {
	db := config.InitDB()
	if err := db.AutoMigrate(&models.User{}, &models.Penulis{}, &models.Buku{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// Tambahkan kode seeder di sini
	seedUser(db)

	// Inisialisasi repository, service, dan controller
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	bookRepo := repositories.NewBookRepository(db)
	bookService := services.NewBookService(bookRepo)
	bukuController := controllers.NewBukuController(bookService)

	penulisRepo := repositories.NewPenulisRepository(db)
	penulisService := services.NewPenulisService(penulisRepo)
	penulisController := controllers.NewPenulisController(penulisService)

	// Inisialisasi router dan jalankan server
	r := server.InitRouter(authController, userController, bukuController, penulisController)
	r.Run() // Menjalankan server pada port default (8080)
}

// Fungsi untuk menjalankan seeder user
func seedUser(db *gorm.DB) {
	// Cek apakah user sudah ada
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		password, err := services.EncryptPassword("123")
		if err != nil {
			log.Fatal("Failed to encrypt password: ", err)
		}

		user := models.User{
			Username: "adminuser",
			Password: password,
			Email:    "admin@example.com",
			FullName: "Admin User",
		}

		if err := db.Create(&user).Error; err != nil {
			log.Fatal("Failed to create user: ", err)
		}

		log.Println("Seeder berhasil dijalankan, user adminuser dibuat dengan password '123'")
	} else {
		log.Println("User sudah ada, seeder tidak dijalankan")
	}
}
