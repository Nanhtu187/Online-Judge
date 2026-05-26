package config

import (
	common "github.com/Nanhtu187/online-judge/src/packages/config"
)

type Config struct {
	Server         common.ServerConfig `mapstructure:"server" json:"server"`
	Kafka          common.KafkaConfig  `mapstructure:"kafka" json:"kafka"`
	ServerEndpoint string              `mapstructure:"server_endpoint" json:"server_endpoint"`
	IAMEndpoint    string              `mapstructure:"iam_endpoint" json:"iam_endpoint"`
	InternalKey    string              `mapstructure:"internal_api_key" json:"internal_api_key"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Server:         common.DefaultServerConfig(),
		Kafka:          common.DefaultKafkaConfig(),
		ServerEndpoint: "localhost:50051",
		IAMEndpoint:    "localhost:50050",
		InternalKey:    "secret",
	}

	if err := common.LoadConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
