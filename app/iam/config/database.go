package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type DatabaseConfig struct {
	Type               string `json:"type" mapstructure:"type"`
	Host               string `json:"host" mapstructure:"host"`
	Port               string `json:"port" mapstructure:"port"`
	User               string `json:"user" mapstructure:"user"`
	Password           string `json:"password" mapstructure:"password"`
	Database           string `json:"database" mapstructure:"database"`
	MaxConnections     int    `json:"max_connections" mapstructure:"max_connections"`
	MaxIdleConnections int    `json:"max_idle_connections" mapstructure:"max_idle_connections"`
}

func DatabaseDefaultConfig() DatabaseConfig {
	return DatabaseConfig{
		Type:               "mysql",
		Host:               "localhost",
		Port:               "3306",
		User:               "root",
		Password:           "1",
		Database:           "iam",
		MaxConnections:     30,
		MaxIdleConnections: 30,
	}
}

// DSN returns data source name
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
}

func (c DatabaseConfig) MustConnect() *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.DSN()))
	if err != nil {
		log.Fatalf("Error when connect database: %s", err.Error())
	}
	DB, err := db.DB()
	DB.SetMaxOpenConns(c.MaxConnections)
	DB.SetMaxIdleConns(c.MaxIdleConnections)
	DB.SetConnMaxIdleTime(4 * time.Hour)
	return db
}
