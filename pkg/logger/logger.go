package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger
type Logger struct {
	*zap.SugaredLogger
}

// New creates a new logger instance
func New() *Logger {
	config := zap.NewProductionConfig()
	
	// Set log level from environment
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	
	// Configure output format
	if os.Getenv("APP_ENV") == "development" {
		config.Encoding = "console"
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	
	// Build logger
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

// WithFields adds fields to logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		SugaredLogger: l.With(args...),
	}
}

// WithError adds error field to logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		SugaredLogger: l.With("error", err.Error()),
	}
}
