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
}
type PostgresConfig struct {
}

func Get() (Config, error) {
	var conf Config
	err := cleanenv.ReadEnv(&conf)

	return conf, err
}
