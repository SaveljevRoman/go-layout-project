package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerAddress       string      `json:"server_address"`
	CacheUpdateInterval int         `json:"cache_update_interval"` // в секундах
	MySQL               MySQLConfig `json:"mysql"`
	Redis               RedisConfig `json:"redis"`
}

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func Load() (*Config, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
