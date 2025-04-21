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
	"pvz/internal/models/reception"
	"pvz/pkg/errors"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pvz/internal/bootstrap"
	"pvz/internal/handlers"
	"pvz/internal/mocks"
)

func TestReceptionHandlers(t *testing.T) {
	e := echo.New()
	apiValidator := handlers.NewApiValidator()
	e.Validator = apiValidator
	logger.Init("debug")

	validTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		setupMock      func(*mocks.ReceptionService)
		expectedStatus int
		wantResponse   interface{}
	}{
		{
			name:   "success create reception",
			method: http.MethodPost,
			path:   "/receptions",
			body: map[string]interface{}{
				"pvzId": testUUID,
			},
			setupMock: func(m *mocks.ReceptionService) {
				m.On("Create", mock.MatchedBy(func(req reception.CreateRequest) bool {
					return req.PvzId == testUUID
				})).Return(reception.CreateResponse{
					Id:       "rec-123",
					DateTime: validTime,
					PvzId:    testUUID,
					Status:   reception.InProgressStatus,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			wantResponse: reception.CreateResponse{
				Id:       "rec-123",
				DateTime: validTime,
				PvzId:    testUUID,
				Status:   reception.InProgressStatus,
			},
		},
		{
			name:   "invalid pvzId format",
			method: http.MethodPost,
			path:   "/receptions",
			body: map[string]interface{}{
				"pvzId": "invalid-uuid",
			},
			setupMock:      func(m *mocks.ReceptionService) {},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "property PvzId has bad value format invalid-uuid"},
		},
		{
			name:           "missing pvzId",
			method:         http.MethodPost,
			path:           "/receptions",
			body:           map[string]interface{}{},
			setupMock:      func(m *mocks.ReceptionService) {},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "property PvzId is missing"},
		},
		{
			name:   "reception already exists",
			method: http.MethodPost,
			path:   "/receptions",
			body: map[string]interface{}{
				"pvzId": testUUID,
			},
			setupMock: func(m *mocks.ReceptionService) {
				m.On("Create", mock.Anything).Return(
					reception.CreateResponse{},
					errors.NewReceptionIsNotClosed(testUUID),
				)
			},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "reception in pvz with id " + testUUID + " is not closed"},
		},
		{
			name:   "pvz not found",
			method: http.MethodPost,
			path:   "/receptions",
			body: map[string]interface{}{
				"pvzId": testUUID,
			},
			setupMock: func(m *mocks.ReceptionService) {
				m.On("Create", mock.Anything).Return(
					reception.CreateResponse{},
					errors.NewObjectNotFound("pvz"),
				)
			},
			expectedStatus: http.StatusNotFound,
			wantResponse:   map[string]string{"message": "pvz not found"},
		},
		{
			name:   "internal error",
			method: http.MethodPost,
			path:   "/receptions",
			body: map[string]interface{}{
				"pvzId": testUUID,
			},
			setupMock: func(m *mocks.ReceptionService) {
				m.On("Create", mock.Anything).Return(
					reception.CreateResponse{},
					errors.NewInternalError(),
				)
			},
			expectedStatus: http.StatusInternalServerError,
			wantResponse:   map[string]string{"message": "internal error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReception := mocks.NewReceptionService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockReception)
			}

			h := handlers.NewReceptionHandler(bootstrap.Deps{
				ReceptionService: mockReception,
			})

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(tt.path)

			err := middleware.HandleError(func(c echo.Context) error {
				return h.Create(c)
			})(c)

			if err != nil {
				e.HTTPErrorHandler(err, c)
			}
			log.Println(err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.wantResponse != nil {
				switch want := tt.wantResponse.(type) {
				case reception.CreateResponse:
					var response reception.CreateResponse
					err := json.Unmarshal(rec.Body.Bytes(), &response)
					require.NoError(t, err)
					assert.Equal(t, want, response)
				case map[string]string:
					var response map[string]string
					err := json.Unmarshal(rec.Body.Bytes(), &response)
					require.NoError(t, err)
					assert.Equal(t, want, response)
				default:
					t.Fatal("unsupported response type")
				}
			}

			mockReception.AssertExpectations(t)
		})
	}
}
