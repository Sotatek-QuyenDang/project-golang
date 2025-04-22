package routes

import (
	"github.com/gin-gonic/gin"
	"user/handlers"
)

func SetupRoutes(r *gin.Engine, h *handlers.UserHandler) {
	v1 := r.Group("/api")
	{
		v1.GET("/users", h.GetUsers)
		v1.POST("/users", h.CreateUser)
		v1.PUT("/users/:id", h.UpdateUser)
		v1.DELETE("/users/:id", h.DeleteUser)
		v1.GET("/users/:id", h.GetUserById)
	}
}
