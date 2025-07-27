package logger

import (
	"io"
	"log/slog"
)

func Configure(production bool, w io.Writer) *slog.Logger {
	var handler slog.Handler
	if production {
		handler = slog.NewJSONHandler(w, nil)
	} else {
		handler = slog.NewTextHandler(w, nil)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
