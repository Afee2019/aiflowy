package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger
var sugar *zap.SugaredLogger

// Config holds logger configuration
type Config struct {
	Level      string
	Format     string
	Output     string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

// Init initializes the logger
func Init(cfg *Config) error {
	// Parse log level
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Configure encoder
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if cfg.Format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Configure output
	var writeSyncer zapcore.WriteSyncer
	if cfg.Output == "file" {
		writeSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   true,
		})
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = log.Sugar()

	return nil
}

// Get returns the logger instance
func Get() *zap.Logger {
	return log
}

// Sugar returns the sugared logger
func Sugar() *zap.SugaredLogger {
	return sugar
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return log.Sync()
}
