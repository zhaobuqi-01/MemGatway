package logger

import (
	"gateway/configs"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger

	logConfig *configs.LogConfig
)

func init() {
	logConfig = configs.GetLogConfig()
	initLogger()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getFileLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   logConfig.Filename,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		LocalTime:  true,
		Compress:   logConfig.Compress,
	}
}

func getErrorFileLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   logConfig.Filename,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		LocalTime:  true,
		Compress:   logConfig.Compress,
	}
}

func getZapCoreForLevel(level zapcore.Level) zapcore.Core {
	encoder := getEncoder()
	fileLogger := getFileLogger()

	logLevel := zap.InfoLevel

	logLevelStr := logConfig.Level
	if logLevelStr == "debug" {
		logLevel = zap.DebugLevel
	}

	if level == zapcore.ErrorLevel {
		errorFileLogger := getErrorFileLogger()
		return zapcore.NewCore(encoder, zapcore.AddSync(errorFileLogger), zapcore.ErrorLevel)
	}
	return zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(fileLogger), logLevel),
	)
}

func initLogger() {
	nonErrorCore := getZapCoreForLevel(zapcore.InfoLevel)
	errorCore := getZapCoreForLevel(zapcore.ErrorLevel)
	core := zapcore.NewTee(nonErrorCore, errorCore)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// 重定向标准库的日志输出到zap
	zap.RedirectStdLog(logger)

	// 替换全局logger
	zap.ReplaceGlobals(logger)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
