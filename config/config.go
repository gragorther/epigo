package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Production    bool       `env:"PROD" env-description:"whether the server is in prod mode"`
	LogLevel      slog.Level `env:"LOG_LEVEL"`
	AdminUsername string     `env:"ADMIN_USERNAME"`
	AdminPassword string     `env:"ADMIN_PASSWORD"`
	JWTSecret     string     `env:"JWT_SECRET"`
	RedisAddress  string     `env:"REDIS_ADDRESS"`
	DatabaseURL   string     `env:"DATABASE_URL"`
	Email         EmailConfig
	BASE_URL      string `env:"BASE_URL" env-description:"the base url of the app, e.g. https://afterwill.life"`
}

func Get() (Config, error) {
	var conf Config
	err := cleanenv.ReadEnv(&conf)

	return conf, err
}

type EmailConfig struct {
	From string `env:"EMAIL_FROM"`
}
