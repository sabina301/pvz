package tokens_test

import (
	"os"
	"pvz/internal/tokens"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvz/configs"
	"pvz/internal/models/auth"
)

// Мокаем конфигурацию
type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) GetJwtSecret() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfig) GetJwtExpiration() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func LoadTestEnv() {
	os.Setenv("DB_HOST", "test_host")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("JWT_SECRET", "test_secret")
}

func TestGenerateJwt(t *testing.T) {
	LoadTestEnv()
	_, _ = configs.LoadConfig()

	userId := "12345"
	role := "moderator"

	token, err := tokens.GenerateJwt(userId, auth.Role(role))

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedClaims, parseErr := tokens.ParseJwt(token)
	assert.NoError(t, parseErr)
	assert.Equal(t, userId, parsedClaims.UserId)
	assert.Equal(t, role, string(parsedClaims.Role))
}

func TestGenerateDummyJwt(t *testing.T) {
	LoadTestEnv()
	_, _ = configs.LoadConfig()

	role := "employee"

	token, err := tokens.GenerateDummyJwt(auth.Role(role))

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedClaims, parseErr := tokens.ParseJwt(token)
	assert.NoError(t, parseErr)
	assert.Equal(t, role, string(parsedClaims.Role))
}

func TestParseJwt(t *testing.T) {
	LoadTestEnv()
	_, _ = configs.LoadConfig()

	role := "employee"
	userId := "12345"
	token, err := tokens.GenerateJwt(userId, auth.Role(role))
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedClaims, parseErr := tokens.ParseJwt(token)
	assert.NoError(t, parseErr)
	assert.Equal(t, userId, parsedClaims.UserId)
	assert.Equal(t, role, string(parsedClaims.Role))

	invalidToken := "invalidtoken"
	parsedClaims, parseErr = tokens.ParseJwt(invalidToken)
	assert.Error(t, parseErr)
	assert.Nil(t, parsedClaims)
}

func TestGenerateJwtWithMockConfig(t *testing.T) {
	mockConfig := new(MockConfig)
	mockConfig.On("GetJwtSecret").Return("mocksecret")
	mockConfig.On("GetJwtExpiration").Return(time.Minute)

	configs.AppConfiguration.Auth.JwtSecret = mockConfig.GetJwtSecret()
	configs.AppConfiguration.Auth.Expiration = mockConfig.GetJwtExpiration()

	role := "moderator"
	userId := "user123"

	token, err := tokens.GenerateJwt(userId, auth.Role(role))

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedClaims, parseErr := tokens.ParseJwt(token)
	assert.NoError(t, parseErr)
	assert.Equal(t, userId, parsedClaims.UserId)
	assert.Equal(t, role, string(parsedClaims.Role))
}
