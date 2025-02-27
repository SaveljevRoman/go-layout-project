package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := "SELECT * FROM users WHERE id = ?"
	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	users := []*models.User{}
	query := "SELECT * FROM users"
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	query := "INSERT INTO users (username, email, created_at, updated_at) VALUES (?, ?, NOW(), NOW())"
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET username = ?, email = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
