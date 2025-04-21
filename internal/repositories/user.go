package repositories

import (
	"database/sql"
	"errors"
	"pvz/internal/models"
)

type UserRepository interface {
	GetByEmail(q Querier, email string) (*models.User, error)
	Create(q Querier, user models.User) (string, error)
}

type userRepositoryPsql struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryPsql{
		db: db,
	}
}

func (ur *userRepositoryPsql) GetByEmail(q Querier, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, role FROM users WHERE email = $1`

	var user models.User
	err := q.QueryRow(query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepositoryPsql) Create(q Querier, user models.User) (string, error) {
	query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id`

	err := q.QueryRow(query, user.Email, user.PasswordHash, user.Role).Scan(&user.Id)
	if err != nil {
		return "", err
	}
	return user.Id, nil
}
