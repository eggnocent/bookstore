package services

import (
	"apimandiri/models"
	"apimandiri/repositories"
	"errors"
)

type AuthService interface {
	Authenticate(username, password string) (models.User, error)
}

type authService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Authenticate(username, password string) (models.User, error) {
	var user models.User
	err := s.repo.FindByUsername(username, &user)
	if err != nil {
		return user, err
	}
	if !CheckPassword(user.Password, password) {
		return user, errors.New("Invalid Credentials")
	}
	return user, nil
}
