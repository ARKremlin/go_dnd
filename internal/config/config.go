package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	DB       DB
	Redis    Redis
	JWT      JWT
	Telegram Telegram
}

type Telegram struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
}

type App struct {
	Port     string `env:"APP_PORT" env-required:"true"`
	Loglevel string `env:"LOG_LEVEL" env-required:"true"`
}

type DB struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     string `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Name     string `env:"POSTGRES_DB" env-required:"true"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST" env-required:"true"`
	Port     string `env:"REDIS_PORT" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
}

type JWT struct {
	Secret string        `env:"JWT_SECRET" env-required:"true"`
	TTL    time.Duration `env:"JWT_TTL" env-default:"15m"`
}

func (d DB) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name)
}

func Load() (*Config, error) {
	var cfg Config

	_ = godotenv.Load()

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env %w", err)
	}
	return &cfg, nil
}
