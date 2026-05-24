package config

import (
	common "github.com/Nanhtu187/online-judge/src/packages/config"
)

type Config struct {
	Server         common.ServerConfig `mapstructure:"server" json:"server"`
	Kafka          common.KafkaConfig  `mapstructure:"kafka" json:"kafka"`
	ServerEndpoint string              `mapstructure:"server_endpoint" json:"server_endpoint"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Server:         common.DefaultServerConfig(),
		Kafka:          common.DefaultKafkaConfig(),
		ServerEndpoint: "localhost:50051",
	}

	if err := common.LoadConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
