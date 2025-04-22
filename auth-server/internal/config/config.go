package config

import (
	"auth-server/internal/models"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB     DBConfig
	Redis  RedisConfig
	JWT    JWTConfig
	Server ServerConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	Expiration int
}

type ServerConfig struct {
	Port string
}

var AppConfig *Config
var Database *gorm.DB

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./internal/config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	AppConfig = &config
	return nil
}

func InitDatabase() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		AppConfig.DB.Host,
		AppConfig.DB.User,
		AppConfig.DB.Password,
		AppConfig.DB.Name,
		AppConfig.DB.Port,
		AppConfig.DB.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	models.Migrate(db)
	Database = db
	return nil
}
