package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HttpPort string `mapstructure:"HTTP_PORT"`

	PsqlHost     string `mapstructure:"PSQL_HOST"`
	PsqlName     string `mapstructure:"PSQL_NAME"`
	PsqlUser     string `mapstructure:"PSQL_USER"`
	PsqlPassword string `mapstructure:"PSQL_PASSWORD"`
	PsqlSSLMode  string `mapstructure:"PSQL_SSLMODE"`
}

func InitConfig() (*Config, error) {
	var cfg Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
