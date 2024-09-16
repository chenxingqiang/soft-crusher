// File: pkg/logging/logging.go

package logging

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// InitLogger initializes the global logger
func InitLogger(debug bool, logFile string) {
	once.Do(func() {
		var cfg zap.Config
		if debug {
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			cfg = zap.NewProductionConfig()
		}
		var err error
		logger, err = cfg.Build()
		if err != nil {
			panic(err)
		}
	})

	if logFile != "" {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			w,
			cfg.Level,
		)
		logger = zap.New(core)
	}
}

// SetLogger sets a custom logger
func SetLogger(l *zap.Logger) {
	logger = l
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Fatal logs a fatal message and then calls os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// With creates a child logger with the given fields
func With(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return logger.Sync()
}

// GetLogger returns the current logger instance
func GetLogger() *zap.Logger {
	return logger
}
