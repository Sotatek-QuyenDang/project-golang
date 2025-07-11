package services_test

import (
	"auth-server/internal/models"
	"auth-server/internal/services"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUserRepoForTest struct {
	users []*models.Users
}

func (m *mockUserRepoForTest) GetAllUsers(ctx context.Context) ([]*models.Users, error) {
	return m.users, nil
}

func (m *mockUserRepoForTest) GetUserByID(ctx context.Context, id uint) (*models.Users, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}
func (m *mockUserRepoForTest) CreateUser(ctx context.Context, user *models.Users) error {
	m.users = append(m.users, user)
	return nil
}

func (m *mockUserRepoForTest) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	for i, user := range m.users {
		if user.ID == id {
			// Apply updates to the user
			for key, value := range updates {
				switch key {
				case "user_name":
					m.users[i].UserName = value.(string)
				case "hashed_password":
					m.users[i].HashedPassword = value.(string)
				case "role":
					m.users[i].Role = value.(string)
				}
			}
			return nil
		}
	}
	return nil
}

func (m *mockUserRepoForTest) DeleteUser(ctx context.Context, id uint) error {
	for i, user := range m.users {
		if user.ID == id {
			// Remove the user from the slice
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockUserRepoForTest) CheckRole(ctx context.Context, id uint) (string, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user.Role, nil
		}
	}
	return "", nil
}

func TestGetAllUsers(t *testing.T) {
	// Setup mock repository
	mockRepo := &mockUserRepoForTest{
		users: []*models.Users{
			{ID: 1, UserName: "user1", HashedPassword: "password1", Role: "user"},
			{ID: 2, UserName: "user2", HashedPassword: "password2", Role: "user"},
		},
	}

	// Create user service with mock
	userService := services.NewUserService(mockRepo)

	// Test GetAllUsers
	users, err := userService.GetAllUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := &mockUserRepoForTest{
		users: []*models.Users{
			{ID: 1, UserName: "user1", HashedPassword: "password1", Role: "user"},
			{ID: 2, UserName: "user2", HashedPassword: "password2", Role: "user"},
		},
	}
	userService := services.NewUserService(mockRepo)

	t.Run("User exists", func(t *testing.T) {
		user, err := userService.GetUserByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(1), user.ID)
	})

	t.Run("User not found", func(t *testing.T) {
		user, err := userService.GetUserByID(context.Background(), 999)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})
}

func TestCreateUser(t *testing.T) {
	mockRepo := &mockUserRepoForTest{}
	userService := services.NewUserService(mockRepo)

	newUser := &models.Users{
		ID:             3,
		UserName:       "newuser",
		HashedPassword: "password3",
		Role:           "admin",
	}

	err := userService.CreateUser(context.Background(), newUser)
	assert.NoError(t, err)
	assert.Len(t, mockRepo.users, 1)
	assert.Equal(t, "newuser", mockRepo.users[0].UserName)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := &mockUserRepoForTest{
		users: []*models.Users{
			{ID: 1, UserName: "user1", HashedPassword: "password1", Role: "user"},
		},
	}
	userService := services.NewUserService(mockRepo)
	updates := map[string]interface{}{
		"user_name": "updateduser",
		"role":      "admin",
	}
	err := userService.UpdateUser(context.Background(), 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", mockRepo.users[0].UserName)
	assert.Equal(t, "admin", mockRepo.users[0].Role)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := &mockUserRepoForTest{
		users: []*models.Users{
			{ID: 1, UserName: "user1", HashedPassword: "password1", Role: "user"},
		},
	}
	userService := services.NewUserService(mockRepo)
	err := userService.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, mockRepo.users, 0)
}

func TestCheckRole(t *testing.T) {
	mockRepo := &mockUserRepoForTest{
		users: []*models.Users{
			{ID: 1, UserName: "user1", HashedPassword: "password1", Role: "user"},
		},
	}
	userService := services.NewUserService(mockRepo)
	role, err := userService.CheckRole(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "user", role)
}
