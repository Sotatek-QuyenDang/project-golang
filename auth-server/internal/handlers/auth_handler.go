package handlers

import (
	"fmt"
	"net/http"

	"auth-server/internal/dto"
	"auth-server/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	if service == nil {
		panic("auth service cannot be nil")
	}
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}
	if req.UserName == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserName and HashedPassword are required"})
		return
	}
	fmt.Printf("Dữ liệu nhận được: %+v\n", req)
	resp, err := h.Service.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("login failed: %v", err)})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	if len(token) <= 7 || token[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		return
	}

	token = token[7:] // trim "Bearer " prefix
	if err := h.Service.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("logout failed: %v", err)})
		return
	}
	c.Status(http.StatusOK)
}
