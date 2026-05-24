package config

import (
	"fmt"
)

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver" json:"driver"`
	Host            string `mapstructure:"host" json:"host"`
	Port            string `mapstructure:"port" json:"port"`
	User            string `mapstructure:"user" json:"user"`
	Password        string `mapstructure:"password" json:"password"`
	Name            string `mapstructure:"name" json:"name"`
	MigrationSource string `mapstructure:"migration_source" json:"migration_source"`
}

func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Driver:          "mysql",
		Host:            "localhost",
		Port:            "3306",
		MigrationSource: "file://src/apps/backend/server/internal/migrations",
	}
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

func (c DatabaseConfig) MigrationDSN() string {
	return fmt.Sprintf("%s://%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Driver, c.User, c.Password, c.Host, c.Port, c.Name)
}
