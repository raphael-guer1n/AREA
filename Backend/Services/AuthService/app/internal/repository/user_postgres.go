package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/raphael-guer1n/AREA/AuthService/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(email, username, passwordHash string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(
		`INSERT INTO users (email, login, hashed_password)
         VALUES ($1, $2, $3)
         RETURNING id, email, login, hashed_password, created_at, updated_at`,
		email, username, passwordHash,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(
		`SELECT id, email, login, hashed_password, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(
		`SELECT id, email, login, hashed_password, created_at, updated_at FROM users WHERE login = $1`,
		username,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByEmailOrUsername(identifier string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(
		`SELECT id, email, login, hashed_password, created_at, updated_at
         FROM users WHERE email = $1 OR login = $1`,
		identifier,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByID(id int) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(
		`SELECT id, email, login, hashed_password, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) DeleteByID(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.Exec(`DELETE FROM user_service_profiles WHERE user_id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete user profiles: %w", err)
	}

	res, err := tx.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read delete result: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
