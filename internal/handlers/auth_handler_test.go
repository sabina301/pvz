package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
	"net/http/httptest"
	"pvz/internal/logger"
	"pvz/internal/middleware"
	"pvz/pkg/errors"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pvz/internal/bootstrap"
	"pvz/internal/handlers"
	"pvz/internal/mocks"
	"pvz/internal/models/auth"
)

func TestAuthHandlers(t *testing.T) {
	e := echo.New()
	apiValidator := handlers.NewApiValidator()
	e.Validator = apiValidator
	logger.Init("debug")
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		setupMock      func(*mocks.AuthService)
		expectedStatus int
		wantResponse   interface{}
	}{
		{
			name:   "success login",
			method: http.MethodPost,
			path:   "/login",
			body: map[string]interface{}{
				"email":    "user@avito.ru",
				"password": "avito12345",
			},
			setupMock: func(m *mocks.AuthService) {
				m.On("Login", mock.MatchedBy(func(req auth.LoginRequest) bool {
					return req.Email == "user@avito.ru" &&
						req.Password == "avito12345"
				})).Return("token-avito", nil)
			},
			expectedStatus: http.StatusOK,
			wantResponse:   "token-avito",
		},
		{
			name:   "login user not found",
			method: http.MethodPost,
			path:   "/login",
			body: map[string]interface{}{
				"email":    "nonexistent@avito.ru",
				"password": "avito12345",
			},
			setupMock: func(m *mocks.AuthService) {
				m.On("Login", mock.Anything).Return("", errors.NewInvalidCredentials())
			},
			expectedStatus: http.StatusUnauthorized,
			wantResponse:   map[string]string{"message": "invalid credentials"},
		},
		{
			name:   "login invalid request",
			method: http.MethodPost,
			path:   "/login",
			body: map[string]interface{}{
				"email": "invalid-email",
			},
			setupMock:      func(m *mocks.AuthService) {},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "property Email has bad value format invalid-email"},
		},
		{
			name:   "success register",
			method: http.MethodPost,
			path:   "/register",
			body: map[string]interface{}{
				"email":    "new@avito.ru",
				"password": "avito12345",
				"role":     "employee",
			},
			setupMock: func(m *mocks.AuthService) {
				m.On("Register", mock.Anything).Return(auth.RegisterResponse{
					Id:    "123",
					Email: "new@avito.ru",
					Role:  "employee",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			wantResponse: auth.RegisterResponse{
				Id:    "123",
				Email: "new@avito.ru",
				Role:  "employee",
			},
		},
		{
			name:   "register existing user",
			method: http.MethodPost,
			path:   "/register",
			body: map[string]interface{}{
				"email":    "exists@avito.ru",
				"password": "avito12345",
				"role":     "employee",
			},
			setupMock: func(m *mocks.AuthService) {
				m.On("Register", mock.Anything).Return(auth.RegisterResponse{}, errors.NewObjectAlreadyExists("user", "email", "exists@avito.ru"))
			},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "user with email exists@avito.ru already exists"},
		},
		{
			name:   "register invalid role",
			method: http.MethodPost,
			path:   "/register",
			body: map[string]interface{}{
				"email":    "new@avito.ru",
				"password": "avito12345",
				"role":     "invalid",
			},
			setupMock:      func(m *mocks.AuthService) {},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "property Role has wrong value invalid"},
		},
		{
			name:   "success dummy login moderator",
			method: http.MethodPost,
			path:   "/dummy-login",
			body: map[string]interface{}{
				"role": "moderator",
			},
			setupMock: func(m *mocks.AuthService) {
				m.On("DummyLogin", mock.Anything).Return("dummy-token-moderator", nil)
			},
			expectedStatus: http.StatusOK,
			wantResponse:   "dummy-token-moderator",
		},
		{
			name:   "dummy login invalid role",
			method: http.MethodPost,
			path:   "/dummy-login",
			body: map[string]interface{}{
				"role": "invalid",
			},
			setupMock:      func(m *mocks.AuthService) {},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "property Role has wrong value invalid"},
		},
		{
			name:   "login internal error",
			method: http.MethodPost,
			path:   "/login",
			body: map[string]interface{}{
				"email":    "error@avito.ru",
				"password": "avito12345",
			},
			setupMock: func(m *mocks.AuthService) {
				m.On("Login", mock.Anything).Return("", errors.NewInternalError())
			},
			expectedStatus: http.StatusInternalServerError,
			wantResponse:   map[string]string{"message": "internal error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := mocks.NewAuthService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockAuth)
			}

			h := handlers.NewAuthHandler(bootstrap.Deps{
				AuthService: mockAuth,
			})

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(tt.path)

			err := middleware.HandleError(func(c echo.Context) error {
				switch tt.path {
				case "/login":
					return h.Login(c)
				case "/register":
					return h.Register(c)
				case "/dummy-login":
					return h.DummyLogin(c)
				default:
					t.Fatalf("unsupported path %s", tt.path)
					return nil
				}
			})(c)

			if err != nil {
				e.HTTPErrorHandler(err, c)
			}
			log.Println(err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.wantResponse != nil {
				switch want := tt.wantResponse.(type) {
				case string:
					// Для строковых ответов удаляем кавычки
					assert.JSONEq(t, `"`+want+`"`, rec.Body.String())
				case map[string]string:
					var response map[string]string
					err := json.Unmarshal(rec.Body.Bytes(), &response)
					require.NoError(t, err)
					assert.Equal(t, want, response)
				case auth.RegisterResponse:
					var response auth.RegisterResponse
					err := json.Unmarshal(rec.Body.Bytes(), &response)
					require.NoError(t, err)
					assert.Equal(t, want, response)
				default:
					t.Fatal("unsupported response type")
				}
			}

			mockAuth.AssertExpectations(t)
		})
	}
}
