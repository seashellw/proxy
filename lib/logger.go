package lib

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

func NewLogger() *Logger {
	logger := slog.NewWithHandlers(
		handler.NewBuilder().
			WithLogfile("./tmp/error.log").
			WithLogLevels(slog.DangerLevels).
			WithBuffSize(1024*10).
			WithBuffMode(handler.BuffModeBite).
			WithRotateTime(rotatefile.EveryHour).
			Build(),

		handler.NewBuilder().
			WithLogfile("./tmp/info.log").
			WithLogLevels(slog.NormalLevels).
			WithBuffSize(1024*10).
			WithBuffMode(handler.BuffModeBite).
			WithRotateTime(rotatefile.EveryHour).
			Build(),
	)
	return &Logger{
		slog: logger,
	}
}

type Logger struct {
	slog *slog.Logger
}

func (logger *Logger) Error(args ...any) {
	logger.slog.Error(args...)
}

func (logger *Logger) Info(args ...any) {
	logger.slog.Info(args...)
}
