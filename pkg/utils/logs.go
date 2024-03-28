package utils

import (
	"log/slog"
	"os"
)

func HandleFatal(err error, msg string) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}
