package services

import (
	"auth-server/internal/models"
	"auth-server/internal/repository"
	"context"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}
func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.Users, error) {
	return s.Repo.GetAllUsers(ctx)
}
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.Users, error) {
	return s.Repo.GetUserByID(ctx, id)
}
func (s *UserService) CreateUser(ctx context.Context, user *models.Users) error {
	return s.Repo.CreateUser(ctx, user)
}
func (s *UserService) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.Repo.UpdateUser(ctx, id, updates)
}
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.Repo.DeleteUser(ctx, id)
}
func (s *UserService) CheckRole(ctx context.Context, id uint) (string, error) {
	role, err := s.Repo.CheckRole(ctx, id)
	if err != nil {
		return "", err
	}
	return role, nil
}
