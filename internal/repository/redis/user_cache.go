package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

type UserCache struct {
	client *redis.Client
}

func NewUserCache(client *redis.Client) *UserCache {
	return &UserCache{
		client: client,
	}
}

func (c *UserCache) getUserKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (c *UserCache) GetByID(ctx context.Context, id int64) (*models.User, error) {
	key := c.getUserKey(id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Кеш пуст
		}
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *UserCache) Set(ctx context.Context, user *models.User, expiration time.Duration) error {
	key := c.getUserKey(user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *UserCache) Delete(ctx context.Context, id int64) error {
	key := c.getUserKey(id)
	return c.client.Del(ctx, key).Err()
}

func (c *UserCache) SetAllUsers(ctx context.Context, users []*models.User, expiration time.Duration) error {
	for _, user := range users {
		if err := c.Set(ctx, user, expiration); err != nil {
			return err
		}
	}

	// Обновляем список всех идентификаторов пользователей
	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = fmt.Sprintf("%d", user.ID)
	}

	return c.client.Set(ctx, "users:all", userIDs, expiration).Err()
}
