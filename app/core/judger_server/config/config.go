package config

import (
	"bytes"
	"encoding/json"
	"strings"

	common "github.com/Nanhtu187/Online-Judge/app/common/config"
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig      `mapstructure:"database" json:"database" yaml:"database"`
	Server   ServerConfig        `mapstructure:"server" json:"server" yaml:"server"`
	Log      LogConfig           `mapstructure:"log" json:"log" yaml:"log"`
	Redis    RedisConfig         `mapstructure:"redis" json:"redis" yaml:"redis"`
	Jaeger   common.JaegerConfig `mapstructure:"jaeger" json:"jaeger" yaml:"jaeger"`
}

func Load() (*Config, error) {
	c := &Config{
		Server:   ServerDefaultConfig(),
		Log:      LogDefaultConfig(),
		Database: DatabaseDefaultConfig(),
		Redis:    RedisDefaultConfig(),
		Jaeger:   common.JaegerDefaultConfig(),
	}
	// --- hacking to load reflect structure config into env ----//
	viper.SetConfigType("json")
	configBuffer, err := json.Marshal(c)

	if err != nil {
		return nil, err
	}

	viper.ReadConfig(bytes.NewBuffer(configBuffer))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	// -- end of hacking --//
	viper.AutomaticEnv()
	err = viper.Unmarshal(c)
	return c, err
}
