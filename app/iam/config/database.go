package config

type DatabaseConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	User     string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
	Database string `json:"database" mapstructure:"database"`
}

func DatabaseDefaultConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "1",
		Database: "iam",
	}
}
