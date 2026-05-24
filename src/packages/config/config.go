package config

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(cfg any) error {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AutomaticEnv()

	buf, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	v.SetConfigType("json")
	if err := v.ReadConfig(bytes.NewBuffer(buf)); err != nil {
		return err
	}

	return v.Unmarshal(cfg)
}
