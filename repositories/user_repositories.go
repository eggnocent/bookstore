package repositories

import (
	"apimandiri/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user models.User) error
	FindAll(users *[]models.User, bukuID uint) error
	FindByID(id uint, user *models.User) error
	Update(user models.User) error
	FindByUsername(username string, user *models.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) FindAll(users *[]models.User, bukuID uint) error {
	// Jika bukuID tidak 0, lakukan filter berdasarkan buku
	if bukuID != 0 {
		// Ambil buku yang sesuai dengan bukuID
		var buku models.Buku
		if err := r.db.Where("id = ?", bukuID).First(&buku).Error; err != nil {
			return err // Kembalikan error jika buku tidak ditemukan
		}

		// Ambil pengguna yang memiliki ID yang sama dengan UserID dari buku
		err := r.db.Where("id = ?", buku.UserID).Find(users).Error
		return err
	}

	// Jika tidak ada filter bukuID, ambil semua pengguna
	err := r.db.Preload("Buku").Find(users).Error
	return err
}

func (r *userRepository) FindByID(id uint, user *models.User) error {
	return r.db.Preload("Buku").First(user, id).Error
}

func (r *userRepository) Update(user models.User) error {
	return r.db.Save(&user).Error
}

func (r *userRepository) FindByUsername(username string, user *models.User) error {
	return r.db.Where("username = ?", username).First(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.User{}, id).Error
}
