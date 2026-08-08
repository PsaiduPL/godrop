package log

import (
	"fmt"
	"log/slog"
)

const (
	DefaultLogLevel    = slog.LevelWarn
	DefaultLogLevelStr = "WARN"
)

func MapStringToSlogLevel(levelStr string) (slog.Level, error) {
	var (
		predicatedLevel slog.Level
		found           = true
	)
	switch levelStr {
	case "DEBUG":
		predicatedLevel = slog.LevelDebug
	case "INFO":
		predicatedLevel = slog.LevelInfo
	case "WARN":
		predicatedLevel = slog.LevelWarn
	case "ERROR":
		predicatedLevel = slog.LevelError
	default:
		found = false
	}
	if found {
		return predicatedLevel, nil
	}
	return slog.LevelInfo, fmt.Errorf("Invalid log level one of DEGUG,INFO,WARN,ERROR avaiable")
}

func SetLogLevel(level slog.Level) {
	slog.SetLogLoggerLevel(level)
}
