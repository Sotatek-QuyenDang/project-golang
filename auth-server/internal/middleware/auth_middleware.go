package middleware

import (
	"auth-server/internal/config"
	"auth-server/internal/models"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var RedisClient *redis.Client

func InitMiddleware(redis *redis.Client) {
	RedisClient = redis
}
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Check Redis for token validity
		redisKey := "token:" + tokenStr
		exists, err := RedisClient.Exists(context.Background(), redisKey).Result()
		if err != nil || exists == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired or logged out"})
			return
		}

		claims := token.Claims.(*Claims)
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != string(models.RoleAdmin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}

func RequireUserOrAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || (role != string(models.RoleAdmin) && role != string(models.RoleUser)) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user or admin access required"})
			return
		}
		c.Next()
	}
}
