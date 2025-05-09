package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"user/config"
	"user/db"
	"user/handlers"
	"user/models"
	"user/repository"
	"user/routes"
	"user/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.LoadConfig()
	db, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.User{})

	userRepository := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()
	routes.SetupRoutes(r, userHandler)

	r.Run(":8080")
}
