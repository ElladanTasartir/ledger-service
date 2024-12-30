package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type KafkaConfig struct {
	URL string `mapstructure:"url"`
}

type Config struct {
	Kafka *KafkaConfig `mapstructure:"kafka"`
}

func NewConfig(path string) (*Config, error) {
	viper.SetConfigFile(fmt.Sprintf("%s/%s.yml", path, getEnv()))

	return getConfig()
}

func getEnv() string {
	return "local"
}

func getConfig() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		Kafka: getKafkaConfig(),
	}, nil
}

func getKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		URL: viper.GetString("kafka.url"),
	}
}
