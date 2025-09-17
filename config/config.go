package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Production               bool   `env:"PROD" env-description:"whether the server is in prod mode"`
	AdminUsername            string `env:"ADMIN_USERNAME"`
	AdminPassword            string `env:"ADMIN_PASSWORD"`
	JWTSecret                string `env:"JWT_SECRET"`
	DatabaseURL              string `env:"DATABASE_URL"`
	Email                    EmailConfig
	Redis                    RedisConfig
	BaseURL                  string        `env:"BASE_URL" env-description:"the base url of the app, e.g. https://afterwill.life"`
	GinMode                  string        `env:"GIN_MODE"`
	MinDurationBetweenEmails time.Duration `env:"MIN_DURATION_BETWEEN_EMAILS"`
}

func Get() (Config, error) {
	var conf Config
	err := cleanenv.ReadEnv(&conf)
	return conf, err
}

type EmailConfig struct {
	FromFormat string `env:"EMAIL_FROM_NAME"`
	From       string `env:"EMAIL_FROM"`
	Host       string `env:"EMAIL_HOST"`
	Port       int    `env:"EMAIL_PORT"`
	Password   string `env:"EMAIL_PASSWORD"`
	Username   string `env:"EMAIL_USERNAME" env-description:"the username used for authenticating with the mail server"`
}

type RedisConfig struct {
	Address  string `env:"REDIS_ADDRESS"`
	Password string `env:"REDIS_PASSWORD"`
	Username string `env:"REDIS_USERNAME"`
	DB       int    `env:"REDIS_DB" env-description:"the redis db number to select"`
}
