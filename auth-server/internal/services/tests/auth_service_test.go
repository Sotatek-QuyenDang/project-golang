package services_test

import (
	"auth-server/internal/config"
	"auth-server/internal/dto"
	"auth-server/internal/models"
	"auth-server/internal/services"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	users []*models.Users
}

func (m *mockUserRepository) GetAllUsers(ctx context.Context) ([]*models.Users, error) {
	return m.users, nil
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id uint) (*models.Users, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *models.Users) error {
	m.users = append(m.users, user)
	return nil
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	return nil
}

func (m *mockUserRepository) DeleteUser(ctx context.Context, id uint) error {
	return nil
}

func (m *mockUserRepository) CheckRole(ctx context.Context, id uint) (string, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user.Role, nil
		}
	}
	return "", nil
}

// Setup Redis mock server
func setupRedis(t *testing.T) (*redis.Client, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, func() {
		client.Close()
		mr.Close()
	}
}

func TestLogin(t *testing.T) {
	// Initialize config for JWT
	config.AppConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret-key",
			Expiration: 15,
		},
	}

	// Setup mock repository with a test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	testUser := &models.Users{
		ID:             1,
		UserName:       "admin01",
		HashedPassword: string(hashedPassword),
		Role:           "admin",
	}

	userRepo := &mockUserRepository{
		users: []*models.Users{testUser},
	}

	// Setup Redis mock
	redisMock, cleanup := setupRedis(t)
	defer cleanup()

	// Create auth service with mocks
	authService := services.NewAuthService(userRepo, redisMock)

	t.Run("Invalid password", func(t *testing.T) {
		loginRequest := &dto.LoginRequest{
			UserName: "admin01",
			Password: "wrongpassword",
		}

		_, err := authService.Login(context.Background(), loginRequest)
		assert.Error(t, err)
		assert.Equal(t, "invalid username or password", err.Error())
	})

	t.Run("Invalid username", func(t *testing.T) {
		loginRequest := &dto.LoginRequest{
			UserName: "nonexistent",
			Password: "testpassword",
		}

		_, err := authService.Login(context.Background(), loginRequest)
		assert.Error(t, err)
		assert.Equal(t, "invalid username or password", err.Error())
	})

	t.Run("Successful login", func(t *testing.T) {
		loginRequest := &dto.LoginRequest{
			UserName: "admin01",
			Password: "testpassword",
		}

		resp, err := authService.Login(context.Background(), loginRequest)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)

		// Verify token
		token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.AppConfig.JWT.Secret), nil
		})
		assert.NoError(t, err)

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, float64(1), claims["user_id"])
		assert.Equal(t, "admin", claims["role"])
	})
}

func TestLogout(t *testing.T) {
	// Setup Redis mock
	redisMock, cleanup := setupRedis(t)
	defer cleanup()

	// Create auth service with mocks
	authService := services.NewAuthService(&mockUserRepository{}, redisMock)

	// Set a token in Redis
	ctx := context.Background()
	testToken := "test-token"
	err := redisMock.Set(ctx, "token:"+testToken, 1, time.Minute).Err()
	assert.NoError(t, err)

	// Test logout
	err = authService.Logout(ctx, testToken)
	assert.NoError(t, err)

	// Verify token was removed
	val, err := redisMock.Get(ctx, "token:"+testToken).Result()
	assert.Error(t, err)
	assert.Empty(t, val)
}
