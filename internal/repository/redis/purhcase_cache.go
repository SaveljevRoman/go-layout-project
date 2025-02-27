package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

type PurchaseCache struct {
	client *redis.Client
}

func NewPurchaseCache(client *redis.Client) *PurchaseCache {
	return &PurchaseCache{
		client: client,
	}
}

func (c *PurchaseCache) getPurchaseKey(id int64) string {
	return fmt.Sprintf("purchase:%d", id)
}

func (c *PurchaseCache) getUserPurchasesKey(userID int64) string {
	return fmt.Sprintf("user:%d:purchases", userID)
}

func (c *PurchaseCache) GetByID(ctx context.Context, id int64) (*models.Purchase, error) {
	key := c.getPurchaseKey(id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Кеш пуст
		}
		return nil, err
	}

	var purchase models.Purchase
	if err := json.Unmarshal(data, &purchase); err != nil {
		return nil, err
	}

	return &purchase, nil
}

func (c *PurchaseCache) Set(ctx context.Context, purchase *models.Purchase, expiration time.Duration) error {
	key := c.getPurchaseKey(purchase.ID)
	data, err := json.Marshal(purchase)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *PurchaseCache) Delete(ctx context.Context, id int64) error {
	key := c.getPurchaseKey(id)
	return c.client.Del(ctx, key).Err()
}

func (c *PurchaseCache) SetUserPurchases(ctx context.Context, userID int64, purchases []*models.Purchase, expiration time.Duration) error {
	key := c.getUserPurchasesKey(userID)
	data, err := json.Marshal(purchases)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *PurchaseCache) GetUserPurchases(ctx context.Context, userID int64) ([]*models.Purchase, error) {
	key := c.getUserPurchasesKey(userID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Кеш пуст
		}
		return nil, err
	}

	var purchases []*models.Purchase
	if err := json.Unmarshal(data, &purchases); err != nil {
		return nil, err
	}

	return purchases, nil
}
