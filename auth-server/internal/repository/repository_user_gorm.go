package repository

import (
	"auth-server/internal/models"
	"context"
	"gorm.io/gorm"
)

type userRepositoryGorm struct {
	DB *gorm.DB
}

// NewUserRepositoryGorm creates a new GORM implementation of UserRepository
func NewUserRepositoryGorm(db *gorm.DB) UserRepository {
	return &userRepositoryGorm{DB: db}
}

func (r *userRepositoryGorm) GetAllUsers(ctx context.Context) ([]*models.Users, error) {
	var users []*models.Users
	err := r.DB.WithContext(ctx).Find(&users).Error
	return users, err
}
func (r *userRepositoryGorm) GetUserByID(ctx context.Context, id uint) (*models.Users, error) {
	var user models.Users
	err := r.DB.WithContext(ctx).First(&user, id).Error
	return &user, err
}
func (r *userRepositoryGorm) CreateUser(ctx context.Context, user *models.Users) error {
	return r.DB.WithContext(ctx).Create(user).Error
}
func (r *userRepositoryGorm) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.DB.WithContext(ctx).Model(&models.Users{}).Where("id = ?", id).Updates(updates).Error
}
func (r *userRepositoryGorm) DeleteUser(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&models.Users{}, id).Error
}
func (r *userRepositoryGorm) CheckRole(ctx context.Context, id uint) (string, error) {
	var user models.Users
	err := r.DB.WithContext(ctx).Select("role").Where("id = ?", id).First(&user).Error
	return user.Role, err
}
