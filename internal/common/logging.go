package common

import (
	"log/slog"
	"os"
)

// Configure logging for `gambol` client.
func ConfigureLogger(verbose bool) {
	loggingLevel := new(slog.LevelVar)

	// Set log level to `INFO` unless `verbose` is true. `DEBUG` is the most verbose.
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	// Configure a new log handler and make it default logger.
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	loggingLevel.Set(level)
}
