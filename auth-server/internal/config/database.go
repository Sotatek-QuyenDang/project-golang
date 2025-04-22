package config

import (
	"auth-server/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() error {
	// Kiểm tra nếu AppConfig.DB.Port là kiểu số (int), chuyển thành chuỗi
	portStr := fmt.Sprintf("%d", AppConfig.DB.Port)

	// Kiểm tra các giá trị cấu hình
	if AppConfig.DB.Host == "" || AppConfig.DB.User == "" || AppConfig.DB.Password == "" || AppConfig.DB.Name == "" || portStr == "" || AppConfig.DB.SSLMode == "" {
		return fmt.Errorf("invalid database configuration: missing one or more required fields")
	}

	// Xây dựng chuỗi kết nối DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		AppConfig.DB.Host, AppConfig.DB.User, AppConfig.DB.Password, AppConfig.DB.Name, portStr, AppConfig.DB.SSLMode)

	// Kết nối với database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	models.Migrate(db)
	Database = db
	return nil
}
