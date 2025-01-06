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

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
}

type Config struct {
	Kafka *KafkaConfig `mapstructure:"kafka"`
	DB    *DBConfig    `mapstructure:"db"`
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
		DB:    getDBConfig(),
	}, nil
}

func getKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		URL:            viper.GetString("kafka.url"),
		GroupID:        viper.GetString("kafka.group_id"),
		ConsumerTopics: viper.GetStringSlice("kafka.consumer_topics"),
	}
}

func getDBConfig() *DBConfig {
	return &DBConfig{
		Host:     viper.GetString("db.host"),
		Database: viper.GetString("db.database"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Port:     viper.GetInt("db.port"),
	}
}
