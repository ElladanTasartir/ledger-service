package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type KafkaConfig struct {
	URL            string   `mapstructure:"url"`
	GroupID        string   `mapstructure:"group_id"`
	ConsumerTopics []string `mapstructure:"consumer_topics"`
}

type Config struct {
	Kafka *KafkaConfig `mapstructure:"kafka"`
}

func NewConfigModule(path string) fx.Option {
	return fx.Options(
		fx.Provide(func() (*Config, error) {
			return NewConfig(path)
		}),
	)
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
		URL:            viper.GetString("kafka.url"),
		GroupID:        viper.GetString("kafka.group_id"),
		ConsumerTopics: viper.GetStringSlice("kafka.consumer_topics"),
	}
}
