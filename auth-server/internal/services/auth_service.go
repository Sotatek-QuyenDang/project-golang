package services

import (
	"auth-server/internal/config"
	"auth-server/internal/dto"
	"auth-server/internal/models"
	"auth-server/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
	Redis    *redis.Client
}

func NewAuthService(userRepo repository.UserRepository, Redis *redis.Client) *AuthService {
	return &AuthService{userRepo: userRepo, Redis: Redis}
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (dto.LoginResponse, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	var user *models.Users
	for _, u := range users {
		if u.UserName == req.UserName {
			user = u
			break
		}
	}

	if user == nil {
		return dto.LoginResponse{}, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		return dto.LoginResponse{}, errors.New("invalid username or password")
	}

	exp := time.Minute * time.Duration(config.AppConfig.JWT.Expiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(exp).Unix(),
	})

	secret := []byte(config.AppConfig.JWT.Secret)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	err = s.Redis.Set(ctx, "token:"+signedToken, user.ID, exp).Err()
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{Token: signedToken}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.Redis.Del(ctx, "token:"+token).Err()
}
