package services

import (
	"auth-server/internal/models"
	"auth-server/internal/repository"
	"context"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*models.Users, error)
	GetUserByID(ctx context.Context, id uint) (*models.Users, error)
	CreateUser(ctx context.Context, user *models.Users) error
	UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, id uint) error
	CheckRole(ctx context.Context, id uint) (string, error)
}

type userService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{Repo: repo}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*models.Users, error) {
	return s.Repo.GetAllUsers(ctx)
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.Users, error) {
	return s.Repo.GetUserByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, user *models.Users) error {
	return s.Repo.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.Repo.UpdateUser(ctx, id, updates)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.Repo.DeleteUser(ctx, id)
}

func (s *userService) CheckRole(ctx context.Context, id uint) (string, error) {
	role, err := s.Repo.CheckRole(ctx, id)
	if err != nil {
		return "", err
	}
	return role, nil
}
