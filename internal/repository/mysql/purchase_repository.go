package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"github.com/jmoiron/sqlx"
)

type PurchaseRepository struct {
	db *sqlx.DB
}

func NewPurchaseRepository(db *sqlx.DB) *PurchaseRepository {
	return &PurchaseRepository{
		db: db,
	}
}

func (r *PurchaseRepository) GetByID(ctx context.Context, id int64) (*models.Purchase, error) {
	purchase := &models.Purchase{}
	query := "SELECT * FROM purchases WHERE id = ?"
	err := r.db.GetContext(ctx, purchase, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Покупка не найдена
		}
		return nil, err
	}
	return purchase, nil
}

func (r *PurchaseRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Purchase, error) {
	purchases := []*models.Purchase{}
	query := "SELECT * FROM purchases WHERE user_id = ?"
	err := r.db.SelectContext(ctx, &purchases, query, userID)
	if err != nil {
		return nil, err
	}
	return purchases, nil
}

func (r *PurchaseRepository) Create(ctx context.Context, purchase *models.Purchase) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Проверяем, достаточно ли товара на складе
	var currentQuantity int
	err = tx.GetContext(ctx, &currentQuantity, "SELECT quantity FROM products WHERE id = ?", purchase.ProductID)
	if err != nil {
		return 0, err
	}

	if currentQuantity < purchase.Quantity {
		return 0, errors.New("недостаточное количество товара")
	}

	// Обновляем количество товара
	_, err = tx.ExecContext(ctx, "UPDATE products SET quantity = quantity - ?, updated_at = NOW() WHERE id = ?",
		purchase.Quantity, purchase.ProductID)
	if err != nil {
		return 0, err
	}

	// Создаем запись о покупке
	query := `INSERT INTO purchases (user_id, product_id, quantity, total_price, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
	result, err := tx.ExecContext(ctx, query,
		purchase.UserID, purchase.ProductID, purchase.Quantity, purchase.TotalPrice, purchase.Status)
	if err != nil {
		return 0, err
	}

	purchaseID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return purchaseID, nil
}

func (r *PurchaseRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := "UPDATE purchases SET status = ?, updated_at = NOW() WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *PurchaseRepository) GetAll(ctx context.Context) ([]*models.Purchase, error) {
	purchases := []*models.Purchase{}
	query := "SELECT * FROM purchases"
	err := r.db.SelectContext(ctx, &purchases, query)
	if err != nil {
		return nil, err
	}
	return purchases, nil
}
