package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Users struct {
	ID             uint      `gorm:"primaryKey"`
	UserName       string    `gorm:"uniqueIndex;not null"`
	HashedPassword string    `gorm:"not null"`
	Role           string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (Users) TableName() string {
	return "users"
}
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Users{})
}
