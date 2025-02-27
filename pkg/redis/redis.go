package redis

import (
	"github.com/SaveljevRoman/go-layout-project/internal/config"
	"github.com/go-redis/redis/v8"
)

func NewConnection(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return client, nil
}
