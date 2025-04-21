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
	"pvz/internal/models/pvz"
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

func TestPvzHandlers(t *testing.T) {
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
		params         map[string]string
		queryParams    map[string]string
		body           interface{}
		setupMock      func(*mocks.PvzService)
		expectedStatus int
		wantResponse   interface{}
	}{
		{
			name:   "success create with generated id",
			method: http.MethodPost,
			path:   "/pvz",
			body: map[string]interface{}{
				"city": "Москва",
			},
			setupMock: func(m *mocks.PvzService) {
				m.On("Create", mock.Anything).Return(pvz.CreateResponse{
					Id:               testUUID,
					RegistrationDate: validTime,
					City:             "Москва",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			wantResponse: pvz.CreateResponse{
				Id:               testUUID,
				RegistrationDate: validTime,
				City:             "Москва",
			},
		},
		{
			name:   "create with existing id",
			method: http.MethodPost,
			path:   "/pvz",
			body: map[string]interface{}{
				"id":   testUUID,
				"city": "Санкт-Петербург",
			},
			setupMock: func(m *mocks.PvzService) {
				m.On("Create", mock.Anything).Return(
					pvz.CreateResponse{},
					errors.NewObjectAlreadyExists("pvz", "id", testUUID),
				)
			},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "pvz with id " + testUUID + " already exists"},
		},
		{
			name:   "create invalid city",
			method: http.MethodPost,
			path:   "/pvz",
			body: map[string]interface{}{
				"city": "Новосибирск",
			},
			setupMock:      func(m *mocks.PvzService) {},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "property City has wrong value Новосибирск"},
		},
		{
			name:   "success delete last product",
			method: http.MethodDelete,
			path:   "/pvz/:pvzId/products/last",
			params: map[string]string{"pvzId": testUUID},
			setupMock: func(m *mocks.PvzService) {
				m.On("DeleteLastProduct", testUUID).Return(pvz.DeleteLastProductResponse{
					Id: "product-123",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			wantResponse:   pvz.DeleteLastProductResponse{Id: "product-123"},
		},
		{
			name:   "delete last product pvz not found",
			method: http.MethodDelete,
			path:   "/pvz/:pvzId/products/last",
			params: map[string]string{"pvzId": testUUID},
			setupMock: func(m *mocks.PvzService) {
				m.On("DeleteLastProduct", testUUID).Return(
					pvz.DeleteLastProductResponse{},
					errors.NewObjectNotFound("pvz"),
				)
			},
			expectedStatus: http.StatusNotFound,
			wantResponse:   map[string]string{"message": "pvz not found"},
		},
		{
			name:   "delete last product invalid uuid",
			method: http.MethodDelete,
			path:   "/pvz/:pvzId/products/last",
			params: map[string]string{"pvzId": "invalid"},
			setupMock: func(m *mocks.PvzService) {
			},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "bad param value invalid"},
		},
		{
			name:   "success close last reception",
			method: http.MethodPost,
			path:   "/pvz/:pvzId/receptions/close",
			params: map[string]string{"pvzId": testUUID},
			setupMock: func(m *mocks.PvzService) {
				m.On("CLoseLastReception", testUUID).Return(pvz.CloseLastProductResponse{
					Id:       "rec-123",
					DateTime: validTime,
					PvzId:    testUUID,
					Status:   "closed",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			wantResponse: pvz.CloseLastProductResponse{
				Id:       "rec-123",
				DateTime: validTime,
				PvzId:    testUUID,
				Status:   "closed",
			},
		},
		{
			name:   "close reception no in-progress",
			method: http.MethodPost,
			path:   "/pvz/:pvzId/receptions/close",
			params: map[string]string{"pvzId": testUUID},
			setupMock: func(m *mocks.PvzService) {
				m.On("CLoseLastReception", testUUID).Return(
					pvz.CloseLastProductResponse{},
					errors.NewNoInProgressReception(),
				)
			},
			expectedStatus: http.StatusNotFound,
			wantResponse:   map[string]string{"message": "no in-progress reception"},
		},
		{
			name:   "success list with filter",
			method: http.MethodGet,
			path:   "/pvz",
			queryParams: map[string]string{
				"startDate": "2024-01-01T00:00:00Z",
				"endDate":   "2024-01-31T23:59:59Z",
				"page":      "2",
				"limit":     "20",
			},
			setupMock: func(m *mocks.PvzService) {
				start, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
				end, _ := time.Parse(time.RFC3339, "2024-01-31T23:59:59Z")
				m.On("ListWithFilterDate", pvz.ListRequest{
					StartDate: &start,
					EndDate:   &end,
					Page:      2,
					Limit:     20,
				}).Return([]pvz.ListResponse{
					{
						Pvz: pvz.Pvz{
							Id:               testUUID,
							RegistrationDate: validTime,
							City:             "Москва",
						},
						Receptions: []pvz.ReceptionProducts{
							{
								Reception: pvz.Reception{
									Id:       "rec-123",
									DateTime: validTime,
									PvzId:    testUUID,
									Status:   "closed",
								},
								Products: []pvz.Product{
									{
										Id:          "product-123",
										DateTime:    validTime,
										Type:        "электроника",
										ReceptionId: "rec-123",
									},
								},
							},
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			wantResponse: []pvz.ListResponse{
				{
					Pvz: pvz.Pvz{
						Id:               testUUID,
						RegistrationDate: validTime,
						City:             "Москва",
					},
					Receptions: []pvz.ReceptionProducts{
						{
							Reception: pvz.Reception{
								Id:       "rec-123",
								DateTime: validTime,
								PvzId:    testUUID,
								Status:   "closed",
							},
							Products: []pvz.Product{
								{
									Id:          "product-123",
									DateTime:    validTime,
									Type:        "электроника",
									ReceptionId: "rec-123",
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "list invalid date format",
			method: http.MethodGet,
			path:   "/pvz",
			queryParams: map[string]string{
				"startDate": "invalid",
				"endDate":   "2024-01-31",
			},
			setupMock:      func(m *mocks.PvzService) {},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "bad param value invalid"},
		},
		{
			name:   "list start after end date",
			method: http.MethodGet,
			path:   "/pvz",
			queryParams: map[string]string{
				"startDate": "2024-02-01T00:00:00Z",
				"endDate":   "2024-01-01T00:00:00Z",
			},
			setupMock: func(m *mocks.PvzService) {
				start, _ := time.Parse(time.RFC3339, "2024-02-01T00:00:00Z")
				end, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
				m.On("ListWithFilterDate", pvz.ListRequest{
					StartDate: &start,
					EndDate:   &end,
					Page:      1,
					Limit:     10,
				}).Return(nil, errors.NewStartDateAfterEndDate())
			},
			expectedStatus: http.StatusBadRequest,
			wantResponse:   map[string]string{"message": "start date after end date"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPvz := mocks.NewPvzService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockPvz)
			}

			h := handlers.NewPvzHandler(bootstrap.Deps{
				PvzService: mockPvz,
			})

			var reqBody []byte
			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
			}
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			for k, v := range tt.params {
				c.SetParamNames(k)
				c.SetParamValues(v)
			}

			err := middleware.HandleError(func(c echo.Context) error {
				switch tt.method + ":" + tt.path {
				case http.MethodPost + ":/pvz":
					return h.Create(c)
				case http.MethodDelete + ":/pvz/:pvzId/products/last":
					return h.DeleteLastProduct(c)
				case http.MethodPost + ":/pvz/:pvzId/receptions/close":
					return h.CloseLastReception(c)
				case http.MethodGet + ":/pvz":
					return h.ListWithFilterDate(c)
				default:
					t.Fatalf("unknown route: %s %s", tt.method, tt.path)
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
				case pvz.CreateResponse:
					var response pvz.CreateResponse
					require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
					assert.Equal(t, want, response)
				case pvz.DeleteLastProductResponse:
					var response pvz.DeleteLastProductResponse
					require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
					assert.Equal(t, want, response)
				case pvz.CloseLastProductResponse:
					var response pvz.CloseLastProductResponse
					require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
					assert.Equal(t, want, response)
				case []pvz.ListResponse:
					var response []pvz.ListResponse
					require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
					assert.Equal(t, want, response)
				case map[string]string:
					var response map[string]string
					require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
					assert.Equal(t, want, response)
				default:
					t.Fatal("unsupported response type")
				}
			}

			mockPvz.AssertExpectations(t)
		})
	}
}
