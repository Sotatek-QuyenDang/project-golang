package repository

import (
	"auth-server/internal/models"
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type userRepositoryGorm struct {
	DB *gorm.DB
}

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]*models.Users, error)
	GetUserByID(ctx context.Context, id uint) (*models.Users, error)
	CreateUser(ctx context.Context, user *models.Users) error
	UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, id uint) error
	CheckRole(ctx context.Context, id uint) (string, error)
}

// NewUserRepositoryGorm creates a new GORM implementation of UserRepository
func NewUserRepositoryGorm(db *gorm.DB) UserRepository {
	return &userRepositoryGorm{DB: db}
}

func (r *userRepositoryGorm) GetAllUsers(ctx context.Context) ([]*models.Users, error) {
	log.Println("Repository: Getting all users")
	var users []*models.Users
	err := r.DB.WithContext(ctx).Find(&users).Error
	if err != nil {
		log.Printf("Repository: Error getting all users: %v", err)
		return nil, err
	}
	log.Printf("Repository: Found %d users", len(users))
	return users, err
}

func (r *userRepositoryGorm) GetUserByID(ctx context.Context, id uint) (*models.Users, error) {
	log.Printf("Repository: Getting user with ID: %d", id)
	var user models.Users
	err := r.DB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		log.Printf("Repository: Error getting user with ID %d: %v", id, err)
		return nil, err
	}
	log.Printf("Repository: Found user: %s, role: %s", user.UserName, user.Role)
	return &user, err
}

func (r *userRepositoryGorm) CreateUser(ctx context.Context, user *models.Users) error {
	log.Printf("Repository: Creating user with username: %s, role: %s", user.UserName, user.Role)
	err := r.DB.WithContext(ctx).Create(user).Error
	if err != nil {
		log.Printf("Repository: Error creating user: %v", err)
		return err
	}
	log.Printf("Repository: User created successfully with ID: %d", user.ID)
	return nil
}

func (r *userRepositoryGorm) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	log.Printf("Repository: Updating user with ID: %d, updates: %+v", id, updates)
	result := r.DB.WithContext(ctx).Model(&models.Users{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		log.Printf("Repository: Error updating user with ID %d: %v", id, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("Repository: No rows affected when updating user with ID %d", id)
		return fmt.Errorf("user with ID %d not found", id)
	}
	log.Printf("Repository: User with ID %d updated successfully", id)
	return nil
}

func (r *userRepositoryGorm) DeleteUser(ctx context.Context, id uint) error {
	log.Printf("Repository: Deleting user with ID: %d", id)
	result := r.DB.WithContext(ctx).Delete(&models.Users{}, id)
	if result.Error != nil {
		log.Printf("Repository: Error deleting user with ID %d: %v", id, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("Repository: No rows affected when deleting user with ID %d", id)
		return fmt.Errorf("user with ID %d not found", id)
	}
	log.Printf("Repository: User with ID %d deleted successfully", id)
	return nil
}

func (r *userRepositoryGorm) CheckRole(ctx context.Context, id uint) (string, error) {
	log.Printf("Repository: Checking role for user with ID: %d", id)
	var user models.Users
	err := r.DB.WithContext(ctx).Select("role").Where("id = ?", id).First(&user).Error
	if err != nil {
		log.Printf("Repository: Error checking role for user with ID %d: %v", id, err)
		return "", err
	}
	log.Printf("Repository: User with ID %d has role: %s", id, user.Role)
	return user.Role, nil
}
