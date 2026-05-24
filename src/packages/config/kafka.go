package config

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers" json:"brokers"`
}

func DefaultKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: []string{"localhost:9092"},
	}
}
