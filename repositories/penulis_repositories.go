package repositories

import (
	"apimandiri/models"

	"gorm.io/gorm"
)

type PenulisRepository interface {
	CreatePenulis(penulis models.Penulis) error
	GetAllPenulis() ([]models.Penulis, error)
	GetPenulisByID(id uint) (models.Penulis, error)
	UpdatePenulis(penulis models.Penulis) error
	DeletePenulis(id uint) error
}

type penulisRepositoryImpl struct {
	db *gorm.DB
}

func NewPenulisRepository(db *gorm.DB) PenulisRepository {
	return &penulisRepositoryImpl{db}
}

func (r *penulisRepositoryImpl) CreatePenulis(penulis models.Penulis) error {
	return r.db.Create(&penulis).Error
}

func (r *penulisRepositoryImpl) GetAllPenulis() ([]models.Penulis, error) {
	var penulis []models.Penulis
	err := r.db.Preload("Buku").Find(&penulis).Error
	if err != nil {
		return nil, err
	}
	return penulis, nil
}

func (r *penulisRepositoryImpl) GetPenulisByID(id uint) (models.Penulis, error) {
	var penulis models.Penulis
	err := r.db.Preload("Buku").First(&penulis, id).Error
	if err != nil {
		return penulis, err
	}
	return penulis, nil

}

func (r *penulisRepositoryImpl) UpdatePenulis(penulis models.Penulis) error {
	return r.db.Save(&penulis).Error
}

func (r *penulisRepositoryImpl) DeletePenulis(id uint) error {
	return r.db.Delete(&models.Penulis{}, id).Error
}
