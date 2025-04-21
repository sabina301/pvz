package services

import (
	"database/sql"
	"pvz/internal/logger"
	"pvz/internal/models"
	"pvz/internal/models/reception"
	"pvz/internal/repositories"
	"pvz/pkg/errors"
)

type ReceptionService interface {
	Create(req reception.CreateRequest) (reception.CreateResponse, error)
}

type receptionServiceImpl struct {
	receptionRepo repositories.ReceptionRepository
	pvzRepo       repositories.PvzRepository
	conn          *sql.DB
}

func NewReceptionService(receptionRepo repositories.ReceptionRepository, pvzRepo repositories.PvzRepository, conn *sql.DB) ReceptionService {
	return &receptionServiceImpl{
		receptionRepo: receptionRepo,
		pvzRepo:       pvzRepo,
		conn:          conn,
	}
}
func (rs *receptionServiceImpl) Create(req reception.CreateRequest) (reception.CreateResponse, error) {
	log := logger.Log.With("pvzId", req.PvzId)
	log.Info("starting Create reception")

	tx, err := rs.conn.Begin()
	if err != nil {
		log.Error("failed to begin transaction", "err", err)
		return reception.CreateResponse{}, errors.InternalError{}
	}
	defer tx.Rollback()

	rec, err := rs.receptionRepo.GetByPvzId(tx, req.PvzId)
	if err != nil {
		log.Error("failed to get reception by pvzId", "err", err)
		return reception.CreateResponse{}, errors.NewInternalError()
	}
	if rec != nil {
		log.Warn("reception already exists and is not closed", "pvzId", req.PvzId)
		return reception.CreateResponse{}, errors.NewReceptionIsNotClosed(req.PvzId)
	}

	pvz, err := rs.pvzRepo.GetById(tx, req.PvzId)
	if err != nil {
		log.Error("failed to get pvz by id", "err", err)
		return reception.CreateResponse{}, errors.NewInternalError()
	}
	if pvz == nil {
		log.Warn("pvz not found", "pvzId", req.PvzId)
		return reception.CreateResponse{}, errors.NewObjectNotFound("pvz")
	}

	newRec := models.Reception{
		PvzId:  req.PvzId,
		Status: reception.InProgressStatus,
	}

	recRep, err := rs.receptionRepo.Create(tx, newRec)
	if err != nil {
		log.Error("failed to create reception", "err", err)
		return reception.CreateResponse{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error("failed to commit transaction", "err", err)
		return reception.CreateResponse{}, errors.NewInternalError()
	}

	log.Info("reception created successfully", "receptionId", recRep.Id, "status", recRep.Status)
	return reception.CreateResponse{
		Id:       recRep.Id,
		DateTime: recRep.DateTime,
		PvzId:    newRec.PvzId,
		Status:   newRec.Status,
	}, nil
}
