package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Host           string `env:"HOST"`
	Port           string `env:"PORT"`
	DB             string `env:"DB"`
	User           string `env:"USER"`
	Password       string `env:"PASSWORD"`
	SSLMode        string `env:"SSLMODE"`
	WorkerInterval int    `env:"CACHE_WORKER_INTERVAL"`
	AdminAuth      string `env:"ADMIN_AUTH"`
	UserAuth       string `env:"USER_AUTH"`
	JWTSecret      string `env:"JWT_SECRET"`
}

func NewConfig() *Config {
	return &Config{}
}
func MustLoad() *Config {
	cfg := NewConfig()

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	return cfg

}
