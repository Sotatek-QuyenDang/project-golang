package services_test

import (
	"auth-server/internal/models"
	"auth-server/internal/repository"
	"auth-server/internal/services"
	"auth-server/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMockUserRepository(t *testing.T) (repository.UserRepository, func()) {
	mockRepo := new(mocks.UserRepository)

	cleanup := func() {
		mockRepo.AssertExpectations(t)
	}

	return mockRepo, cleanup
}
func TestUserServiceSuite(t *testing.T) {
	t.Run("TestGetAllUsers", TestUserService_GetAllUsers)
	t.Run("TestGetUserByID", TestUserService_GetUserByID)
	t.Run("TestCreateUser", TestUserService_CreateUser)
	t.Run("TestUpdateUser", TestUserService_UpdateUser)
	t.Run("TestDeleteUser", TestUserService_DeleteUser)
	t.Run("TestCheckRole", TestUserService_CheckRole)
}
func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo, cleanup := setupMockUserRepository(t)
	defer cleanup()

	userService := services.NewUserService(mockRepo)

	expectedUsers := []*models.Users{
		{
			ID:       1,
			UserName: "John Doe",
			Role:     "admin",
		},
		{
			ID:       2,
			UserName: "Jane Smith",
			Role:     "user",
		},
	}

	mockUserRepo := mockRepo.(*mocks.UserRepository)
	mockUserRepo.On("GetAllUsers", mock.Anything).Return(expectedUsers, nil)

	// Test success case
	actualUsers, err := userService.GetAllUsers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, actualUsers)

	// Create a new mock for the error case to avoid expectation conflicts
	mockRepoError, cleanupError := setupMockUserRepository(t)
	defer cleanupError()

	userServiceError := services.NewUserService(mockRepoError)
	mockUserRepoError := mockRepoError.(*mocks.UserRepository)
	mockUserRepoError.On("GetAllUsers", mock.Anything).Return(nil, errors.New("database error"))

	// Test error case
	actualUsersError, errError := userServiceError.GetAllUsers(context.Background())
	assert.Error(t, errError)
	assert.Nil(t, actualUsersError)
}
func TestUserService_GetUserByID(t *testing.T) {
	mockRepo, cleanup := setupMockUserRepository(t)
	defer cleanup()

	userService := services.NewUserService(mockRepo)

	expectedUser := &models.Users{
		ID:       1,
		UserName: "John Doe",
		Role:     "admin",
	}

	// The mockRepo is actually a *mocks.UserRepository, not repository.UserRepository
	mockUserRepo := mockRepo.(*mocks.UserRepository)
	mockUserRepo.On("GetUserByID", mock.Anything, uint(1)).Return(expectedUser, nil)

	// Test success case
	actualUser, err := userService.GetUserByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)

	// Test error case
	mockUserRepo.On("GetUserByID", mock.Anything, uint(2)).Return(nil, errors.New("user not found"))
	actualUser, err = userService.GetUserByID(context.Background(), 2)
	assert.Error(t, err)
	assert.Nil(t, actualUser)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo, cleanup := setupMockUserRepository(t)
	defer cleanup()

	userService := services.NewUserService(mockRepo)

	expectedUser := &models.Users{
		ID:             1,
		UserName:       "John Doe",
		HashedPassword: "hashed_password_123",
		Role:           "admin",
	}

	mockUserRepo := mockRepo.(*mocks.UserRepository)

	// Test success case
	mockUserRepo.On("CreateUser", mock.Anything, expectedUser).Return(nil)
	err := userService.CreateUser(context.Background(), expectedUser)
	assert.NoError(t, err)

	// Test error case
	mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.Users")).Return(errors.New("database error"))
	err = userService.CreateUser(context.Background(), &models.Users{UserName: "Invalid User"})
	assert.Error(t, err)
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo, cleanup := setupMockUserRepository(t)
	defer cleanup()

	userService := services.NewUserService(mockRepo)

	updates := map[string]interface{}{
		"user_name":       "John Doe",
		"hashed_password": "new_hashed_password",
		"role":            "admin",
	}

	mockUserRepo := mockRepo.(*mocks.UserRepository)

	// Test success case
	mockUserRepo.On("UpdateUser", mock.Anything, uint(1), updates).Return(nil)
	err := userService.UpdateUser(context.Background(), 1, updates)
	assert.NoError(t, err)

	// Test error case
	mockUserRepo.On("UpdateUser", mock.Anything, uint(999), mock.Anything).Return(errors.New("user not found"))
	err = userService.UpdateUser(context.Background(), 999, updates)
	assert.Error(t, err)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo, cleanup := setupMockUserRepository(t)
	defer cleanup()

	userService := services.NewUserService(mockRepo)

	mockUserRepo := mockRepo.(*mocks.UserRepository)

	// Test success case
	mockUserRepo.On("DeleteUser", mock.Anything, uint(1)).Return(nil)
	err := userService.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)

	// Test error case
	mockUserRepo.On("DeleteUser", mock.Anything, uint(999)).Return(errors.New("user not found"))
	err = userService.DeleteUser(context.Background(), 999)
	assert.Error(t, err)
}

func TestUserService_CheckRole(t *testing.T) {
	mockRepo, cleanup := setupMockUserRepository(t)
	defer cleanup()

	userService := services.NewUserService(mockRepo)
	mockUserRepo := mockRepo.(*mocks.UserRepository)

	// Test success case
	mockUserRepo.On("CheckRole", mock.Anything, uint(1)).Return("admin", nil)
	role, err := userService.CheckRole(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "admin", role)

	// Test error case
	mockUserRepo.On("CheckRole", mock.Anything, uint(999)).Return("", errors.New("user not found"))
	role, err = userService.CheckRole(context.Background(), 999)
	assert.Error(t, err)
	assert.Equal(t, "", role)
}
