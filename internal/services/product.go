package services

import (
	"database/sql"
	"pvz/internal/logger"
	"pvz/internal/models"
	"pvz/internal/models/product"
	"pvz/internal/models/reception"
	"pvz/internal/repositories"
	"pvz/pkg/errors"
)

type ProductService interface {
	AddInReception(req product.AddInReceptionRequest) (product.AddInReceptionResponse, error)
}

type productServiceImpl struct {
	productRepo   repositories.ProductRepository
	pvzRepo       repositories.PvzRepository
	receptionRepo repositories.ReceptionRepository
	conn          *sql.DB
}

func NewProductService(productRepo repositories.ProductRepository, pvzRepo repositories.PvzRepository, recRepo repositories.ReceptionRepository, conn *sql.DB) ProductService {
	return &productServiceImpl{
		productRepo:   productRepo,
		pvzRepo:       pvzRepo,
		receptionRepo: recRepo,
		conn:          conn,
	}
}
func (ps *productServiceImpl) AddInReception(req product.AddInReceptionRequest) (product.AddInReceptionResponse, error) {
	log := logger.Log.With("pvz_id", req.PvzId, "product_type", req.Type)
	log.Info("starting AddInReception")

	tx, err := ps.conn.Begin()
	if err != nil {
		log.Error("failed to begin transaction", "err", err)
		return product.AddInReceptionResponse{}, errors.NewInternalError()
	}
	defer tx.Rollback()

	pvz, err := ps.pvzRepo.GetById(tx, req.PvzId)
	if err != nil {
		log.Error("failed to fetch pvz", "err", err)
		return product.AddInReceptionResponse{}, errors.NewInternalError()
	}
	if pvz == nil {
		log.Warn("pvz not found")
		return product.AddInReceptionResponse{}, errors.NewObjectNotFound("pvz")
	}

	rec, err := ps.receptionRepo.GetByPvzId(tx, req.PvzId)
	if err != nil {
		log.Error("failed to fetch reception", "err", err)
		return product.AddInReceptionResponse{}, errors.NewInternalError()
	}
	if rec == nil {
		log.Warn("reception not found for pvz")
		return product.AddInReceptionResponse{}, errors.NewObjectNotFound("pvz")
	}
	if rec.Status != reception.InProgressStatus {
		log.Warn("reception is not in progress", "reception_id", rec.Id, "status", rec.Status)
		return product.AddInReceptionResponse{}, errors.NewReceptionIsNotInProgress(rec.Id)
	}

	reqProduct := models.Product{
		Type:        req.Type,
		ReceptionId: rec.Id,
	}

	productResp, err := ps.productRepo.AddInReception(tx, reqProduct)
	if err != nil {
		log.Error("failed to add product to reception", "err", err)
		return product.AddInReceptionResponse{}, errors.NewInternalError()
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed to commit transaction", "err", err)
		return product.AddInReceptionResponse{}, errors.NewInternalError()
	}

	log.Info("product successfully added to reception", "product_id", productResp.Id, "reception_id", productResp.ReceptionId)

	return product.AddInReceptionResponse{
		Id:          productResp.Id,
		DateTime:    productResp.DateTime,
		Type:        productResp.Type,
		ReceptionId: productResp.ReceptionId,
	}, nil
}
