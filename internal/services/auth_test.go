package services_test

import (
	"database/sql"
	"os"
	"pvz/configs"
	"pvz/internal/logger"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"pvz/internal/models"
	"pvz/internal/models/auth"
	"pvz/internal/repositories"
	"pvz/internal/services"
	"pvz/pkg/errors"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByEmail(q repositories.Querier, email string) (*models.User, error) {
	args := m.Called(q, email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(q repositories.Querier, user models.User) (string, error) {
	args := m.Called(q, user)
	return args.String(0), args.Error(1)
}

func LoadTestEnv() {
	os.Setenv("DB_HOST", "test_host")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("JWT_SECRET", "test_secret")
}

func TestAuthService_Login(t *testing.T) {
	logger.Init("debug")
	LoadTestEnv()
	_, _ = configs.LoadConfig()
	mockRepo := new(MockUserRepository)
	db, _ := sql.Open("postgres", "") // dummy connection
	service := services.NewAuthService(mockRepo, db)

	t.Run("success login", func(t *testing.T) {
		email := "user@avito.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		mockRepo.On("GetByEmail", db, email).Return(&models.User{
			Id:           "123",
			Email:        email,
			PasswordHash: string(hashedPassword),
			Role:         "employee",
		}, nil)

		token, err := service.Login(auth.LoginRequest{
			Email:    email,
			Password: password,
		})

		require.NoError(t, err)
		require.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		email := "user@avito.com"
		mockRepo.On("GetByEmail", db, email).Return(&models.User{
			PasswordHash: "invalid_hash",
		}, nil)

		_, err := service.Login(auth.LoginRequest{
			Email:    email,
			Password: "wrong_password",
		})

		require.ErrorIs(t, err, errors.NewInvalidCredentials())
	})

	t.Run("user not found", func(t *testing.T) {
		email := "notfound@avito.com"
		mockRepo.On("GetByEmail", db, email).Return((*models.User)(nil), nil)

		_, err := service.Login(auth.LoginRequest{Email: email})

		require.ErrorIs(t, err, errors.NewInvalidCredentials())
	})

	t.Run("database error", func(t *testing.T) {
		email := "error@avito.com"
		mockRepo.On("GetByEmail", db, email).Return((*models.User)(nil), errors.NewInternalError())

		_, err := service.Login(auth.LoginRequest{Email: email})

		require.ErrorIs(t, err, errors.NewInternalError())
	})
}

func NewTestDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Failed to create mock database")
	}
	return db, mock
}

func TestAuthService_DummyLogin(t *testing.T) {
	db, _ := NewTestDB()
	defer db.Close()
	logger.Init("debug")
	LoadTestEnv()
	_, _ = configs.LoadConfig()
	mockRepo := new(MockUserRepository)
	service := services.NewAuthService(mockRepo, db)

	t.Run("success moderator login", func(t *testing.T) {
		token, err := service.DummyLogin(auth.DummyLoginRequest{Role: "moderator"})
		require.NoError(t, err)
		require.Contains(t, token, "ey")
	})

}

func TestAuthService_Register(t *testing.T) {
	logger.Init("debug")
	LoadTestEnv()
	_, _ = configs.LoadConfig()

	t.Run("successful registration", func(t *testing.T) {
		db, mockDB := NewTestDB()
		defer db.Close()

		mockRepo := new(MockUserRepository)
		service := services.NewAuthService(mockRepo, db)

		req := auth.RegisterRequest{
			Email:    "test@avito.com",
			Password: "validpassword123",
			Role:     "employee",
		}

		mockDB.ExpectBegin()

		mockRepo.On("GetByEmail", mock.AnythingOfType("*sql.Tx"), req.Email).
			Return((*models.User)(nil), nil).Once()

		mockRepo.On("Create", mock.AnythingOfType("*sql.Tx"), mock.MatchedBy(func(user models.User) bool {
			err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
			return err == nil
		})).Return("user123", nil).Once()

		mockDB.ExpectCommit()

		resp, err := service.Register(req)

		require.NoError(t, err)
		require.Equal(t, "user123", resp.Id)
		mockRepo.AssertExpectations(t)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("transaction begin error", func(t *testing.T) {
		db, mockDB := NewTestDB()
		defer db.Close()

		service := services.NewAuthService(nil, db)
		req := auth.RegisterRequest{Email: "error@example.com"}

		mockDB.ExpectBegin().WillReturnError(errors.NewInternalError())

		_, err := service.Register(req)

		require.Error(t, err)
		require.Contains(t, err.Error(), "internal error")
		require.NoError(t, mockDB.ExpectationsWereMet())
	})
}
