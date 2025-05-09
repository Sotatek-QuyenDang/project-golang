package seed

import (
	"log"

	"auth-server/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	users := []struct {
		UserName string
		Password string
		Role     string
	}{
		{"admin01", "admin123", "admin"},
		{"user01", "user123", "user"},
		{"user02", "user456", "user"},
	}

	for _, u := range users {
		// Kiểm tra user đã tồn tại chưa
		var existing models.Users
		if err := db.Where("user_name = ?", u.UserName).First(&existing).Error; err == nil {
			log.Printf("User %s đã tồn tại, bỏ qua", u.UserName)
			continue
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Không thể hash password cho %s: %v", u.UserName, err)
			continue
		}

		user := models.Users{
			UserName:       u.UserName,
			HashedPassword: string(hashed),
			Role:           u.Role,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Không thể tạo user %s: %v", u.UserName, err)
		} else {
			log.Printf("Đã thêm user %s thành công", u.UserName)
		}
	}
}
