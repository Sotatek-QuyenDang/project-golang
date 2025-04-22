// services/user_service.go
package services

import (
	"user/models"
	"user/repository"
)

type UserService interface {
	GetAll() ([]models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(user models.User) (models.User, error)
	GetByID(id int) (models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) Create(user models.User) (models.User, error) {
	return s.repo.Create(user)
}

func (s *userService) Update(user models.User) (models.User, error) {
	return s.repo.Update(user)
}

func (s *userService) Delete(user models.User) (models.User, error) {
	return s.repo.Delete(user)
}

func (s *userService) GetByID(id int) (models.User, error) {
	return s.repo.GetByID(id)
}
