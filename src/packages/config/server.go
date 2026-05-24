package config

import (
	"fmt"
)

type ServerConfig struct {
	HTTPPort string `mapstructure:"http_port" json:"http_port"`
	GrpcPort string `mapstructure:"grpc_port" json:"grpc_port"`
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		HTTPPort: "8080",
		GrpcPort: "50050",
	}
}

func (c ServerConfig) HTTPAddr() string {
	return fmt.Sprintf(":%s", c.HTTPPort)
}

func (c ServerConfig) GrpcAddr() string {
	return fmt.Sprintf(":%s", c.GrpcPort)
}
