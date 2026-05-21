package corezaplogger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	corelogger "github.com/egotk/golang-advert-app/internal/core/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zapLogger *zap.Logger

	file *os.File
}

func New(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())
	return &Logger{
		zapLogger: zapLogger,
		file:      logFile,
	}, nil
}

func (l *Logger) Debug(msg string, fields ...corelogger.Field) {
	fieldsZap := fieldsToZap(fields...)
	l.zapLogger.Debug(msg, fieldsZap...)
}

func (l *Logger) Info(msg string, fields ...corelogger.Field) {
	fieldsZap := fieldsToZap(fields...)
	l.zapLogger.Info(msg, fieldsZap...)
}

func (l *Logger) Warn(msg string, fields ...corelogger.Field) {
	fieldsZap := fieldsToZap(fields...)
	l.zapLogger.Warn(msg, fieldsZap...)
}

func (l *Logger) Error(msg string, fields ...corelogger.Field) {
	fieldsZap := fieldsToZap(fields...)
	l.zapLogger.Error(msg, fieldsZap...)
}

func (l *Logger) DPanic(msg string, fields ...corelogger.Field) {
	fieldsZap := fieldsToZap(fields...)
	l.zapLogger.DPanic(msg, fieldsZap...)
}

func (l *Logger) Panic(msg string, fields ...corelogger.Field) {
	fieldsZap := fieldsToZap(fields...)
	l.zapLogger.Panic(msg, fieldsZap...)
}

func (l *Logger) Fatal(msg string, fields ...corelogger.Field) {
	fieldsZap := fieldsToZap(fields...)
	l.zapLogger.Fatal(msg, fieldsZap...)
}

func (l *Logger) With(field ...corelogger.Field) corelogger.Logger {
	fieldsZap := fieldsToZap(field...)

	return &Logger{
		zapLogger: l.zapLogger.With(fieldsZap...),
		file:      l.file,
	}
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}

func fieldsToZap(fields ...corelogger.Field) []zap.Field {
	fieldsZap := make([]zap.Field, len(fields))

	for idx, field := range fields {
		fieldsZap[idx] = zap.Any(field.Key, field.Value)
	}

	return fieldsZap
}
