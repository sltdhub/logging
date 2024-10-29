package logging

import (
	"context"
	"log/slog"
	"os"
)

const (
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelFatal: "FATAL",
}

type Logger struct {
	*slog.Logger
}

func NewLogger() *Logger {
	logLevel := new(slog.LevelVar)

	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		logLevel.Set(slog.LevelDebug)
	case "INFO":
		logLevel.Set(slog.LevelInfo)
	case "WARN":
		logLevel.Set(slog.LevelWarn)
	case "ERROR":
		logLevel.Set(slog.LevelError)
	case "FATAL":
		logLevel.Set(LevelFatal)
	default:
		logLevel.Set(slog.LevelInfo)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	return &Logger{
		slog.New(handler),
	}
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.Log(context.TODO(), LevelFatal, msg, args...)
	os.Exit(1)
}
