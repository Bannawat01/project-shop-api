package logger

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm/logger"
)

type EchoLoggerAdapter struct {
	logger echo.Logger
}

func NewEchoLoggerAdapter(echoLogger echo.Logger) logger.Interface {
	return &EchoLoggerAdapter{logger: echoLogger}
}

func (l *EchoLoggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	// For simplicity, return self. Could adjust echo logger level if needed.
	return l
}

func (l *EchoLoggerAdapter) Info(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l *EchoLoggerAdapter) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

func (l *EchoLoggerAdapter) Error(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

func (l *EchoLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// For GORM trace, log the SQL if no error, or error.
	sql, rows := fc()
	if err != nil {
		l.logger.Errorf("SQL Error: %s, Rows: %d, Error: %v", sql, rows, err)
	} else {
		l.logger.Debugf("SQL: %s, Rows: %d", sql, rows)
	}
}
