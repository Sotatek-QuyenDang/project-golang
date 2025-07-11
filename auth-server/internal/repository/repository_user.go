package repository

import (
	"auth-server/internal/models"
	"context"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]*models.Users, error)
	GetUserByID(ctx context.Context, id uint) (*models.Users, error)
	CreateUser(ctx context.Context, user *models.Users) error
	UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, id uint) error
	CheckRole(ctx context.Context, id uint) (string, error)
}
