package services_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"pvz/internal/logger"
	"pvz/internal/models"
	"pvz/internal/models/pvz"
	"pvz/internal/models/reception"
	"pvz/internal/services"
	"pvz/pkg/errors"
)

func TestPvzService_Create(t *testing.T) {
	logger.Init("debug")

	t.Run("duplicate ID error", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		service := services.NewPvzService(pvzRepo, nil, nil, db)

		pvzID := "existing-id"
		req := pvz.CreateRequest{
			Id:   &pvzID,
			City: "Санкт-Петербург",
		}

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.AnythingOfType("*sql.Tx"), pvzID).
			Return(&models.Pvz{Id: pvzID}, nil).
			Once()
		mockDB.ExpectRollback()

		_, err := service.Create(req)

		require.Error(t, err)
		require.IsType(t, errors.ObjectAlreadyExists{}, err)
		pvzRepo.AssertExpectations(t)
	})

	t.Run("transaction begin error", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		service := services.NewPvzService(nil, nil, nil, db)
		mockDB.ExpectBegin().WillReturnError(errors.NewInternalError())

		_, err := service.Create(pvz.CreateRequest{City: "Казань"})

		require.Error(t, err)
		require.IsType(t, errors.InternalError{}, err)
	})
}

func TestPvzService_DeleteLastProduct(t *testing.T) {
	logger.Init("debug")

	t.Run("successful product deletion", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		productRepo := new(MockProductRepo)
		recRepo := new(MockReceptionRepo)
		service := services.NewPvzService(pvzRepo, productRepo, recRepo, db)

		pvzID := "pvz123"
		recID := "rec456"
		product := &models.Product{Id: "prod789"}

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.AnythingOfType("*sql.Tx"), pvzID).
			Return(&models.Pvz{Id: pvzID}, nil)
		recRepo.On("GetByPvzId", mock.AnythingOfType("*sql.Tx"), pvzID).
			Return(&models.Reception{Id: recID, Status: reception.InProgressStatus}, nil)
		productRepo.On("DeleteLast", mock.AnythingOfType("*sql.Tx"), recID).
			Return(product, nil)
		mockDB.ExpectCommit()

		resp, err := service.DeleteLastProduct(pvzID)

		require.NoError(t, err)
		require.Equal(t, product.Id, resp.Id)
		mockDB.ExpectationsWereMet()
	})

	t.Run("no products to delete", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		productRepo := new(MockProductRepo)
		recRepo := new(MockReceptionRepo)
		service := services.NewPvzService(pvzRepo, productRepo, recRepo, db)

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.Anything, mock.Anything).Return(&models.Pvz{}, nil)
		recRepo.On("GetByPvzId", mock.Anything, mock.Anything).Return(&models.Reception{
			Id:     "uuid",
			Status: reception.InProgressStatus,
		}, nil)
		productRepo.On("DeleteLast", mock.Anything, mock.Anything).Return((*models.Product)(nil), nil)
		mockDB.ExpectRollback()

		_, err := service.DeleteLastProduct("pvz123")

		require.Error(t, err)
		require.IsType(t, errors.ObjectHasNotSubObjects{}, err)
	})
}

func TestPvzService_CLoseLastReception(t *testing.T) {
	logger.Init("debug")

	t.Run("successful reception closing", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		recRepo := new(MockReceptionRepo)
		service := services.NewPvzService(pvzRepo, nil, recRepo, db)

		pvzID := "pvz123"
		rec := &models.Reception{
			Id:     "rec456",
			Status: reception.CloseStatus,
		}

		mockDB.ExpectBegin()
		pvzRepo.On("GetById", mock.Anything, pvzID).Return(&models.Pvz{}, nil)
		recRepo.On("GetByPvzId", mock.Anything, pvzID).Return(&models.Reception{
			Id:     "uuid",
			Status: reception.InProgressStatus,
		}, nil)
		recRepo.On("SetStatus", mock.Anything, reception.CloseStatus, pvzID).Return(rec, nil)
		mockDB.ExpectCommit()

		resp, err := service.CLoseLastReception(pvzID)

		require.NoError(t, err)
		require.Equal(t, rec.Status, resp.Status)
	})

}

func TestPvzService_ListWithFilterDate(t *testing.T) {
	logger.Init("debug")

	t.Run("successful list with products", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		defer db.Close()

		pvzRepo := new(MockPvzRepo)
		service := services.NewPvzService(pvzRepo, nil, nil, db)

		now := time.Now()
		req := pvz.ListRequest{
			StartDate: &now,
			EndDate:   &now,
			Page:      1,
			Limit:     10,
		}

		rawData := []pvz.RawList{
			{
				PvzId:           "pvz1",
				ReceptionId:     "rec1",
				ProductId:       sql.NullString{String: "prod1", Valid: true},
				ProductType:     sql.NullString{String: "электроника", Valid: true},
				ReceptionStatus: reception.InProgressStatus,
			},
		}

		pvzRepo.On("ListWithFilterDate", db, req, 0).Return(rawData, nil)

		result, err := service.ListWithFilterDate(req)

		require.NoError(t, err)
		require.Len(t, result, 1)
		require.Len(t, result[0].Receptions, 1)
		require.Len(t, result[0].Receptions[0].Products, 1)
	})

	t.Run("invalid date range", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		defer db.Close()

		service := services.NewPvzService(nil, nil, nil, db)
		start := time.Now()
		end := start.Add(-time.Hour)

		_, err := service.ListWithFilterDate(pvz.ListRequest{
			StartDate: &start,
			EndDate:   &end,
		})

		require.Error(t, err)
		require.IsType(t, errors.StartDateAfterEndDate{}, err)
	})

}
