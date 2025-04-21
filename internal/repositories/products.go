package repositories

import (
	"database/sql"
	"errors"
	"pvz/internal/models"
)

type ProductRepository interface {
	AddInReception(q Querier, reqProduct models.Product) (*models.Product, error)
	DeleteLast(q Querier, recId string) (*models.Product, error)
}

type productRepositoryPsql struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepositoryPsql{
		db: db,
	}
}

func (pr *productRepositoryPsql) AddInReception(q Querier, reqProduct models.Product) (*models.Product, error) {
	query := `INSERT INTO products (receptionId, type) VALUES ($1, $2) RETURNING id, receivedAt, receptionId, type`
	var product models.Product
	err := q.QueryRow(query, reqProduct.ReceptionId, reqProduct.Type).Scan(&product.Id, &product.DateTime, &product.ReceptionId, &product.Type)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *productRepositoryPsql) DeleteLast(q Querier, recId string) (*models.Product, error) {
	query := `DELETE FROM products WHERE id = (SELECT id FROM products WHERE receptionId = $1 ORDER BY receivedAt DESC LIMIT 1) RETURNING id, receivedAt, receptionId, type`
	var product models.Product
	err := q.QueryRow(query, recId).Scan(&product.Id, &product.DateTime, &product.ReceptionId, &product.Type)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}
