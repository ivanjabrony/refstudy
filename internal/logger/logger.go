package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type LoggerLevel string

var (
	Debug LoggerLevel = "debug"
	Prod  LoggerLevel = "prod"
	Test  LoggerLevel = "test"
)

const (
	LogFormatText = "text"
	LogFormatJson = "json"
)

type MyLogger struct {
	*slog.Logger
}

func New(level LoggerLevel, format string) *MyLogger {

	var loggerLevel slog.Level

	switch level {
	case Debug:
		loggerLevel = slog.LevelDebug
	case Prod:
		loggerLevel = slog.LevelInfo
	case Test:
		loggerLevel = slog.LevelDebug
	}

	var logger *slog.Logger
	switch format {
	case LogFormatJson:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: loggerLevel}))
	case LogFormatText:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: loggerLevel}))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: loggerLevel}))
		logger.Warn(fmt.Sprintf("unsupported logging format %s, using default format instead", format))
	}

	return &MyLogger{logger}
}

func (l *MyLogger) WrapError(msg string, err error, args ...any) {
	args = append(args, slog.String("error", err.Error()))
	l.Error(msg, args...)
}
