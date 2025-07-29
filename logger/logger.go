package logger

import (
	"io"
	"log/slog"
)

func Configure(production bool, w io.Writer) *slog.Logger {
	var handler slog.Handler

	// if we're in production mode, make the log handler a json handler for observability
	if production {
		handler = slog.NewJSONHandler(w, nil)
	} else {
		// if we're in debug mode, text handler for easier log reading
		handler = slog.NewTextHandler(w, nil)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
