package services_test

import (
	"pvz/internal/models/pvz"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pvz/internal/logger"
	"pvz/internal/models"
	"pvz/internal/models/product"
	"pvz/internal/models/reception"
	"pvz/internal/repositories"
	"pvz/internal/services"
	"pvz/pkg/errors"
)

type MockProductRepo struct{ mock.Mock }
type MockPvzRepo struct{ mock.Mock }

type MockReceptionRepo struct{ mock.Mock }

func (m *MockReceptionRepo) Create(q repositories.Querier, req models.Reception) (*models.Reception, error) {
	args := m.Called(q, req)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (m *MockPvzRepo) GetById(q repositories.Querier, id string) (*models.Pvz, error) {
	args := m.Called(q, id)
	return args.Get(0).(*models.Pvz), args.Error(1)
}

func (m *MockPvzRepo) Create(q repositories.Querier, req pvz.CreateRequest) (*models.Pvz, error) {
	args := m.Called(q, req)
	return args.Get(0).(*models.Pvz), args.Error(1)
}

func (m *MockPvzRepo) ListWithFilterDate(q repositories.Querier, req pvz.ListRequest, offset int) ([]pvz.RawList, error) {
	args := m.Called(q, req, offset)
	return args.Get(0).([]pvz.RawList), args.Error(1)
}

func (m *MockReceptionRepo) SetStatus(q repositories.Querier, status string, pvzId string) (*models.Reception, error) {
	args := m.Called(q, status, pvzId)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (m *MockProductRepo) AddInReception(q repositories.Querier, p models.Product) (*models.Product, error) {
	args := m.Called(q, p)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepo) DeleteLast(q repositories.Querier, recId string) (*models.Product, error) {
	args := m.Called(q, recId)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockReceptionRepo) GetByPvzId(q repositories.Querier, pvzId string) (*models.Reception, error) {
	args := m.Called(q, pvzId)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func TestProductService_AddInReception(t *testing.T) {
	logger.Init("debug")

	t.Run("successful product addition", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		productRepo := new(MockProductRepo)
		pvzRepo := new(MockPvzRepo)
		recRepo := new(MockReceptionRepo)
		service := services.NewProductService(productRepo, pvzRepo, recRepo, db)

		req := product.AddInReceptionRequest{
			PvzId: "550e8400-e29b-41d4-a716-446655440000",
			Type:  "электроника",
		}

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return(&models.Pvz{Id: req.PvzId}, nil).
			Once()

		recRepo.On("GetByPvzId", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return(&models.Reception{
				Id:     "rec123",
				Status: reception.InProgressStatus,
			}, nil).
			Once()

		expectedProduct := &models.Product{
			Id:          "prod456",
			DateTime:    time.Now(),
			Type:        req.Type,
			ReceptionId: "rec123",
		}

		productRepo.On("AddInReception", mock.AnythingOfType("*sql.Tx"), mock.MatchedBy(func(p models.Product) bool {
			return p.ReceptionId == "rec123" && p.Type == req.Type
		})).Return(expectedProduct, nil).
			Once()

		mockDB.ExpectCommit()

		resp, err := service.AddInReception(req)

		require.NoError(t, err)
		require.Equal(t, expectedProduct.Id, resp.Id)
		require.Equal(t, expectedProduct.ReceptionId, resp.ReceptionId)
		productRepo.AssertExpectations(t)
		pvzRepo.AssertExpectations(t)
		recRepo.AssertExpectations(t)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("pvz not found", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		service := services.NewProductService(nil, pvzRepo, nil, db)

		req := product.AddInReceptionRequest{
			PvzId: "invalid-pvz",
			Type:  "одежда",
		}

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return((*models.Pvz)(nil), nil).
			Once()
		mockDB.ExpectRollback()

		_, err := service.AddInReception(req)

		require.Error(t, err)
		require.IsType(t, errors.ObjectNotFound{}, err)
		pvzRepo.AssertExpectations(t)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("database error during pvz check", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		service := services.NewProductService(nil, pvzRepo, nil, db)

		req := product.AddInReceptionRequest{
			PvzId: "550e8400-e29b-41d4-a716-446655440000",
			Type:  "электроника",
		}

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return((*models.Pvz)(nil), errors.NewInternalError()).
			Once()
		mockDB.ExpectRollback()

		_, err := service.AddInReception(req)

		require.Error(t, err)
		require.IsType(t, errors.InternalError{}, err)
		pvzRepo.AssertExpectations(t)
	})

	t.Run("transaction commit failure", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		recRepo := new(MockReceptionRepo)
		productRepo := new(MockProductRepo)
		service := services.NewProductService(productRepo, pvzRepo, recRepo, db)

		req := product.AddInReceptionRequest{
			PvzId: "550e8400-e29b-41d4-a716-446655440000",
			Type:  "электроника",
		}

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return(&models.Pvz{Id: req.PvzId}, nil)
		recRepo.On("GetByPvzId", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return(&models.Reception{Id: "rec123", Status: reception.InProgressStatus}, nil)
		productRepo.On("AddInReception", mock.Anything, mock.Anything).
			Return(&models.Product{Id: "prod123"}, nil)
		mockDB.ExpectCommit().WillReturnError(errors.NewInternalError())

		_, err := service.AddInReception(req)

		require.Error(t, err)
		require.IsType(t, errors.InternalError{}, err)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("reception not found", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		recRepo := new(MockReceptionRepo)
		service := services.NewProductService(nil, pvzRepo, recRepo, db)

		req := product.AddInReceptionRequest{
			PvzId: "550e8400-e29b-41d4-a716-446655440000",
			Type:  "одежда",
		}

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return(&models.Pvz{Id: req.PvzId}, nil)
		recRepo.On("GetByPvzId", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return((*models.Reception)(nil), nil)
		mockDB.ExpectRollback()

		_, err := service.AddInReception(req)

		require.Error(t, err)
		require.IsType(t, errors.ObjectNotFound{}, err)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("product insertion failure", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		recRepo := new(MockReceptionRepo)
		productRepo := new(MockProductRepo)
		service := services.NewProductService(productRepo, pvzRepo, recRepo, db)

		req := product.AddInReceptionRequest{
			PvzId: "550e8400-e29b-41d4-a716-446655440000",
			Type:  "обувь",
		}

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return(&models.Pvz{Id: req.PvzId}, nil)
		recRepo.On("GetByPvzId", mock.AnythingOfType("*sql.Tx"), req.PvzId).
			Return(&models.Reception{Id: "rec123", Status: reception.InProgressStatus}, nil)
		productRepo.On("AddInReception", mock.Anything, mock.Anything).
			Return((*models.Product)(nil), errors.NewInternalError())
		mockDB.ExpectRollback()

		_, err := service.AddInReception(req)

		require.Error(t, err)
		require.IsType(t, errors.InternalError{}, err)
		productRepo.AssertExpectations(t)
	})

	t.Run("transaction begin failure", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		service := services.NewProductService(nil, nil, nil, db)
		req := product.AddInReceptionRequest{}

		mockDB.ExpectBegin().WillReturnError(errors.NewInternalError())

		_, err := service.AddInReception(req)

		require.Error(t, err)
		require.IsType(t, errors.InternalError{}, err)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})
}
