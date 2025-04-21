package services_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"

	"pvz/internal/logger"
	"pvz/internal/models"
	"pvz/internal/models/reception"
	"pvz/internal/services"
	"pvz/pkg/errors"
)

func TestReceptionService_Create(t *testing.T) {
	logger.Init("debug")

	t.Run("successful reception creation", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		recRepo := new(MockReceptionRepo)
		pvzRepo := new(MockPvzRepo)
		service := services.NewReceptionService(recRepo, pvzRepo, db)

		pvzID := "550e8400-e29b-41d4-a716-446655440000"
		req := reception.CreateRequest{PvzId: pvzID}

		mockDB.ExpectBegin()
		recRepo.On("GetByPvzId", mock.AnythingOfType("*sql.Tx"), pvzID).
			Return((*models.Reception)(nil), nil).
			Once()
		pvzRepo.On("GetById", mock.AnythingOfType("*sql.Tx"), pvzID).
			Return(&models.Pvz{Id: pvzID}, nil).
			Once()
		recRepo.On("Create", mock.AnythingOfType("*sql.Tx"), mock.MatchedBy(func(r models.Reception) bool {
			return r.PvzId == pvzID && r.Status == reception.InProgressStatus
		})).Return(&models.Reception{Id: "rec123"}, nil).
			Once()
		mockDB.ExpectCommit()

		resp, err := service.Create(req)

		require.NoError(t, err)
		require.Equal(t, "rec123", resp.Id)
		recRepo.AssertExpectations(t)
		pvzRepo.AssertExpectations(t)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("transaction begin error", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		recRepo := new(MockReceptionRepo)
		pvzRepo := new(MockPvzRepo)
		service := services.NewReceptionService(recRepo, pvzRepo, db)
		mockDB.ExpectBegin().WillReturnError(errors.NewInternalError())

		_, err := service.Create(reception.CreateRequest{PvzId: "a"})

		require.Error(t, err)
		require.IsType(t, errors.InternalError{}, err)
	})

	t.Run("existing active reception", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		recRepo := new(MockReceptionRepo)
		service := services.NewReceptionService(recRepo, nil, db)

		pvzID := "550e8400-e29b-41d4-a716-446655440000"
		req := reception.CreateRequest{PvzId: pvzID}

		mockDB.ExpectBegin()
		recRepo.On("GetByPvzId", mock.AnythingOfType("*sql.Tx"), pvzID).
			Return(&models.Reception{Status: reception.InProgressStatus}, nil).
			Once()
		mockDB.ExpectRollback()

		_, err := service.Create(req)

		require.Error(t, err)
		require.IsType(t, errors.ReceptionIsNotClosed{}, err)
		recRepo.AssertExpectations(t)
	})

	t.Run("database error on reception check", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		recRepo := new(MockReceptionRepo)
		service := services.NewReceptionService(recRepo, nil, db)

		req := reception.CreateRequest{PvzId: "pvz123"}

		mockDB.ExpectBegin()
		recRepo.On("GetByPvzId", mock.AnythingOfType("*sql.Tx"), "pvz123").
			Return((*models.Reception)(nil), errors.NewInternalError()).
			Once()
		mockDB.ExpectRollback()

		_, err := service.Create(req)

		require.Error(t, err)
		require.IsType(t, errors.InternalError{}, err)
		recRepo.AssertExpectations(t)
	})

	t.Run("transaction begin error", func(t *testing.T) {
		db, mockDB, _ := sqlmock.New()
		defer db.Close()

		service := services.NewReceptionService(nil, nil, db)
		mockDB.ExpectBegin().WillReturnError(errors.NewInternalError())

		_, err := service.Create(reception.CreateRequest{})

		require.Error(t, err)
		require.IsType(t, errors.InternalError{}, err)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})
}
