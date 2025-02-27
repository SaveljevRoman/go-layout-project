package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

type ProductCache struct {
	client *redis.Client
}

func NewProductCache(client *redis.Client) *ProductCache {
	return &ProductCache{
		client: client,
	}
}

func (c *ProductCache) getProductKey(id int64) string {
	return fmt.Sprintf("product:%d", id)
}

func (c *ProductCache) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	key := c.getProductKey(id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Кеш пуст
		}
		return nil, err
	}

	var product models.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *ProductCache) Set(ctx context.Context, product *models.Product, expiration time.Duration) error {
	key := c.getProductKey(product.ID)
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *ProductCache) Delete(ctx context.Context, id int64) error {
	key := c.getProductKey(id)
	return c.client.Del(ctx, key).Err()
}

func (c *ProductCache) SetAllProducts(ctx context.Context, products []*models.Product, expiration time.Duration) error {
	for _, product := range products {
		if err := c.Set(ctx, product, expiration); err != nil {
			return err
		}
	}

	// Обновляем список всех идентификаторов продуктов
	productIDs := make([]string, len(products))
	for i, product := range products {
		productIDs[i] = fmt.Sprintf("%d", product.ID)
	}

	return c.client.Set(ctx, "products:all", productIDs, expiration).Err()
}
