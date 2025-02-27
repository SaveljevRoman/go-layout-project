package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	product := &models.Product{}
	query := "SELECT * FROM products WHERE id = ?"
	err := r.db.GetContext(ctx, product, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Продукт не найден
		}
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]*models.Product, error) {
	products := []*models.Product{}
	query := "SELECT * FROM products"
	err := r.db.SelectContext(ctx, &products, query)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) (int64, error) {
	query := "INSERT INTO products (name, description, price, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())"
	result, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Quantity)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	query := "UPDATE products SET name = ?, description = ?, price = ?, quantity = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Quantity, product.ID)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM products WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
