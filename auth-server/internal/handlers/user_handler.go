package handlers

import (
	"net/http"
	"strconv"

	"auth-server/internal/dto"
	"auth-server/internal/models"
	"auth-server/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Service.GetUsersList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := services.HashPassword(req.HashedPassword)
	user := models.Users{
		UserName:       req.UserName,
		HashedPassword: hashed,
		Role:           req.Role,
	}
	if err := h.Service.CreateUserRequest(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create user failed"})
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
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if req.HashedPassword != "" {
		hash, err := services.HashPassword(req.HashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		updates["hashed_password"] = hash
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if err := h.Service.UpdateUserByID(c.Request.Context(), uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update user failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Service.DeleteUserByID(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete user failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.Service.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Get user successfully",
		"data":    user,
	})
}
