package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Postgres struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

type Server struct {
	Port int
}

type Config struct {
	DB     Postgres
	Server Server
}

func New(dirname, filename string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(dirname)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := envconfig.Process("db", &cfg.DB); err != nil {
		return nil, err
	}

	return cfg, nil
}
