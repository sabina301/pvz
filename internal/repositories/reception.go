package repositories

import (
	"database/sql"
	"errors"
	"pvz/internal/models"
)

type ReceptionRepository interface {
	Create(q Querier, req models.Reception) (*models.Reception, error)
	GetByPvzId(q Querier, pvzId string) (*models.Reception, error)
	SetStatus(q Querier, status, pvzId string) (*models.Reception, error)
}

type receptionRepositoryPsql struct {
	db *sql.DB
}

func NewReceptionRepository(db *sql.DB) ReceptionRepository {
	return &receptionRepositoryPsql{
		db: db,
	}
}

func (rr *receptionRepositoryPsql) GetByPvzId(q Querier, pvzId string) (*models.Reception, error) {
	query := `SELECT * FROM receptions WHERE pvzId = $1 ORDER BY createdAt DESC LIMIT 1`
	var scanReception models.Reception
	err := q.QueryRow(query, pvzId).Scan(&scanReception.Id, &scanReception.DateTime, &scanReception.PvzId, &scanReception.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &scanReception, nil
}

func (rr *receptionRepositoryPsql) Create(q Querier, req models.Reception) (*models.Reception, error) {
	query := `INSERT INTO receptions (pvzId, status) VALUES ($1,$2) RETURNING id, createdAt, pvzId, status`
	var reception models.Reception
	err := q.QueryRow(query, req.PvzId, req.Status).Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reception, nil
}

func (rr *receptionRepositoryPsql) SetStatus(q Querier, status, pvzId string) (*models.Reception, error) {
	query := `UPDATE receptions SET status = $1 WHERE pvzId = $2 RETURNING id, createdAt, pvzId, status`
	var reception models.Reception
	err := q.QueryRow(query, status, pvzId).Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reception, nil
}
