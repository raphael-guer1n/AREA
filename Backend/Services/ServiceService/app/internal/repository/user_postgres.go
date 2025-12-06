package repository

import (
	"database/sql"

	"github.com/raphael-guer1n/AREA/ServiceService/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) List() ([]domain.User, error) {
	rows, err := r.db.Query(`SELECT id, email, first_name, last_name, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) Create(email, firstName, lastName string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(
		`INSERT INTO users (email, first_name, last_name)
         VALUES ($1, $2, $3)
         RETURNING id, email, first_name, last_name, created_at`,
		email, firstName, lastName,
	).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)
	return u, err
}
