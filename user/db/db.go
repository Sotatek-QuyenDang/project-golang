package db

import (
	"fmt"
	"user/config"
	model "user/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Sql struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
}

func Connect(cfg config.Config) (*gorm.DB, error) {
	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fmt.Println("DSN:", dsn)
	fmt.Println("Error:", err)
	return db, err
}
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.User{})
}
