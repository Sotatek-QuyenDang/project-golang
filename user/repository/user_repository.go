package repository

import (
	"gorm.io/gorm"
	"user/models"
)

type UserRepository interface {
	GetAll() ([]models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(user models.User) (models.User, error)
	GetByID(id int) (models.User, error)
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (r *userRepository) Create(user models.User) (models.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (r *userRepository) Update(user models.User) (models.User, error) {
	if err := r.db.Model(&user).Updates(user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (r *userRepository) Delete(user models.User) (models.User, error) {
	if err := r.db.Delete(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (r *userRepository) GetByID(id int) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
