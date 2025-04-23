package services

import (
	"auth-server/internal/models"
	"context"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}
func (s *UserService) GetUsersList(ctx context.Context) ([]models.Users, error) {
	var users []models.Users
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (s *UserService) CreateUserRequest(ctx context.Context, user *models.Users) error {
	return s.DB.Create(user).Error
}
func (s *UserService) UpdateUserByID(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.DB.Model(&models.Users{}).Where("id = ?", id).Updates(updates).Error
}
func (s *UserService) DeleteUserByID(ctx context.Context, id uint) error {
	return s.DB.Delete(&models.Users{}, id).Error
}
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.Users, error) {
	var user models.Users
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckRole verifies the user's role
func (s *UserService) CheckRole(ctx context.Context, id uint, required models.Role) (bool, error) {
	var user models.Users
	if err := s.DB.Select("role").Where("id = ?", id).First(&user).Error; err != nil {
		return false, err
	}
	return user.Role == string(required), nil
}
