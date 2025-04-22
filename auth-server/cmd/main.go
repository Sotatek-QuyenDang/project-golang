package main

import (
	"context"
	"log"

	"auth-server/internal/config"
	"auth-server/internal/handlers"
	"auth-server/internal/middleware"
	"auth-server/internal/routes"
	"auth-server/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	middleware.InitMiddleware(rdb)

	authService := services.NewAuthService(config.Database, rdb)
	userService := services.NewUserService(config.Database)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()
	routes.SetupRoutes(r, authHandler, userHandler)

	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
