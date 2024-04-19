package config

// JaegerConfig ...
type JaegerConfig struct {
	Host          string   `json:"host" mapstructure:"host"`
	Port          uint16   `json:"port" mapstructure:"port"`
	Ratio         float64  `json:"ratio" mapstructure:"ratio"`
	ExcludedPaths []string `json:"excluded_paths" mapstructure:"excluded_paths"`
	Enabled       bool     `json:"enabled" mapstructure:"enabled"`
}

func JaegerDefaultConfig() JaegerConfig {
	return JaegerConfig{
		Host:          "localhost",
		Port:          6831,
		Ratio:         1,
		ExcludedPaths: []string{},
		Enabled:       true,
	}
}
