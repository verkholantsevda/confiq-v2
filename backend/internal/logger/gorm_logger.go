package logger

import (
	"context"
	"log/slog"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct{}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	slog.Debug(msg, "args", args)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	slog.Warn(msg, "args", args)
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	slog.Error(msg, "args", args)
}

func (l *GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	if !DebugEnabled() {
		return
	}

	sql, rows := fc()

	elapsed := time.Since(begin)

	if err != nil {
		slog.Error(
			"sql failed",
			"duration", elapsed,
			"rows", rows,
			"query", sql,
			"error", err,
		)
		return
	}

	slog.Debug(
		"sql",
		"duration", elapsed,
		"rows", rows,
		"query", sql,
	)
}
