package repositories

import (
	"database/sql"
	"errors"
	"pvz/internal/models"
	"pvz/internal/models/pvz"
)

type PvzRepository interface {
	GetById(q Querier, id string) (*models.Pvz, error)
	Create(q Querier, reqPvz pvz.CreateRequest) (*models.Pvz, error)
	ListWithFilterDate(q Querier, reqPvz pvz.ListRequest, offset int) ([]pvz.RawList, error)
}

type pvzRepositoryPsql struct {
	db *sql.DB
}

func NewPvzRepository(db *sql.DB) PvzRepository {
	return &pvzRepositoryPsql{
		db: db,
	}
}

func (pr *pvzRepositoryPsql) GetById(q Querier, id string) (*models.Pvz, error) {
	query := `SELECT * FROM pvzs WHERE id = $1`

	var scanPvz models.Pvz
	err := q.QueryRow(query, id).Scan(&scanPvz.Id, &scanPvz.RegistrationDate, &scanPvz.City)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &scanPvz, nil
}

func (pr *pvzRepositoryPsql) Create(q Querier, reqPvz pvz.CreateRequest) (*models.Pvz, error) {
	query := `INSERT INTO pvzs(city) VALUES ($1) RETURNING id, registrationDate, city`

	var scanPvz models.Pvz
	err := q.QueryRow(query, reqPvz.City).Scan(&scanPvz.Id, &scanPvz.RegistrationDate, &scanPvz.City)
	if err != nil {
		return nil, err
	}
	return &scanPvz, nil
}

func (pr *pvzRepositoryPsql) ListWithFilterDate(q Querier, reqPvz pvz.ListRequest, offset int) ([]pvz.RawList, error) {
	query := `SELECT p.id, p.registrationDate, p.city,
					 r.id, r.createdAt, r.status,
					 pr.id, pr.receivedAt, pr.type
		FROM pvzs p
		LEFT JOIN receptions r ON r.pvzId = p.id
		LEFT JOIN products pr ON pr.receptionId = r.id
		WHERE r.createdAt BETWEEN $1 AND $2
		ORDER BY r.createdAt DESC
		OFFSET $3 LIMIT $4;
	`

	rows, err := q.Query(query, reqPvz.StartDate, reqPvz.EndDate, offset, reqPvz.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []pvz.RawList
	for rows.Next() {
		var row pvz.RawList
		err := rows.Scan(
			&row.PvzId, &row.PvzRegDate, &row.PvzCity,
			&row.ReceptionId, &row.ReceptionDate, &row.ReceptionStatus,
			&row.ProductId, &row.ProductDate, &row.ProductType,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}
