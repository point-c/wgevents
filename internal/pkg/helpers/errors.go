package helpers

import (
	"log/slog"
	"os"
)

var exit = os.Exit

func Must[T any](t T, err error) T {
	if err != nil {
		slog.Error("failed to get value", "error", err)
		exit(1)
	}
	return t
}

func Check(err error) {
	if err != nil {
		slog.Error("operation failed", "error", err)
		exit(1)
	}
}
