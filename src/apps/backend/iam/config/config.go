package config

import (
	common "github.com/Nanhtu187/online-judge/src/packages/config"
)

type Config struct {
	Server    common.ServerConfig   `mapstructure:"server" json:"server"`
	Database  common.DatabaseConfig `mapstructure:"database" json:"database"`
	JWTSecret string                `mapstructure:"jwt_secret" json:"jwt_secret"`
}

func Load() (*Config, error) {
	dbCfg := common.DefaultDatabaseConfig()
	dbCfg.MigrationSource = "file://src/apps/backend/iam/internal/migrations"

	cfg := &Config{
		Server:    common.DefaultServerConfig(),
		Database:  dbCfg,
		JWTSecret: "secret",
	}

	if err := common.LoadConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
