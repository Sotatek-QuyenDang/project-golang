package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"auth-server/internal/dto"
	"auth-server/internal/models"
	"auth-server/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	if service == nil {
		panic("user service cannot be nil")
	}
	return &UserHandler{Service: *service}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get users: %v", err)})
		return
	}

	var resp []dto.UserResponse
	for _, user := range users {
		resp = append(resp, dto.UserResponse{
			ID:       user.ID,
			UserName: user.UserName,
			Role:     string(user.Role),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Get users successfully",
		"data":    resp,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if req.UserName == "" || req.HashedPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	hashed, err := services.HashPassword(req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to hash password: %v", err)})
		return
	}

	user := models.Users{
		UserName:       req.UserName,
		HashedPassword: hashed,
		Role:           req.Role,
	}
	if err := h.Service.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("create user failed: %v", err)})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data": dto.UserResponse{
			ID:       user.ID,
			UserName: user.UserName,
			Role:     user.Role,
		},
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid ID format: %v", err)})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	updates := map[string]interface{}{}
	if req.HashedPassword != "" {
		hash, err := services.HashPassword(req.HashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to hash password: %v", err)})
			return
		}
		updates["hashed_password"] = hash
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := h.Service.UpdateUser(c.Request.Context(), uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("update user failed: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid ID format: %v", err)})
		return
	}

	if err := h.Service.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("delete user failed: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid ID format: %v", err)})
		return
	}

	user, err := h.Service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get user: %v", err)})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get user successfully",
		"data": dto.UserResponse{
			ID:       user.ID,
			UserName: user.UserName,
			Role:     user.Role,
		},
	})
}
