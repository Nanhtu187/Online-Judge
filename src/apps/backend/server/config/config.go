package config

import (
	common "github.com/Nanhtu187/online-judge/src/packages/config"
)

type Config struct {
	Server      common.ServerConfig   `mapstructure:"server" json:"server"`
	Database    common.DatabaseConfig `mapstructure:"database" json:"database"`
	Kafka       common.KafkaConfig    `mapstructure:"kafka" json:"kafka"`
	IAMEndpoint string                `mapstructure:"iam_endpoint" json:"iam_endpoint"`
	InternalKey string                `mapstructure:"internal_api_key" json:"internal_api_key"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Server:      common.DefaultServerConfig(),
		Database:    common.DefaultDatabaseConfig(),
		Kafka:       common.DefaultKafkaConfig(),
		IAMEndpoint: "localhost:50050",
		InternalKey: "secret",
	}

	if err := common.LoadConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
