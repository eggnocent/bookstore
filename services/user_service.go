package services

import (
	"apimandiri/models"
	"apimandiri/repositories"
	"errors"
	"strconv"
)

type UserService interface {
	CreateUser(user models.User) error
	GetAllUsers(bukuID uint) ([]models.User, error)
	GetUserByID(id string) (models.User, error)
	UpdateUser(user models.User) error
	DeleteUser(id string) error
}

type userServiceImpl struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userServiceImpl{repo}
}

func (s *userServiceImpl) CreateUser(user models.User) error {
	encryptedPassword, err := EncryptPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = encryptedPassword

	return s.repo.Create(user)
}

func (s *userServiceImpl) GetAllUsers(bukuID uint) ([]models.User, error) {
	var users []models.User
	err := s.repo.FindAll(&users, bukuID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userServiceImpl) GetUserByID(id string) (models.User, error) {
	var user models.User

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return user, errors.New("Invalid user ID")
	}
	err = s.repo.FindByID(uint(userID), &user)
	if err != nil {
		return user, errors.New("ID Tidak Ditemukan")
	}
	return user, nil
}

func (s *userServiceImpl) UpdateUser(user models.User) error {
	var existingUser models.User

	err := s.repo.FindByID(user.ID, &existingUser)
	if err != nil {
		return errors.New("ID tidak ditemukan")
	}
	if user.Password != "" {
		EncryptedPassword, err := EncryptPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = EncryptedPassword
	}
	return s.repo.Update(user)
}

func (s *userServiceImpl) DeleteUser(id string) error {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.New("ID tidak valid")
	}

	var user models.User
	err = s.repo.FindByID(uint(userID), &user)
	if err != nil {
		return errors.New("ID tidak ditemukan")
	}

	// Panggil repository untuk menghapus user dan tangani errornya
	return s.repo.Delete(uint(userID))
}
