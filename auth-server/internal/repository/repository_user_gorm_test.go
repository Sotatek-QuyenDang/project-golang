package repository_test

import (
	"auth-server/internal/models"
	"auth-server/internal/repository"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	cleanup := func() {
		sqlDB.Close()
	}
	return db, mock, cleanup
}
func TestUserRepositorySuite(t *testing.T) {
	t.Run("TestCreateUser", TestCreateUser)
	t.Run("TestUpdateUser", TestUpdateUser)
	t.Run("TestGetAllUsers", TestGetAllUsers)
	t.Run("TestDeleteUser", TestDeleteUser)
	t.Run("TestCheckRole", TestCheckRole)
}
func TestCreateUser(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepositoryGorm(db)

	// Create a test user with all required fields
	user := models.Users{
		UserName:       "John Doe",
		HashedPassword: "hashed_password_123",
		Role:           "admin",
	}

	// Add context
	ctx := context.Background()

	// Debug SQL query expectation
	t.Logf("Setting up SQL mock expectations")
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" \("user_name","hashed_password","role","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
		WithArgs(user.UserName, user.HashedPassword, user.Role, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	t.Logf("Testing CreateUser with user: %+v", user)
	err := repo.CreateUser(ctx, &user)
	if err != nil {
		t.Logf("CreateUser error: %v", err)
		t.Logf("SQL error details: %v", err)
	} else {
		t.Logf("CreateUser successful, user ID: %d", user.ID)
	}
	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)

	// Verify all expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Logf("Unfulfilled expectations: %v", err)
		// Print actual queries that were executed
		t.Logf("Actual queries might differ from expected ones")
	}
	assert.NoError(t, err, "There were unfulfilled expectations")
}

func TestUpdateUser(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepositoryGorm(db)

	userID := uint(1)
	updates := map[string]interface{}{
		"user_name":       "John Doe",
		"hashed_password": "hashed_password_123",
		"role":            "admin",
	}

	ctx := context.Background()

	t.Logf("Testing UpdateUser with user ID: %d", userID)
	t.Logf("Setting up SQL mock expectations")

	// Bắt đầu transaction
	mock.ExpectBegin()

	// Khớp đúng SQL thực tế GORM sinh ra
	mock.ExpectExec(`UPDATE "users" SET (.+) WHERE id = \$5`).
		WithArgs(
			updates["hashed_password"],
			updates["role"],
			updates["user_name"],
			sqlmock.AnyArg(), // updated_at
			userID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Kết thúc transaction
	mock.ExpectCommit()

	t.Logf("Updating user with values: %+v", updates)

	err := repo.UpdateUser(ctx, userID, updates)
	if err != nil {
		t.Logf("UpdateUser error: %v", err)
		t.Logf("SQL error details: %v", err)
	} else {
		t.Logf("UpdateUser successful for user ID: %d", userID)
	}
	assert.NoError(t, err)

	// Kiểm tra tất cả expectations đã được thực hiện
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Logf("Unfulfilled expectations: %v", err)
		t.Logf("Actual queries might differ from expected ones")
	}
	assert.NoError(t, err, "There were unfulfilled expectations")
}

func TestGetAllUsers(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepositoryGorm(db)
	ctx := context.Background()

	// Setup mock expectations
	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "role"}).
			AddRow(1, "John Doe", "admin").
			AddRow(2, "Jane Smith", "user"))

	users, err := repo.GetAllUsers(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, users, "Users should not be nil")
	assert.Equal(t, 2, len(users), "Should return 2 users")
	assert.Equal(t, uint(1), users[0].ID)
	assert.Equal(t, "John Doe", users[0].UserName)
	assert.Equal(t, "admin", users[0].Role)
	assert.Equal(t, uint(2), users[1].ID)
	assert.Equal(t, "Jane Smith", users[1].UserName)
	assert.Equal(t, "user", users[1].Role)

	// Verify all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "There were unfulfilled expectations")
}

func TestDeleteUser(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepositoryGorm(db)

	userID := int64(1)

	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "users" WHERE "users"\."id" = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	t.Logf("Testing DeleteUser with user ID: %d", userID)
	err := repo.DeleteUser(ctx, uint(userID))
	if err != nil {
		t.Logf("DeleteUser error: %v", err)
		t.Logf("SQL error details: %v", err)
	} else {
		t.Logf("DeleteUser successful for user ID: %d", userID)
	}
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Logf("Unfulfilled expectations: %v", err)
		// Print actual queries that were executed
		t.Logf("Actual queries might differ from expected ones")
	}
	assert.NoError(t, err, "There were unfulfilled expectations")
}

func TestCheckRole(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepositoryGorm(db)

	user := models.Users{
		ID:             1,
		UserName:       "John Doe",
		HashedPassword: "hashed_password_123",
		Role:           "admin",
	}

	ctx := context.Background()

	t.Logf("Testing CheckRole with user ID: %d", user.ID)

	mock.ExpectQuery(`SELECT "role" FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))
	role, err := repo.CheckRole(ctx, uint(user.ID))
	if err != nil {
		t.Logf("CheckRole error: %v", err)
	} else {
		t.Logf("CheckRole successful. Expected role: %s, Got role: %s", user.Role, role)
	}

	assert.NoError(t, err)
	assert.Equal(t, user.Role, role)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Logf("Unfulfilled expectations: %v", err)
		t.Logf("Actual queries might differ from expected ones")
	}
	assert.NoError(t, err, "There were unfulfilled expectations")
}
