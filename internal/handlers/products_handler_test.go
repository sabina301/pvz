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
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pvz/internal/bootstrap"
	"pvz/internal/handlers"
	"pvz/internal/mocks"
	"pvz/internal/models/product"
)

func TestProductHandlers(t *testing.T) {
	e := echo.New()
	apiValidator := handlers.NewApiValidator()
	e.Validator = apiValidator
	logger.Init("debug")

	validTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		setupMock      func(*mocks.ProductService)
		expectedStatus int
		wantResponse   interface{}
	}{
		{
			name:   "success add product",
			method: http.MethodPost,
			path:   "/products/reception",
			body: map[string]interface{}{
				"type":  "электроника",
				"pvzId": "550e8400-e29b-41d4-a716-446655440000",
			},
			setupMock: func(m *mocks.ProductService) {
				m.On("AddInReception", mock.MatchedBy(func(req product.AddInReceptionRequest) bool {
					return req.Type == "электроника" &&
						req.PvzId == "550e8400-e29b-41d4-a716-446655440000"
				})).Return(product.AddInReceptionResponse{
					Id:          "123",
					DateTime:    validTime,
					Type:        "электроника",
					ReceptionId: "456",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			wantResponse: product.AddInReceptionResponse{
				Id:          "123",
				DateTime:    validTime,
				Type:        "электроника",
				ReceptionId: "456",
			},
		},
		{
			name:   "invalid product type",
			method: http.MethodPost,
			path:   "/products/reception",
			body: map[string]interface{}{
				"type":  "мебель",
				"pvzId": "550e8400-e29b-41d4-a716-446655440000",
			},
			setupMock:      func(m *mocks.ProductService) {},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "property Type has wrong value мебель"},
		},
		{
			name:   "missing pvzId",
			method: http.MethodPost,
			path:   "/products/reception",
			body: map[string]interface{}{
				"type": "одежда",
			},
			setupMock:      func(m *mocks.ProductService) {},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "property PvzId is missing"},
		},
		{
			name:   "pvz not found",
			method: http.MethodPost,
			path:   "/products/reception",
			body: map[string]interface{}{
				"type":  "обувь",
				"pvzId": "550e8400-e29b-41d4-a716-446655440000",
			},
			setupMock: func(m *mocks.ProductService) {
				m.On("AddInReception", mock.Anything).Return(
					product.AddInReceptionResponse{},
					errors.NewObjectNotFound("pvz"),
				)
			},
			expectedStatus: http.StatusNotFound,
			wantResponse:   map[string]string{"message": "pvz not found"},
		},
		{
			name:   "reception not in progress",
			method: http.MethodPost,
			path:   "/products/reception",
			body: map[string]interface{}{
				"type":  "электроника",
				"pvzId": "550e8400-e29b-41d4-a716-446655440000",
			},
			setupMock: func(m *mocks.ProductService) {
				m.On("AddInReception", mock.Anything).Return(
					product.AddInReceptionResponse{},
					errors.NewReceptionIsNotInProgress("456"),
				)
			},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "reception with id 456 is not closed"},
		},
		{
			name:   "internal error",
			method: http.MethodPost,
			path:   "/products/reception",
			body: map[string]interface{}{
				"type":  "электроника",
				"pvzId": "550e8400-e29b-41d4-a716-446655440000",
			},
			setupMock: func(m *mocks.ProductService) {
				m.On("AddInReception", mock.Anything).Return(
					product.AddInReceptionResponse{},
					errors.NewInternalError(),
				)
			},
			expectedStatus: http.StatusInternalServerError,
			wantResponse:   map[string]string{"message": "internal error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := mocks.NewProductService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockProduct)
			}

			h := handlers.NewProductHandler(bootstrap.Deps{
				ProductService: mockProduct,
			})

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(tt.path)

			err := middleware.HandleError(func(c echo.Context) error {
				return h.AddInReception(c)
			})(c)

			if err != nil {
				e.HTTPErrorHandler(err, c)
			}
			log.Println(err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.wantResponse != nil {
				switch want := tt.wantResponse.(type) {
				case product.AddInReceptionResponse:
					var response product.AddInReceptionResponse
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
			mockProduct.AssertExpectations(t)
		})
	}
}
