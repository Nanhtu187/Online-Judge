package config

import (
	common "github.com/Nanhtu187/online-judge/src/packages/config"
)

type Config struct {
	Server   common.ServerConfig   `mapstructure:"server" json:"server"`
	Database common.DatabaseConfig `mapstructure:"database" json:"database"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Server:   common.DefaultServerConfig(),
		Database: common.DefaultDatabaseConfig(),
	}

	if err := common.LoadConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
