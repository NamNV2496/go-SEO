package configs

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	AppPort        string `env:"APP_PORT" envDefault:"8080"`
	BuildUrlTool   string `env:"BUILD_URL_TOOL" envDefault:"template"` // template or regex or openai
	DatabaseConfig DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"root"`
	Password string `env:"DB_PASSWORD" envDefault:"root"`
	DBName   string `env:"DB_NAME" envDefault:"postgres"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

func LoadConfig() *Config {
	var conf Config
	if err := env.Parse(&conf); err != nil {
		panic("cannot load env")
	}
	return &conf
}
