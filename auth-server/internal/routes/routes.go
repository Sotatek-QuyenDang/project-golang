package routes

import (
	"auth-server/internal/handlers"
	"auth-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) {
	r.POST("/api/login", authHandler.Login)
	r.POST("/api/logout", middleware.JWTAuthMiddleware(), authHandler.Logout)
	userRoutes := r.Group("/api/users")
	userRoutes.Use(middleware.JWTAuthMiddleware())
	{
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
		userRoutes.GET("/:id", userHandler.GetUserByID)
	}
}
