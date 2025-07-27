package config

import "log/slog"

type Config struct {
	Production bool
	LogLevel   slog.Level
}

func GetConfig(getenv func(key string) string) {

}
