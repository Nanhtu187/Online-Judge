package config

import "github.com/redis/go-redis/v9"

type RedisConfig struct {
	Address  string `json:"address" mapstructure:"address"`
	Password string `json:"password" mapstructure:"password"`
	DB       int    `json:"db" mapstructure:"db"`
}

func RedisDefaultConfig() RedisConfig {
	return RedisConfig{
		Address:  "localhost:6379",
		Password: "",
		DB:       0,
	}
}

func (rd RedisConfig) MustConnect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     rd.Address,
		Password: rd.Password,
		DB:       rd.DB,
	})
}
